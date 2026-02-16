package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ProductsID   int            `gorm:"primaryKey;uniqueIndex" json:"products_id"`
	ProductsUID  string         `gorm:"varchar(255);uniqueIndex" json:"products_uid"`
	KodeProduk   string         `gorm:"size:50;unique" json:"kode_produk"`
	CategoryUID  string         `gorm:"varchar(255);uniqueIndex" json:"category_uid"`
	NamaProducts string         `gorm:"varchar(100)" json:"nama_products"`
	Description  string         `gorm:"varchar(255)" json:"description"`
	HargaJual    int64          `gorm:"int" json:"harga_jual"`
	HargaBeli    int64          `gorm:"int" json:"harga_beli"`
	Stock        int64          `gorm:"int" json:"stock"`
	StockMin     int64          `gorm:"int" json:"stock_min"`
	Image        string         `gorm:"varchar(255)" json:"image"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}
