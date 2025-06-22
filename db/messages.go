package db

import (
	"github.com/gragorther/epigo/models"
)

func (h *DBHandler) CreateLastMessage(lastMessage *models.LastMessage) error {
	res := h.DB.Create(&lastMessage)
	if res.Error != nil {
		return res.Error
	}
	return nil
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
type lastMessageOut struct {
	ID      uint   `gorm:"primarykey"`
	Title   string `json:"title"`
	Groups  []uint `json:"groups" gorm:"many2many:group_last_messages;"`
	Content string `json:"content"`
}

func (h *DBHandler) FindLastMessagesByUserID(userID uint) ([]lastMessageOut, error) {
	var lastMessages []lastMessage
	res := h.DB.Model(&models.LastMessage{}).Preload("Groups").Where("user_id = ?", userID).Find(&lastMessages)

	var lastMessagesOut []lastMessageOut
	for _, lastMessage := range lastMessages {
		var groups []uint
		for _, group := range lastMessage.Groups {
			groups = append(groups, group.ID)
		}

		lastMessagesOut = append(lastMessagesOut, lastMessageOut{
			ID:      lastMessage.ID,
			Title:   lastMessage.Title,
			Content: lastMessage.Content,
			Groups:  groups,
		})
	}
	return lastMessagesOut, res.Error
}

func (h *DBHandler) UpdateLastMessage(newMessage *models.LastMessage) error {
	res := h.DB.Model(&models.LastMessage{ID: newMessage.ID}).Updates(newMessage)
	h.DB.Model(&newMessage).Association("Groups").Replace(newMessage.Groups)
	return res.Error
}
func (h *DBHandler) DeleteLastMessageByID(lastMessageID uint) error {
	res := h.DB.Delete(&models.LastMessage{ID: lastMessageID})
	return res.Error
}
