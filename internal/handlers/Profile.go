package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// 1. GET: Tampilkan Halaman Profil
func GetProfile(c *fiber.Ctx) error {
	// Ambil UID dari session (Pastikan key-nya sesuai dengan middleware Anda, misal "UsersUID")
	userVal := c.Locals("users_uid")
	if userVal == nil {
		return c.Status(401).SendString("Unauthorized")
	}
	userUID := userVal.(string)

	var user models.User
	if err := config.DB.Where("users_uid = ?", userUID).First(&user).Error; err != nil {
		return c.Status(404).SendString("User tidak ditemukan")
	}

	data := fiber.Map{
		"Title":       "Profil Pengguna",
		"User":        user,
		"Role":        c.Locals("Role"),
		"NamaLengkap": c.Locals("NamaLengkap"),
	}

	// Dukungan untuk HTMX
	if c.Get("HX-Request") == "true" {
		return c.Render("pages/profile", data)
	}
	return c.Render("pages/profile", data, "layouts/main")
}

// 2. POST: Update Informasi Profil (Data Diri & Akun)
func UpdateProfile(c *fiber.Ctx) error {
	userUID := c.Locals("users_uid").(string)

	type UpdateInput struct {
		NamaLengkap string `json:"nama_lengkap" form:"nama_lengkap"`
		Username    string `json:"username" form:"username"`
		Email       string `json:"email" form:"email"`
		NoHp        string `json:"no_hp" form:"no_hp"`
		Alamat      string `json:"alamat" form:"alamat"`
	}

	input := new(UpdateInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format input tidak valid"})
	}

	var user models.User
	if err := config.DB.Where("users_uid = ?", userUID).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Update field jika tidak kosong
	if input.NamaLengkap != "" {
		user.NamaLengkap = input.NamaLengkap
	}
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.NoHp != "" {
		user.NoHp = input.NoHp
	}
	if input.Alamat != "" {
		user.Alamat = input.Alamat
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan perubahan profil"})
	}

	return c.JSON(fiber.Map{"message": "Profil berhasil diperbarui!"})
}

// 3. POST: Update Password
func UpdatePassword(c *fiber.Ctx) error {
	userUID := c.Locals("UsersUID").(string)

	type PasswordInput struct {
		OldPassword string `json:"old_password" form:"old_password"`
		NewPassword string `json:"new_password" form:"new_password"`
	}

	input := new(PasswordInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	var user models.User
	if err := config.DB.Where("users_uid = ?", userUID).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Verifikasi password lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Password saat ini salah!"})
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memproses password baru"})
	}

	user.Password = string(hashedPassword)
	config.DB.Save(&user)

	return c.JSON(fiber.Map{"message": "Password berhasil diubah!"})
}

// 4. POST: Upload Foto Profil
func UploadPhoto(c *fiber.Ctx) error {
	userUID := c.Locals("UsersUID").(string)

	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal membaca file upload"})
	}

	// Buat folder jika belum ada
	uploadDir := "./uploads/profiles"
	os.MkdirAll(uploadDir, os.ModePerm)

	// Generate nama file unik
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan foto ke server"})
	}

	// Update DB
	var user models.User
	if err := config.DB.Where("users_uid = ?", userUID).First(&user).Error; err == nil {
		// Hapus foto lama jika ada
		if user.Foto != "" {
			os.Remove(filepath.Join(uploadDir, user.Foto))
		}
		user.Foto = filename
		config.DB.Save(&user)
	}

	return c.JSON(fiber.Map{"message": "Foto profil berhasil diperbarui!"})
}
