package models

import "gorm.io/gorm"

type LastMessage struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"` // markdown
	UserID  uint   `json:"userID"`
	User    User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
