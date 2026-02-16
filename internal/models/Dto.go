package models

import "time"

type CheckoutRequest struct {
	MetodePembayaran string  `json:"metode_pembayaran"`
	UangDiterima     float64 `json:"uang_diterima"`
	Items            []struct {
		ProductsUID string `json:"products_uid"`
		Quantity    int64  `json:"quantity"`
	} `json:"items"`
}

type MutasiHistoryDTO struct {
	Tipe       string    `json:"tipe"`
	Jumlah     int64     `json:"jumlah"`
	Catatan    string    `json:"catatan"`
	CreatedAt  time.Time `json:"created_at"`
	NamaProduk string    `json:"nama_produk"`
	NamaUser   string    `json:"nama_user"`
}
