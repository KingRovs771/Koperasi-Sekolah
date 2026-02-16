package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetStock(c *fiber.Ctx) error {
	var products []models.Product

	if err := config.DB.Order("created_at desc").Find(&products).Error; err != nil {
		return c.Status(500).SendString("Gagal memuat produk")
	}

	var statAvailable, statLow, statOut int
	for _, p := range products {
		if p.Stock <= 0 {
			statOut++
		} else if p.Stock <= p.StockMin {
			statLow++
		} else {
			statAvailable++
		}
	}

	var history []models.MutasiHistoryDTO
	config.DB.Table("mutasi_stocks").
		Select("mutasi_stocks.*, products.nama_products as nama_produk, users.nama_lengkap as nama_user").
		Joins("LEFT JOIN products ON products.products_uid = mutasi_stocks.products_uid").
		Joins("LEFT JOIN users ON users.users_uid = mutasi_stocks.user_uid").
		Order("mutasi_stocks.created_at desc").
		Limit(10). // Batasi 10 riwayat terakhir
		Scan(&history)

	data := fiber.Map{
		"Title":         "Manajemen Stok",
		"Products":      products,
		"History":       history,
		"StatAvailable": statAvailable,
		"StatLow":       statLow,
		"StatOut":       statOut,
		"Role":          c.Locals("Role"),
		"NamaLengkap":   c.Locals("NamaLengkap"),
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("pages/manage_stock", data)
	}
	return c.Render("pages/manage_stock", data, "layouts/main")
}

func AddStock(c *fiber.Ctx) error {
	type AddStockInput struct {
		ProductsUID string `json:"product_uid" form:"products_uid"`
		Quantity    int64  `json:"quantity" form:"quantity"`
		Notes       string `json:"notes" form:"notes"`
	}

	input := new(AddStockInput)

	if err := c.BodyParser(input); err != nil {
		log.Println("Error Parsing:", err.Error())
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	log.Printf("DEBUG ADD STOCK: UID=%s, Qty=%d", input.ProductsUID, input.Quantity)

	if input.ProductsUID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Produk tidak boleh kosong"})
	}
	if input.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Jumlah harus lebih dari 0"})
	}

	userVal := c.Locals("users_uid")
	if userVal == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Sesi habis, silakan login ulang"})
	}
	userUID, ok := userVal.(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca data user"})
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var product models.Product

		if err := tx.Where("products_uid = ?", input.ProductsUID).First(&product).Error; err != nil {
			log.Println("Product Not Found:", err.Error())
			return err
		}

		product.Stock += input.Quantity

		if err := tx.Save(&product).Error; err != nil {
			log.Println("Save Product Error:", err.Error())
			return err
		}

		mutasi := models.MutasiStock{
			MutasiStocksUID: uuid.New().String(),
			ProductsUID:     input.ProductsUID,
			UserUID:         userUID,
			Tipe:            "masuk",
			Jumlah:          input.Quantity,
			Catatan:         input.Notes,
			CreatedAt:       time.Now(),
		}

		if err := tx.Create(&mutasi).Error; err != nil {
			log.Println("Create Mutasi Error:", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memproses transaksi: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Stok berhasil ditambahkan",
	})
}

func ReduceStock(c *fiber.Ctx) error {
	// 1. Struct Input (Sama seperti AddStock)
	type ReduceStockInput struct {
		ProductUID string `json:"product_uid"`
		Quantity   int64  `json:"quantity"`
		Notes      string `json:"notes"`
	}

	input := new(ReduceStockInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// 2. Validasi Dasar
	if input.ProductUID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Produk harus dipilih"})
	}
	if input.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Jumlah harus lebih dari 0"})
	}

	// 3. Ambil User Session
	userVal := c.Locals("users_uid")
	if userVal == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Sesi habis, silakan login ulang"})
	}
	userUID := userVal.(string)

	// 4. Transaksi Database
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var product models.Product

		// Lock baris produk
		if err := tx.Where("products_uid = ?", input.ProductUID).First(&product).Error; err != nil {
			return err
		}

		// --- VALIDASI STOK CUKUP ATAU TIDAK ---
		if product.Stock < input.Quantity {
			// Return error khusus jika stok kurang
			return fiber.NewError(fiber.StatusBadRequest, "Stok tidak mencukupi! Sisa stok: "+string(rune(product.Stock)))
		}

		// Kurangi Stok
		product.Stock -= input.Quantity

		// Simpan Perubahan Stok
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// Catat Mutasi (Tipe: keluar)
		mutasi := models.MutasiStock{
			MutasiStocksUID: uuid.New().String(),
			ProductsUID:     input.ProductUID,
			UserUID:         userUID,
			Tipe:            "keluar", // <--- PENTING
			Jumlah:          input.Quantity,
			Catatan:         input.Notes,
			CreatedAt:       time.Now(),
		}

		if err := tx.Create(&mutasi).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		// Cek jika errornya karena stok kurang
		if e, ok := err.(*fiber.Error); ok {
			return c.Status(e.Code).JSON(fiber.Map{"error": e.Message})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update stok"})
	}

	return c.JSON(fiber.Map{"message": "Stok berhasil dikurangi"})
}
