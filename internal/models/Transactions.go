package models

import "time"

type Transactions struct {
	TransactionsID   int64     `gorm:"primary_key;AUTO_INCREMENT;UniqueIndex" json:"transactions_id"`
	TransactionsUID  string    `gorm:"varchar(255);uniqueIndex" json:"transactions_uid"`
	NoFaktur         string    `gorm:"varchar(50);not null" json:"no_faktur"`
	UserUID          string    `gorm:"varchar(255);;not null" json:"user_uid"`
	TotalBelanja     float64   `gorm:"decimal(15,2);not null" json:"total_belanja"`
	TotalModal       float64   `gorm:"decimal(15,2);not null" json:"total_modal"`
	MetodePembayaran string    `gorm:"enum('tunai','qris','transfer','tabungan');default:'tunai'" json:"metode_pembayaran"`
	UangDiterima     float64   `gorm:"decimal(15,2)" json:"uang_diterima"`
	Kembalian        float64   `gorm:"decimal(15,2)" json:"kembalian"`
	CreatedAt        time.Time `json:"created_at"`
}
