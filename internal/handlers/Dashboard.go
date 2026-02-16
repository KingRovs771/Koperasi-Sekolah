package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TopProductDTO struct {
	NamaProducts string  `gorm:"column:nama_products"`
	Kategori     string  `gorm:"column:kategori"`
	Terjual      int64   `gorm:"column:terjual"`
	Pendapatan   float64 `gorm:"column:pendapatan"`
}

func GetDashboard(c *fiber.Ctx) error {
	// 1. Waktu Batasan (Hari Ini & Minggu Ini)
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)
	weekStart := todayStart.AddDate(0, 0, -7)

	// 2. Metrik Kartu Atas
	var salesToday float64
	var countTrx int64
	var productsSold int64
	var lowStockCount int64

	// Penjualan & Transaksi Hari Ini
	config.DB.Model(&models.Transactions{}).
		Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
		Select("COALESCE(SUM(total_belanja), 0)").Scan(&salesToday)

	config.DB.Model(&models.Transactions{}).
		Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
		Count(&countTrx)

	// Produk Terjual Hari Ini
	config.DB.Table("transactions_items").
		Joins("JOIN transactions ON transactions.transactions_uid = transactions_items.transactions_uid").
		Where("transactions.created_at >= ? AND transactions.created_at < ?", todayStart, todayEnd).
		Select("COALESCE(SUM(transactions_items.jumlah), 0)").Scan(&productsSold)

	// Hitung Stok Rendah (Misal: logic stock <= stock_min)
	// Jika tabel tidak ada stock_min, asumsikan stock <= 5
	config.DB.Model(&models.Product{}).Where("stock <= stock_min OR stock <= 5").Count(&lowStockCount)

	// 3. Peringatan Stok Rendah (List Maksimal 5)
	var lowStockAlerts []models.Product
	config.DB.Where("stock <= stock_min OR stock <= 5").Limit(5).Find(&lowStockAlerts)

	// 4. Produk Terlaris Minggu Ini (Top 5)
	var topProducts []TopProductDTO
	config.DB.Table("transactions_items").
		Select("products.nama_products, products.category_uid as kategori, SUM(transactions_items.jumlah) as terjual, SUM(transactions_items.sub_total) as pendapatan").
		Joins("JOIN products ON products.products_uid = transactions_items.products_uid").
		Joins("JOIN transactions ON transactions.transactions_uid = transactions_items.transactions_uid").
		Where("transactions.created_at >= ?", weekStart).
		Group("products.products_uid, products.nama_products, products.category_uid").
		Order("terjual DESC").
		Limit(5).
		Scan(&topProducts)

	// 5. Data Chart Penjualan (14 Hari Terakhir)
	fourteenDaysAgo := todayStart.AddDate(0, 0, -13)
	var trxLast14 []models.Transactions
	config.DB.Where("created_at >= ?", fourteenDaysAgo).Find(&trxLast14)

	dailyRevenue := make(map[string]float64)
	for _, t := range trxLast14 {
		dateStr := t.CreatedAt.Format("02 Jan")
		dailyRevenue[dateStr] += t.TotalBelanja
	}

	var chartLabels []string
	var chartData []float64

	for i := 13; i >= 0; i-- {
		dateStr := time.Now().AddDate(0, 0, -i).Format("02 Jan")
		chartLabels = append(chartLabels, dateStr)
		chartData = append(chartData, dailyRevenue[dateStr])
	}

	// 6. Siapkan Data untuk HTML
	data := fiber.Map{
		"Title":          "Dashboard",
		"Role":           c.Locals("Role"),
		"Nama_lengkap":   c.Locals("NamaLengkap"),
		"SalesToday":     salesToday,
		"ProductsSold":   productsSold,
		"LowStockCount":  lowStockCount,
		"CountTrx":       countTrx,
		"LowStockAlerts": lowStockAlerts,
		"TopProducts":    topProducts,
		"ChartLabels":    chartLabels,
		"ChartData":      chartData,
	}

	// Pengecekan HTMX agar tidak merender ulang seluruh layout
	if c.Get("HX-Request") == "true" {
		return c.Render("pages/dashboard", data)
	}
	return c.Render("pages/dashboard", data, "layouts/main")
}
