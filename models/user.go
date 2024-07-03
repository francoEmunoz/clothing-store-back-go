package models

import (
	"cs-go/db"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string `json:"name" gorm:"type:varchar(50);not null" validate:"required"`
	Lastname      string `json:"lastname" gorm:"type:varchar(50);not null" validate:"required"`
	Password      string `json:"password" gorm:"type:varchar(100);not null" validate:"required"`
	PlainPassword string `json:"-"`
	Email         string `json:"email" gorm:"not null" validate:"required,email"`
	Country       string `json:"country" gorm:"not null" validate:"required"`
	Role          string `json:"role" gorm:"not null" validate:"required"`
	Logged        bool   `json:"logged" gorm:"default:false"`
}

type UserRegistration struct {
	Name          string `json:"name" validate:"required"`
	Lastname      string `json:"lastname" validate:"required"`
	PlainPassword string `json:"plainPassword" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Country       string `json:"country" validate:"required"`
	Role          string `json:"role" validate:"required"`
}

type UserLogin struct {
	Email         string `json:"email" validate:"required,email"`
	PlainPassword string `json:"plainPassword" validate:"required"`
}

type Users []User

func MigrarUser() {
	db.Database.AutoMigrate(User{})
}
