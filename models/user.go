package models

import "time"

type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"username" gorm:"unique"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	//	Country      string //should probably be a foreign key of another table
	LastLogin *time.Time `json:"lastLogin"`
}
type AuthInput struct { // prevents client from modifying everything in the users table
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
