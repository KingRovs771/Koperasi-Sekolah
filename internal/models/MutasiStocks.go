package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MutasiStock struct {
	MutasiStocksID  int64     `gorm:"primary_key;AUTO_INCREMENT;UniqueIndex" json:"mutasi_stocks_id"`
	MutasiStocksUID string    `gorm:"varchar(255)" json:"mutasi_stocks_uid"`
	ProductsUID     string    `gorm:"varchar(255)" json:"products_uid"`
	UserUID         string    `gorm:"varchar(255)" json:"user_uid"`
	Tipe            string    `gorm:"varcher(20);not null" json:"tipe"`
	Jumlah          int64     `gorm:"int" json:"jumlah"`
	Catatan         string    `gorm:"text" json:"catatan"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (ms *MutasiStock) BeforeSave(tx *gorm.DB) (err error) {
	if ms.MutasiStocksUID == "" {
		ms.MutasiStocksUID = uuid.New().String()
	}

	ms.CreatedAt = time.Now()
	return
}
