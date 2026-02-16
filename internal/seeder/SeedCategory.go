package seeder

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"log"

	"github.com/google/uuid"
)

func SeedCategory() {
	var count int64

	config.DB.Find(&models.Category{}).Count(&count)

	if count > 0 {
		log.Println("Data Category Sudah Ada")
		return
	}

	category := []models.Category{
		{
			CategoryUID:  uuid.New().String(),
			NamaCategory: "Alat Tulis",
			Icon:         "Pencil",
			Description:  "Berisi Alat Tulis Seperti Pensil, Penggaris, Pulpen, Penghapus dan lain sebagainya",
			Status:       "active",
		},
		{
			CategoryUID:  uuid.New().String(),
			NamaCategory: "Seragam",
			Icon:         "Cloth",
			Description:  "Berisi Perlengkapan Seragam seperti topi, dasi, ikat pinggang dan lain sebagainya, ",
			Status:       "active",
		},
		{
			CategoryUID:  uuid.New().String(),
			NamaCategory: "Pramuka",
			Icon:         "Scout",
			Description:  "Berisi Alat Pramuka seperti semaphore, Peluit, Tali, dan lain sebagainya ",
			Status:       "active",
		},
	}
	for _, categorys := range category {
		if err := config.DB.Create(&categorys).Error; err != nil {
			log.Fatalf("❌ Gagal membuat Category %s: %v", categorys.NamaCategory, err)
		}
		log.Printf("✅ Category dibuat: %s (%s)", categorys.NamaCategory, categorys.Status)
	}

	log.Println("🌱 Seeder Category berhasil dijalankan!")
}
