package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	UserID          uint             `json:"userID"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	RecipientEmails []RecipientEmail `json:"recipientEmails"`
}

type RecipientEmail struct {
	gorm.Model
	GroupID uint   `json:"groupID"` // group of the email
	Email   string `json:"email"`
}
