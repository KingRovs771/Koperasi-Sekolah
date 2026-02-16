package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetManageProducts(c *fiber.Ctx) error {
	var products []models.Product
	var categories []models.Category // <--- Tambahkan Variable ini

	// 1. Ambil Produk
	if err := config.DB.Order("created_at desc").Find(&products).Error; err != nil {
		return c.Status(500).SendString("Gagal memuat data produk")
	}

	// 2. AMBIL KATEGORI DARI DATABASE
	if err := config.DB.Find(&categories).Error; err != nil {
		return c.Status(500).SendString("Gagal memuat data kategori")
	}

	// 3. Perhitungan Statistik (Kode lama Anda tetap di sini)
	var totalProduk int64
	var stokMenipis int64
	var stokHabis int64
	var totalAset int64

	for _, p := range products {
		totalProduk++
		if p.Stock == 0 {
			stokHabis++
		} else if p.Stock <= p.StockMin {
			stokMenipis++
		}
		totalAset += (p.HargaBeli * p.Stock)
	}

	// 4. Masukkan ke Fiber Map
	data := fiber.Map{
		"Title":        "Manajemen Produk",
		"Products":     products,
		"Categories":   categories, // <--- Kirim data kategori ke HTML
		"TotalProduk":  totalProduk,
		"StokMenipis":  stokMenipis,
		"StokHabis":    stokHabis,
		"TotalAset":    totalAset,
		"Role":         c.Locals("Role"),
		"Nama_lengkap": c.Locals("NamaLengkap"),
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("pages/manage_produk", data)
	}
	return c.Render("pages/manage_produk", data, "layouts/main")
}

func CreateProduct(c *fiber.Ctx) error {
	// 1. Parsing Form Data manual karena ada File Upload
	nama := c.FormValue("nama_products")
	kode := c.FormValue("kode_produk")
	kategori := c.FormValue("category_uid")
	hargaJual, _ := strconv.Atoi(c.FormValue("harga_jual"))
	hargaBeli, _ := strconv.Atoi(c.FormValue("harga_beli"))
	stock, _ := strconv.Atoi(c.FormValue("stock"))
	stockMin, _ := strconv.Atoi(c.FormValue("stock_min"))
	desc := c.FormValue("description")

	// 2. Handle Image Upload
	var imagePath string
	file, err := c.FormFile("image")
	if err == nil {
		// Generate unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		// Simpan ke folder uploads (Pastikan folder ini ada!)
		if err := c.SaveFile(file, fmt.Sprintf("./uploads/%s", filename)); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Gagal Upload Gambar",
			})
		}
		imagePath = filename
	}

	// 3. Create Struct
	product := models.Product{
		ProductsUID:  uuid.New().String(),
		KodeProduk:   kode,
		CategoryUID:  kategori,
		NamaProducts: nama,
		Description:  desc,
		HargaJual:    int64(hargaJual),
		HargaBeli:    int64(hargaBeli),
		Stock:        int64(stock),
		StockMin:     int64(stockMin),
		Image:        imagePath,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Gagal menyimpan database (Kode Produk mungkin duplikat)",
		})
	}

	return c.JSON(fiber.Map{"message": "Produk berhasil ditambahkan"})
}

func GetProductByUID(c *fiber.Ctx) error {
	uid := c.Params("uid")
	var product models.Product
	if err := config.DB.Where("products_uid = ?", uid).First(&product).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Gagal menyimpan database",
		})
	}
	return c.JSON(product)
}

// UPDATE: Edit Produk
func UpdateProduct(c *fiber.Ctx) error {
	uid := c.Params("uid")
	var product models.Product

	if err := config.DB.Where("products_uid = ?", uid).First(&product).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Produk tidak ditemukan",
		})
	}

	// Update Field Dasar
	product.NamaProducts = c.FormValue("nama_products")
	product.KodeProduk = c.FormValue("kode_produk")
	product.CategoryUID = c.FormValue("category_uid")
	product.Description = c.FormValue("description")

	hJual, _ := strconv.Atoi(c.FormValue("harga_jual"))
	product.HargaJual = int64(hJual)

	hBeli, _ := strconv.Atoi(c.FormValue("harga_beli"))
	product.HargaBeli = int64(hBeli)

	stk, _ := strconv.Atoi(c.FormValue("stock"))
	product.Stock = int64(stk)

	stkMin, _ := strconv.Atoi(c.FormValue("stock_min"))
	product.StockMin = int64(stkMin)

	// Cek apakah ada gambar baru yg diupload
	file, err := c.FormFile("image")
	if err == nil {
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		if err := c.SaveFile(file, fmt.Sprintf("./uploads/%s", filename)); err == nil {
			product.Image = filename // Ganti nama file di DB
			// Note: Idealnya hapus file lama di sini
		}
	}

	config.DB.Save(&product)
	return c.JSON(fiber.Map{"message": "Produk berhasil diperbarui"})
}

// DELETE: Hapus Produk
func DeleteProduct(c *fiber.Ctx) error {
	uid := c.Params("uid")
	if err := config.DB.Where("products_uid = ?", uid).Delete(&models.Product{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal Menghapus Produk",
		})
	}
	return c.JSON(fiber.Map{"message": "Produk dihapus"})
}
