package handlers

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"

	"github.com/gofiber/fiber/v2"
)

func GetCategories(c *fiber.Ctx) error {

	var categories []models.Category
	var totalProduk int64

	if err := config.DB.Order("created_at desc").Find(&categories).Error; err != nil {
		return c.Status(500).SendString("Gagal memuat kategori")
	}

	var catAktif, catNonaktif int
	for _, cat := range categories {
		if cat.Status == "active" {
			catAktif++
		} else {
			catNonaktif++
		}
	}

	if err := config.DB.Model(&models.Product{}).Count(&totalProduk).Error; err != nil {
		totalProduk = 0
	}

	// 4. Siapkan Data Map
	data := fiber.Map{
		"Title":            "Manajemen Kategori",
		"Data":             categories,
		"TotalKategori":    len(categories),
		"KategoriAktif":    catAktif,
		"KategoriNonaktif": catNonaktif,
		"TotalProduk":      totalProduk,
		"Role":             c.Locals("Role"),
		"Nama_lengkap":     c.Locals("NamaLengkap"),
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("pages/manage_category", data)
	}

	return c.Render("pages/manage_category", data, "layouts/main")
}

func CreateCategory(c *fiber.Ctx) error {
	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Error": err,
		})
	}

	if err := config.DB.Create(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Error": err,
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"Message": "Category Created Successfully",
		"Data":    category,
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	type UpdateInputCategory struct {
		NameCategory string `json:"name_category"`
		Icon         string `json:"icon"`
		Description  string `json:"description"`
		Status       string `json:"status"`
	}

	uid := c.Params("uid")
	var category models.Category

	if err := config.DB.Where("uid = ?", uid).First(&category).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"Error": err,
		})
	}

	input := new(UpdateInputCategory)
	if err := c.BodyParser(input); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Error": err,
		})
	}

	category.NamaCategory = input.NameCategory
	category.Icon = input.Icon
	category.Description = input.Description
	category.Status = input.Status

	config.DB.Save(&category)

	return c.JSON(fiber.Map{
		"Message": "Category Updated Successfully",
		"Data":    category,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	uid := c.Params("uid")
	var category models.Category

	if err := config.DB.Where("uid = ?", uid).First(&category).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"Error":   err,
			"Message": "Category Not Found",
		})
	}

	return c.JSON(fiber.Map{
		"Message": "Category Deleted Successfully",
		"Data":    category,
	})
}
