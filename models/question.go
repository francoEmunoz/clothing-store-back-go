package models

import (
	"cs-go/db"

	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	Content   string `json:"content" gorm:"type:varchar(250);not null" validate:"required"`
	UserID    uint   `json:"user_id" validate:"required"`
	User      User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ProductID uint   `json:"product_id" validate:"required"`
}

type Questions []Question

func MigrarQuestion() {
	db.Database.AutoMigrate(Question{})
}
