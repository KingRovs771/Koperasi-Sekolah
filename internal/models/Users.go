package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	UsersID     int64          `gorm:"primary_key;AUTO_INCREMENT;uniqueIndex" json:"users_id"`
	UsersUID    string         `gorm:"varchar(255);uniqueIndex" json:"users_uid"`
	NamaLengkap string         `gorm:"varchar(100);not null" json:"nama_lengkap"`
	Username    string         `gorm:"varchar(100);uniques" json:"username"`
	Email       string         `gorm:"varchar(100)" json:"email"`
	Password    string         `gorm:"varchar(255)" json:"password"`
	NoHp        string         `gorm:"varchar(20)" json:"no_hp"`
	Alamat      string         `gorm:"text" json:"alamat"`
	Role        string         `gorm:"enum('administrator', 'kasir', 'supervisor');default:'kasir'"  json:"role"`
	Status      string         `gorm:"enum('active','inactive');default:'active'"  json:"status"`
	Foto        string         `gorm:"varchar(255)" json:"foto"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdateAt    time.Time      `json:"update_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.UsersUID == "" {
		u.UsersUID = uuid.New().String()
	}
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hash)
	}
	return
}
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hash)
	}
	return
}
