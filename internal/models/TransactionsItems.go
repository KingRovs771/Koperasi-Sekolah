package models

import "time"

type TransactionsItems struct {
	TransactionsItemsID int64     `gorm:"primaryKey;AUTO_INCREMENT;uniqueIndex" json:"transactions_items_id"`
	TransactionsUID     string    `gorm:"varchar(255);" json:"transactions_uid"`
	ProductsUID         string    `gorm:"varchar(255);not null" json:"products_uid"`
	Jumlah              int64     `gorm:"not null" json:"jumlah"`
	HargaSaatIni        float64   `gorm:"decimal(15,2);not null" json:"harga_saat_ini"`
	ModalSaatIni        float64   `gorm:"decimal(15,2);not null" json:"modal_saat_ini"`
	SubTotal            float64   `gorm:"decimal(15,2);not null" json:"sub_total"`
	CreatedAt           time.Time `json:"created_at"`
}
