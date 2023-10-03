package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
