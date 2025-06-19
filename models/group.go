package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	UserID          uint `json:"userID"`
	Name            string
	Description     string
	RecipientEmails []RecipientEmail
}

type RecipientEmail struct {
	gorm.Model
	GroupID uint   `json:"groupID"` // group of the email
	Email   string `json:"email"`
}
