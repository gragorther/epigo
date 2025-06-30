package gormdb

import (
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/types"
	"gorm.io/gorm"
)

type MessageDB struct {
	DB *gorm.DB
}

func (m *MessageDB) CreateLastMessage(lastMessage *models.LastMessage) error {
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&lastMessage)
		return res.Error
	})
	return err
}

type group struct {
	ID uint `gorm:"primarykey"`
}

type lastMessage struct {
	ID      uint    `gorm:"primarykey"`
	Title   string  `json:"title"`
	Groups  []group `json:"groups" gorm:"many2many:group_last_messages;"`
	Content string  `json:"content"`
}

func (m *MessageDB) FindLastMessagesByUserID(userID uint) ([]types.LastMessageOut, error) {
	var lastMessages []lastMessage
	res := m.DB.Model(&models.LastMessage{}).Preload("Groups").Where("user_id = ?", userID).Find(&lastMessages)

	var lastMessagesOut []types.LastMessageOut
	for _, lastMessage := range lastMessages {
		var groups []uint
		for _, group := range lastMessage.Groups {
			groups = append(groups, group.ID)
		}

		lastMessagesOut = append(lastMessagesOut, types.LastMessageOut{
			ID:      lastMessage.ID,
			Title:   lastMessage.Title,
			Content: lastMessage.Content,
			Groups:  groups,
		})
	}
	return lastMessagesOut, res.Error
}

func (m *MessageDB) UpdateLastMessage(newMessage *models.LastMessage) error {
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&models.LastMessage{ID: newMessage.ID}).Updates(newMessage)
		if res.Error != nil {
			return res.Error
		}
		err := tx.Model(&newMessage).Association("Groups").Replace(newMessage.Groups)
		return err
	})
	return err
}
func (m *MessageDB) DeleteLastMessageByID(lastMessageID uint) error {
	err := m.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&models.LastMessage{ID: lastMessageID})
		return res.Error
	})
	return err
}
