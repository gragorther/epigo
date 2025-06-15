package models

import "gorm.io/gorm"

type LastMessage struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"` // markdown
	UserID  uint   `json:"userID"`
	User    User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type SendToGroup struct {
	gorm.Model
	UserID      uint `json:"userID"`
	User        User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Name        string
	Description string
}

type RecipientEmail struct {
	gorm.Model
	UserID        uint        `json:"userID"`
	User          User        `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	SendToGroupID uint        `json:"sendToGroupID"`
	SendToGroup   SendToGroup `json:"sendToGroup"`
	Email         string      `json:"email"`
}

type SendToGroupInput struct {
	RecipientEmails []string `json:"recipients"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
}
