package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	CategoryID   int64     `gorm:"primary_key;AUTO_INCREMENT;uniqueIndex" json:"category_id"`
	CategoryUID  string    `gorm:"varchar(255);unique" json:"category_uid"`
	NamaCategory string    `gorm:"varchar(100);not null" json:"nama_category"`
	Icon         string    `gorm:"size:50;default:'box'" json:"icon"`
	Description  string    `gorm:"varchar(255)" json:"description"`
	Status       string    `gorm:"enum('active','inactive');default:'active'"  json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (c *Category) BeforeSave(tx *gorm.DB) (err error) {
	if c.CategoryUID == "" {
		c.CategoryUID = uuid.New().String()
	}

	c.CreatedAt = time.Now()
	return
}
func (c *Category) BeforeUpdate(tx *gorm.DB) (err error) {
	if c.CategoryUID == "" {
		c.CategoryUID = uuid.New().String()
	}
	return
}
