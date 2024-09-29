package models

import (
	"time"
)

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
	ImageURL    []string  `json:"image_url" gorm:"type:text[]"`
	CreatedAt   time.Time `json:"created_at"`
}
