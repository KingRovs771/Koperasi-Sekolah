package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func GetPos(c *fiber.Ctx) error {
	sess, _ := config.Store.Get(c)
	namaLengkap := sess.Get("NamaLengkap")
	if namaLengkap == nil {
		namaLengkap = "Kasir"
	}
	return c.Render("pages/kasir", fiber.Map{
		"title":       "Kasir POS",
		"namaLengkap": namaLengkap,
		"RandomID":    time.Now().Unix(),
	})
}

func GetProductsPOS(c *fiber.Ctx) error {
	var products []models.Product

	if err := config.DB.Where("stock > 0").Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.JSON(products)
}

func CheckOutProcess(c *fiber.Ctx) error {
	req := new(models.CheckoutRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Data"})
	}

	sess, _ := config.Store.Get(c)
	userUID := sess.Get("users_uid")

	tx := config.DB.Begin()

	trxUID := uuid.New().String()
	noFaktur := fmt.Sprintf("TRX-%d", time.Now().Unix())

	var totalBelanja float64 = 0
	var totalModal float64 = 0

	for _, item := range req.Items {
		var product models.Product

		// Lock row product agar tidak race condition (rebutan stok)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("products_uid = ?", item.ProductsUID).
			First(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"error": "Produk tidak ditemukan: " + item.ProductsUID})
		}

		if product.Stock < item.Quantity {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Stok tidak cukup untuk: " + product.NamaProducts})
		}

		subTotal := float64(product.HargaJual) * float64(item.Quantity)
		subModal := float64(product.HargaBeli) * float64(item.Quantity)

		totalBelanja += subTotal
		totalModal += subModal

		product.Stock -= item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal update stok"})
		}

		trxItem := models.TransactionsItems{
			TransactionsUID: trxUID,
			ProductsUID:     product.ProductsUID,
			Jumlah:          item.Quantity,
			HargaSaatIni:    float64(product.HargaJual),
			ModalSaatIni:    float64(product.HargaBeli),
			SubTotal:        subTotal,
			CreatedAt:       time.Now(),
		}
		if err := tx.Create(&trxItem).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan item transaksi"})
		}
	}

	kembalian := req.UangDiterima - totalBelanja
	if kembalian < 0 {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Uang pembayaran kurang"})
	}

	trx := models.Transactions{
		TransactionsUID:  trxUID,
		NoFaktur:         noFaktur,
		UserUID:          userUID,
		TotalBelanja:     totalBelanja,
		TotalModal:       totalModal,
		MetodePembayaran: req.MetodePembayaran,
		UangDiterima:     req.UangDiterima,
		Kembalian:        kembalian,
		CreatedAt:        time.Now(),
	}

	if err := tx.Create(&trx).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan transaksi"})
	}

	tx.Commit()

	return c.JSON(fiber.Map{
		"message":   "Transaksi Berhasil",
		"no_faktur": noFaktur,
		"kembalian": kembalian,
		"total":     totalBelanja,
	})
}
