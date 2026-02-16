package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type ReportSummary struct {
	TotalPenjualan  float64
	LabaBersih      float64
	JumlahTransaksi int64
}

type CategoryProfit struct {
	NamaCategory string  `gorm:"column:nama_category"`
	Pendapatan   float64 `gorm:"column:pendapatan"`
	Biaya        float64 `gorm:"column:biaya"`
	Laba         float64 `gorm:"column:laba"`
	Margin       float64 `gorm:"-"`
}

type RecentTrxDTO struct {
	CreatedAt        time.Time
	NoFaktur         string
	UserUID          string
	MetodePembayaran string
	TotalBelanja     float64
	Laba             float64
}

func GetReports(c *fiber.Ctx) error {
	var summary ReportSummary

	// A. Hitung Total Penjualan & Laba
	config.DB.Model(&models.Transactions{}).
		Select("COALESCE(SUM(total_belanja), 0) as total_penjualan, COALESCE(SUM(total_belanja - total_modal), 0) as laba_bersih, COUNT(*) as jumlah_transaksi").
		Scan(&summary)

	// B. Ambil 10 Transaksi Terakhir (Mapping ke DTO agar Laba terhitung otomatis)
	var recentTrx []models.Transactions
	config.DB.Order("created_at desc").Limit(10).Find(&recentTrx)

	var recentTrxList []RecentTrxDTO
	for _, t := range recentTrx {
		recentTrxList = append(recentTrxList, RecentTrxDTO{
			CreatedAt:        t.CreatedAt,
			NoFaktur:         t.NoFaktur,
			UserUID:          t.UserUID,
			MetodePembayaran: t.MetodePembayaran,
			TotalBelanja:     t.TotalBelanja,
			Laba:             t.TotalBelanja - t.TotalModal,
		})
	}

	// C. Hitung Laba Rugi Per Kategori
	var catProfits []CategoryProfit
	config.DB.Table("transactions_items").
		Select("products.category_uid as nama_category, COALESCE(SUM(transactions_items.sub_total), 0) as pendapatan, COALESCE(SUM(transactions_items.modal_saat_ini * transactions_items.jumlah), 0) as biaya, COALESCE(SUM(transactions_items.sub_total) - SUM(transactions_items.modal_saat_ini * transactions_items.jumlah), 0) as laba").
		Joins("JOIN products ON products.products_uid = transactions_items.products_uid").
		Group("products.category_uid").
		Scan(&catProfits)

	var totalPendapatan, totalBiaya, totalLaba, totalMargin float64
	for i := range catProfits {
		if catProfits[i].Pendapatan > 0 {
			catProfits[i].Margin = (catProfits[i].Laba / catProfits[i].Pendapatan) * 100
		}
		totalPendapatan += catProfits[i].Pendapatan
		totalBiaya += catProfits[i].Biaya
		totalLaba += catProfits[i].Laba
	}
	if totalPendapatan > 0 {
		totalMargin = (totalLaba / totalPendapatan) * 100
	}

	// D. Data Grafik Tren (14 Hari Terakhir)
	fourteenDaysAgo := time.Now().AddDate(0, 0, -14)
	var trxLast14 []models.Transactions
	config.DB.Where("created_at >= ?", fourteenDaysAgo).Find(&trxLast14)

	type DailyAgg struct {
		Revenue float64
		Profit  float64
	}
	dailyData := make(map[string]DailyAgg)

	for _, t := range trxLast14 {
		dateStr := t.CreatedAt.Format("02 Jan")
		agg := dailyData[dateStr]
		agg.Revenue += t.TotalBelanja
		agg.Profit += (t.TotalBelanja - t.TotalModal)
		dailyData[dateStr] = agg
	}

	var labels []string
	var revenues []float64
	var profits []float64

	for i := 13; i >= 0; i-- {
		dateStr := time.Now().AddDate(0, 0, -i).Format("02 Jan")
		labels = append(labels, dateStr)
		revenues = append(revenues, dailyData[dateStr].Revenue)
		profits = append(profits, dailyData[dateStr].Profit)
	}

	if c.Get("HX-Request") == "true" {
		// Jika via HTMX (klik dari sidebar), render isinya saja tanpa layout utama
		return c.Render("pages/report", fiber.Map{
			"Title":           "Laporan & Analytics",
			"TotalPenjualan":  summary.TotalPenjualan,
			"LabaBersih":      summary.LabaBersih,
			"JumlahTransaksi": summary.JumlahTransaksi,
			"RecentTrx":       recentTrxList,

			"CatProfits":      catProfits,
			"TotalPendapatan": totalPendapatan,
			"TotalBiaya":      totalBiaya,
			"TotalLaba":       totalLaba,
			"TotalMargin":     totalMargin,

			"ChartLabels":   labels,
			"ChartRevenues": revenues,
			"ChartProfits":  profits,

			"Role":        c.Locals("Role"),
			"NamaLengkap": c.Locals("NamaLengkap"),
		})
	}

	// E. Render HTML
	return c.Render("pages/report", fiber.Map{
		"Title":           "Laporan & Analytics",
		"TotalPenjualan":  summary.TotalPenjualan,
		"LabaBersih":      summary.LabaBersih,
		"JumlahTransaksi": summary.JumlahTransaksi,
		"RecentTrx":       recentTrxList,

		"CatProfits":      catProfits,
		"TotalPendapatan": totalPendapatan,
		"TotalBiaya":      totalBiaya,
		"TotalLaba":       totalLaba,
		"TotalMargin":     totalMargin,

		"ChartLabels":   labels,
		"ChartRevenues": revenues,
		"ChartProfits":  profits,

		"Role":        c.Locals("Role"),
		"NamaLengkap": c.Locals("NamaLengkap"),
	}, "layouts/main")
}

