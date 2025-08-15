package gormdb

import (
	"context"

	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

func (g *GormDB) CreateLastMessage(ctx context.Context, lastMessage *models.LastMessage) error {
	return gorm.G[models.LastMessage](g.db).Create(ctx, lastMessage)
}

type group struct {
	ID uint `gorm:"primarykey"`
}

func (g *GormDB) FindLastMessagesByUserID(userID uint) ([]models.LastMessage, error) {
	var lastMessages []models.LastMessage
	res := g.db.Model(&models.LastMessage{}).Preload("Groups").Where("user_id = ?", userID).Find(&lastMessages)

	return lastMessages, res.Error
}

func (g *GormDB) CreateLastMessages(ctx context.Context, lastMessages *[]models.LastMessage) error {
	// this wants me to specify the size of batches, so 200 it is i guess
	return gorm.G[models.LastMessage](g.db).CreateInBatches(ctx, lastMessages, 200)
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

func (g *GormDB) GetLastMessageByID(ctx context.Context, id uint) (models.LastMessage, error) {
	return gorm.G[models.LastMessage](g.db).Where("id = ?", id).First(ctx)
}
