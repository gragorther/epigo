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
	Groups    Group
}
