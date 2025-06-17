package models

import "gorm.io/gorm"

type SendToGroup struct {
	gorm.Model
	UserID          uint `json:"userID"`
	Name            string
	Description     string
	RecipientEmails []RecipientEmail
}

type RecipientEmail struct {
	gorm.Model
	SendToGroupID uint   `json:"sendToGroupID"` // group of the email
	Email         string `json:"email"`
}

type SendToGroupInput struct {
	RecipientEmails []string `json:"recipientEmails"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
}
