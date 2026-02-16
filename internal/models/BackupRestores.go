package models

import "time"

type BackupRestore struct {
	BackupRestoreID  int64     `gorm:"primary_key;AUTO_INCREMENT;uniqueIndex'" json:"backup_restore_id"`
	BackupRestoreUID string    `gorm:"varchar(255)" json:"backup_restore_uid"`
	NamaFile         string    `gorm:"varchar(100)" json:"nama_file"`
	UkuranFile       int64     `gorm:"int" json:"ukuran_file"`
	CreatedAt        time.Time `json:"created_at"`
}