func ExportReport(c *fiber.Ctx) error {
	reportType := c.Query("type") // daily, weekly, monthly
	format := c.Query("format")   // pdf, excel
	dateStr := c.Query("date")    // 2026-01-22 atau 2026-01 (untuk monthly)

	var transactions []models.Transactions
	var startDate, endDate time.Time
	var err error

	// Filter Logika Waktu
	switch reportType {
	case "daily":
		startDate, err = time.Parse("2006-01-02", dateStr)
		endDate = startDate.Add(24 * time.Hour)
	case "monthly":
		startDate, err = time.Parse("2006-01", dateStr)
		endDate = startDate.AddDate(0, 1, 0)
	default:
		// Default hari ini
		startDate = time.Now().Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	}

	if err != nil {
		return c.Status(400).SendString("Format tanggal salah")
	}

	// Ambil Data dari DB
	config.DB.Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Find(&transactions)

	if len(transactions) == 0 {
		return c.Status(404).SendString("Tidak ada data transaksi pada periode ini")
	}

	// Generate File berdasarkan Format
	if format == "excel" {
		return generateExcel(c, transactions, reportType, dateStr)
	} else if format == "pdf" {
		return generatePDF(c, transactions, reportType, dateStr)
	}

	return c.Status(400).SendString("Format tidak didukung")
}

func generateExcel(c *fiber.Ctx, data []models.Transactions, reportType, dateStr string) error {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// Header
	headers := []string{"No", "Tanggal", "No Faktur", "Metode", "Total Modal", "Total Belanja", "Laba"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	// Isi Data
	var totalOmset, totalLaba float64
	for i, trx := range data {
		row := i + 2
		laba := trx.TotalBelanja - trx.TotalModal
		totalOmset += trx.TotalBelanja
		totalLaba += laba

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), trx.CreatedAt.Format("2006-01-02 15:04"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), trx.NoFaktur)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), trx.MetodePembayaran)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), trx.TotalModal)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), trx.TotalBelanja)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), laba)
	}

	// Footer Total
	lastRow := len(data) + 2
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", lastRow), "TOTAL")
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", lastRow), totalOmset)
	f.SetCellValue(sheetName, fmt.Sprintf("G%d", lastRow), totalLaba)

	// Set Header Response
	filename := fmt.Sprintf("Laporan_%s_%s.xlsx", reportType, dateStr)
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return f.Write(c.Response().BodyWriter())
}

// --- FUNGSI GENERATE PDF ---
func generatePDF(c *fiber.Ctx, data []models.Transactions, reportType, dateStr string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Judul
	pdf.Cell(40, 10, fmt.Sprintf("Laporan Penjualan %s", reportType))
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, fmt.Sprintf("Periode: %s", dateStr))
	pdf.Ln(12)

	// Header Tabel
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(10, 10, "No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 10, "Waktu", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "No Faktur", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 10, "Metode", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 10, "Total", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 10, "Laba", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Isi Data
	pdf.SetFont("Arial", "", 9)
	var totalOmset, totalLaba float64

	for i, trx := range data {
		laba := trx.TotalBelanja - trx.TotalModal
		totalOmset += trx.TotalBelanja
		totalLaba += laba

		pdf.CellFormat(10, 8, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, trx.CreatedAt.Format("02/01 15:04"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(45, 8, trx.NoFaktur, "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, trx.MetodePembayaran, "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%.0f", trx.TotalBelanja), "1", 0, "R", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%.0f", laba), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}

	// Total
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(115, 10, "TOTAL PENDAPATAN", "1", 0, "R", true, 0, "")
	pdf.CellFormat(35, 10, fmt.Sprintf("%.0f", totalOmset), "1", 0, "R", true, 0, "")
	pdf.CellFormat(35, 10, fmt.Sprintf("%.0f", totalLaba), "1", 0, "R", true, 0, "")

	// Output
	filename := fmt.Sprintf("Laporan_%s_%s.pdf", reportType, dateStr)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return pdf.Output(c.Response().BodyWriter())
}
