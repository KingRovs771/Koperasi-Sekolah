package seeder

import (
	"Kopsis-Spensa/internal/config"
	"Kopsis-Spensa/internal/models"
	"log"

	"github.com/google/uuid"
)

func SeedUsers() {
	var count int64

	config.DB.Model(&models.User{}).Count(&count)

	if count > 0 {
		log.Println("User already exists")
		return
	}

	users := []models.User{
		{
			UsersUID:    uuid.NewString(),
			NamaLengkap: "Rizky Budiarto",
			Username:    "kingrovs",
			Email:       "kingrovs@smpn1sragen.sch.id",
			Password:    "yolnDA2623*KingRovs771",
			NoHp:        "082142896072",
			Alamat:      "Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen",
			Role:        "administrator",
			Status:      "active",
		}, {
			UsersUID:    uuid.NewString(),
			NamaLengkap: "Mukhlis Royyani",
			Username:    "siroyy",
			Email:       "siroyy@smpn1sragen.sch.id",
			Password:    "coys162*",
			NoHp:        "08893884991",
			Alamat:      "Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen",
			Role:        "supervisor",
			Status:      "active",
		},
		{
			UsersUID:    uuid.NewString(),
			NamaLengkap: "Riyana Lili Lestari",
			Username:    "riyanalili",
			Email:       "riyanalili@smpn1sragen.sch.id",
			Password:    "garvi123*",
			NoHp:        "089577389893",
			Alamat:      "Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen",
			Role:        "supervisor",
			Status:      "active",
		},
		{
			UsersUID:    uuid.NewString(),
			NamaLengkap: "AdminSuper",
			Username:    "adminsuperspensa",
			Email:       "adminsuperspensa@smpn1sragen.sch.id",
			Password:    "smpn1sragen162*joss",
			NoHp:        "08779947729",
			Alamat:      "Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen",
			Role:        "administrator",
			Status:      "active",
		},
		{
			UsersUID:    uuid.NewString(),
			NamaLengkap: "cashier",
			Username:    "cashier",
			Email:       "cashier@smpn1sragen.sch.id",
			Password:    "smpn1sragen162*joss",
			NoHp:        "08779947729",
			Alamat:      "Jalan Raya Sukowati No. 162, Sragen Kulon, Sragen",
			Role:        "kasir",
			Status:      "active",
		},
	}
	for _, user := range users {
		if err := config.DB.Create(&user).Error; err != nil {
			log.Fatalf("❌ Gagal membuat user %s: %v", user.Email, err)
		}
		log.Printf("✅ User dibuat: %s (%s)", user.Email, user.Role)
	}

	log.Println("🌱 Seeder users berhasil dijalankan!")
}
