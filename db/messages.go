package db

import "github.com/gragorther/epigo/models"

type Messages interface {
	CreateLastMessage(lastMessage *models.LastMessage) error
	FindLastMessagesByUserID(userID uint) ([]LastMessageOut, error)
	UpdateLastMessage(newMessage *models.LastMessage) error
	DeleteLastMessageByID(lastMessageID uint) error
}
type LastMessageOut struct {
	ID      uint   `gorm:"primarykey"`
	Title   string `json:"title"`
	Groups  []uint `json:"groups" gorm:"many2many:group_last_messages;"`
	Content string `json:"content"`
}
