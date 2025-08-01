package gormdb

import (
	"fmt"
	"time"

	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

func (g *GormDB) DeleteGroupByID(id uint) error {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		//res := tx.Delete(&models.Group{}, id)
		group := models.Group{ID: id}
		res := tx.Select("RecipientEmails").Delete(&group)
		return res.Error
	})
	return err
}

type recipient struct {
	Email string `json:"email"`
}

type listGroupsDTO struct {
	ID          uint        `gorm:"primarykey"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Recipients  []recipient `json:"recipients" gorm:"foreignKey:GroupID"`
}

func (g *GormDB) FindGroupsAndRecipientsByUserID(userID uint) ([]models.Group, error) {
	var groups []models.Group
	res := g.db.Select("id", "name", "description", "recipients").Where("user_id = ?", userID).Preload("Recipients").Find(&groups)
	if res.Error != nil {
		return nil, res.Error
	}

	return groups, nil
}

func (g *GormDB) CreateGroupAndRecipientEmails(group *models.Group) error {

	err := g.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(group).Error
		return err
	})
	return err
}

func (g *GormDB) UpdateGroup(group *models.Group) error {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		output := tx.Updates(group)
		if output.Error != nil {
			return output.Error
		}
		if output.RowsAffected < 1 {
			return fmt.Errorf("failed to update group: less than 1 rows affected")
		}
		err := tx.Model(group).Association("Recipients").Replace(group.Recipients)

		return err
	})
	return err

}
