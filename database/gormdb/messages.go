package gormdb

import (
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

func (g *GormDB) CreateLastMessage(lastMessage *models.LastMessage) error {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&lastMessage)
		return res.Error
	})
	return err
}

type group struct {
	ID uint `gorm:"primarykey"`
}

func (g *GormDB) FindLastMessagesByUserID(userID uint) ([]models.LastMessage, error) {
	var lastMessages []models.LastMessage
	res := g.db.Model(&models.LastMessage{}).Preload("Groups").Where("user_id = ?", userID).Find(&lastMessages)

	return lastMessages, res.Error
}

func (g *GormDB) UpdateLastMessage(newMessage *models.LastMessage) error {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&models.LastMessage{ID: newMessage.ID}).Updates(newMessage)
		if res.Error != nil {
			return res.Error
		}
		err := tx.Model(&newMessage).Association("Groups").Replace(newMessage.Groups)
		return err
	})
	return err
}
func (g *GormDB) DeleteLastMessageByID(lastMessageID uint) error {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&models.LastMessage{ID: lastMessageID})
		return res.Error
	})
	return err
}
