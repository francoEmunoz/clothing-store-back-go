package models

import (
	"cs-go/db"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name      string `json:"name" gorm:"type:varchar(50);not null" validate:"required"`
	Category  string `json:"category" gorm:"type:varchar(50);not null" validate:"required"`
	Price     int    `json:"price" gorm:"not null" validate:"required"`
	Stock     int    `json:"stock" gorm:"not null" validate:"required"`
	Photo     string `json:"photo" gorm:"not null" validate:"required"`
	Questions []Question
}

type Products []Product

func MigrarProduct() {
	db.Database.AutoMigrate(Product{})
}
