package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func ShowLogin(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)
	if err == nil && sess.Get("user_uid") != nil {
		return c.Redirect("/dashboard")
	}
	return c.Render("pages/login", fiber.Map{
		"title": "Login Sistem POS Koperasi",
	})
}

func ProccessLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	var user models.User

	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.Render("pages/login", fiber.Map{
			"Error": "Username tidak ditemukan Silakan Login Kembali",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println(err.Error())
		return c.Render("pages/login", fiber.Map{
			"Error":   "Password Salah Silakan Login Kembali",
			"Message": err.Error(),
		})
	}

	if user.Status != "active" {
		return c.Render("pages/login", fiber.Map{
			"Error": "Akun Anda Dinonaktifkan Silakan Hubungi Administrator",
		})
	}

	sess, err := config.Store.Get(c)
	if err != nil {
		log.Printf("❌ DEBUG: Gagal Get Session dari Store: %v\n", err)
		return c.Status(500).SendString("Gagal Koneksi Sesi :" + err.Error())
	}
	sess.Set("users_uid", user.UsersUID)
	sess.Set("Username", user.Username)
	sess.Set("Role", user.Role)
	sess.Set("NamaLengkap", user.NamaLengkap)

	if err := sess.Save(); err != nil {
		log.Printf("❌ DEBUG: Gagal Save Session ke Redis: %v\n", err)
		return c.Status(500).SendString("Gagal Menyimpan Sesi: " + err.Error())
	}

	if user.Role == "kasir" {
		return c.Redirect("/pos")
	}

	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)

	if err != nil {
		return c.Redirect("/login")
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(500).SendString("Gagal Logout")
	}
	return c.Redirect("/login")
}
