package models

import (
	"time"
)

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement:true"`
	UserID      uint      `gorm:"user_id" json:"user_id"`
	Name        string    `gorm:"name" json:"name"`
	Description string    `gorm:"description" json:"description"`
	Price       float64   `gorm:"price" json:"price"`
	Quantity    int       `gorm:"quantity" json:"quantity"`
	Status      string    `gorm:"status" json:"status"`
	ImageURL    string    `gorm:"image_url" json:"image_url"`
	CreatedAt   time.Time `gorm:"created_at" json:"created_at"`
}
