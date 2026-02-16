package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const BackupDir = "./storage/backups"

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// 1. GET: Tampilkan Halaman
func GetBackupRestore(c *fiber.Ctx) error {
	var backups []models.BackupRestore
	config.DB.Order("created_at desc").Find(&backups)

	var totalSizeBytes int64
	for _, b := range backups {
		totalSizeBytes += b.UkuranFile
	}

	totalSizeMB := float64(totalSizeBytes) / (1024 * 1024)
	maxSizeMB := 10240.0
	progressPercent := (totalSizeMB / maxSizeMB) * 100
	data := fiber.Map{
		"Title":           "Backup & Restore",
		"Backups":         backups,
		"TotalSizeMB":     fmt.Sprintf("%.2f", totalSizeMB),
		"ProgressPercent": fmt.Sprintf("%.2f", progressPercent),
		"Role":            c.Locals("Role"),
		"NamaLengkap":     c.Locals("NamaLengkap"),
	}

	// --- PENTING: PENGECEKAN HTMX ---
	if c.Get("HX-Request") == "true" {
		// Jika via HTMX (klik dari sidebar), render isinya saja tanpa layout utama
		return c.Render("pages/backup_restore", data)
	}

	// Jika di-refresh manual (F5), render dengan layout utama
	return c.Render("pages/backup_restore", data, "layouts/main")
}

// 2. POST: CREATE BACKUP VIA DOCKER
func CreateBackup(c *fiber.Ctx) error {
	os.MkdirAll(BackupDir, os.ModePerm)

	fileNameInput := c.FormValue("filename")
	if fileNameInput == "" {
		fileNameInput = fmt.Sprintf("backup-%s", time.Now().Format("2006-01-02-150405"))
	}
	fileName := fileNameInput + ".sql"
	filePath := filepath.Join(BackupDir, fileName)

	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "password_db_anda")
	dbName := getEnv("DB_NAME", "kopsis_db")
	containerName := getEnv("DB_CONTAINER_NAME", "postgres_container") // Wajib sama dengan di docker ps

	// Perintah: docker exec -i -e PGPASSWORD=pass container_name pg_dump -U user -d dbname --clean
	cmd := exec.Command("docker", "exec", "-i",
		"-e", fmt.Sprintf("PGPASSWORD=%s", dbPass),
		containerName,
		"pg_dump", "-U", dbUser, "-d", dbName, "--clean",
	)

	// Arahkan hasil output dari Docker langsung ke file lokal Windows
	outFile, err := os.Create(filePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat file backup lokal"})
	}
	defer outFile.Close()

	cmd.Stdout = outFile // Pipe output

	if err := cmd.Run(); err != nil {
		os.Remove(filePath)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal backup! Pastikan nama DB_CONTAINER_NAME benar dan Docker menyala. " + err.Error()})
	}

	// Simpan riwayat ke Database
	fileInfo, _ := os.Stat(filePath)
	backupRecord := models.BackupRestore{
		BackupRestoreUID: uuid.New().String(),
		NamaFile:         fileName,
		UkuranFile:       fileInfo.Size(),
		CreatedAt:        time.Now(),
	}

	if err := config.DB.Create(&backupRecord).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan histori ke database"})
	}

	return c.JSON(fiber.Map{"message": "Backup database berhasil dibuat dari Docker!"})
}

// 3. POST: RESTORE BACKUP VIA DOCKER
func RestoreBackup(c *fiber.Ctx) error {
	fileName := c.FormValue("filename")
	if fileName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Nama file tidak boleh kosong"})
	}

	filePath := filepath.Join(BackupDir, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{"error": "File backup tidak ditemukan"})
	}

	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "password_db_anda")
	dbName := getEnv("DB_NAME", "kopsis_db")
	containerName := getEnv("DB_CONTAINER_NAME", "postgres_container")

	// Perintah: docker exec -i -e PGPASSWORD=pass container_name psql -U user -d dbname
	cmd := exec.Command("docker", "exec", "-i",
		"-e", fmt.Sprintf("PGPASSWORD=%s", dbPass),
		containerName,
		"psql", "-U", dbUser, "-d", dbName,
	)

	// Buka file backup dari lokal Windows
	inFile, err := os.Open(filePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuka file backup lokal"})
	}
	defer inFile.Close()

	// Masukkan isi file Windows ke dalam input perintah Docker psql
	cmd.Stdin = inFile

	if err := cmd.Run(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal restore database! " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Database berhasil di-restore!"})
}

// 4. DELETE BACKUP
func DeleteBackup(c *fiber.Ctx) error {
	fileName := c.FormValue("filename")

	if err := config.DB.Where("nama_file = ?", fileName).Delete(&models.BackupRestore{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus log dari database"})
	}

	filePath := filepath.Join(BackupDir, fileName)
	os.Remove(filePath)

	return c.JSON(fiber.Map{"message": "File backup berhasil dihapus"})
}

// 5. DOWNLOAD BACKUP
func DownloadBackup(c *fiber.Ctx) error {
	fileName := c.Params("filename")
	filePath := filepath.Join(BackupDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).SendString("File tidak ditemukan")
	}

	return c.Download(filePath, fileName)
}
func UploadAndRestore(c *fiber.Ctx) error {
	// 1. Terima file dari form
	file, err := c.FormFile("backup_file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal membaca file upload"})
	}

	// Pastikan folder penyimpanan ada
	os.MkdirAll(BackupDir, os.ModePerm)
	filePath := filepath.Join(BackupDir, file.Filename)

	// 2. Simpan file fisik ke server Windows
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan file ke folder server"})
	}

	// 3. Catat di Database (agar muncul di tabel daftar backup)
	backupRecord := models.BackupRestore{
		BackupRestoreUID: uuid.New().String(),
		NamaFile:         file.Filename,
		UkuranFile:       file.Size,
		CreatedAt:        time.Now(),
	}
	config.DB.Create(&backupRecord) // Hiraukan error minor jika nama file duplikat

	// 4. Eksekusi Restore via Docker (Sama seperti Restore biasa)
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "password_db_anda")
	dbName := getEnv("DB_NAME", "kopsis_db")
	containerName := getEnv("DB_CONTAINER_NAME", "postgres_container")

	cmd := exec.Command("docker", "exec", "-i",
		"-e", fmt.Sprintf("PGPASSWORD=%s", dbPass),
		containerName,
		"psql", "-U", dbUser, "-d", dbName,
	)

	// Buka file yang baru saja diupload
	inFile, err := os.Open(filePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca file lokal untuk di-restore"})
	}
	defer inFile.Close()

	cmd.Stdin = inFile // Masukkan file ke Docker

	if err := cmd.Run(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "File terupload, tapi gagal di-restore ke database: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "File berhasil di-upload dan Database sukses di-restore!"})
}
