package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"uniqueIndex;not null"`
}
