package db

import (
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/types"
)

type Messages interface {
	CreateLastMessage(lastMessage *models.LastMessage) error
	FindLastMessagesByUserID(userID uint) ([]types.LastMessageOut, error)
	UpdateLastMessage(newMessage *models.LastMessage) error
	DeleteLastMessageByID(lastMessageID uint) error
}
