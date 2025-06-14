package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"unique"`
	PasswordHash string `json:"passwordHash"`
	//	Country      string //should probably be a foreign key of another table
	LastLogin *time.Time `json:"lastLogin"`
}
type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
