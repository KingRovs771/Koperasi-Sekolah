package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"Kopsis-Spensa/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	var users []models.User

	if err := config.DB.Order("created_at desc").Find(&users).Error; err != nil {
		return c.Status(500).SendString("Gagal memuat data pengguna")
	}

	var totalAdmin, totalKasir, totalSpv int

	for _, u := range users {
		switch u.Role {
		case "administrator":
			totalAdmin++
		case "kasir":
			totalKasir++
		case "supervisor":
			totalSpv++
		}
	}

	data := fiber.Map{
		"Title":        "Manajemen Pengguna",
		"Data":         users,
		"TotalAdmin":   totalAdmin,
		"TotalKasir":   totalKasir,
		"TotalSpv":     totalSpv,
		"Role":         c.Locals("Role"),
		"Nama_Lengkap": c.Locals("NamaLengkap"),
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("pages/manage_users", data)
	}

	return c.Render("pages/manage_users", data, "layouts/main")
}

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Error": err,
		})
	}

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = hashed

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan user (Username/Email mungkin duplikat)"})
	}

	return c.Status(200).JSON(fiber.Map{
		"Success": true,
		"Data":    user,
	})
}

func GetUserByUID(c *fiber.Ctx) error {
	uid := c.Params("uid")
	var user models.User

	if err := config.DB.Where("users_uid = ?", uid).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"Error": err,
		})
	}

	user.Password = ""
	return c.Status(200).JSON(fiber.Map{
		"Success": true,
		"Data":    user,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	type UpdateInput struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		IsActive string `json:"is_active"`
	}

	input := new(UpdateInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	user.NamaLengkap = input.Name
	user.Username = input.Username
	user.Email = input.Email
	user.Role = input.Role
	user.Status = input.IsActive

	if input.Password != "" {
		hashed, _ := utils.HashPassword(input.Password)
		user.Password = hashed
	}

	config.DB.Save(&user)

	return c.JSON(fiber.Map{"message": "User berhasil diperbarui", "user": user})
}

func DeleteUser(c *fiber.Ctx) error {
	uid := c.Params("id")
	var user models.User

	if err := config.DB.Where("users_uid = ?", uid).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"Error": "User tidak ditemukan",
		})
	}
	config.DB.Delete(&user)
	return c.JSON(fiber.Map{"message": "User berhasil dihapuskan"})
}
