package gormdb

import (
	"context"
	"fmt"

	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

func (g *GormDB) DeleteGroupByID(ctx context.Context, id uint) error {
	return g.db.Transaction(func(tx *gorm.DB) error {

		err := gorm.G[models.Recipient](tx).Exec(ctx, "DELETE FROM recipients WHERE group_id = ?", id)
		if err != nil {
			return err
		}
		err = gorm.G[models.Group](tx).Exec(ctx, "DELETE FROM groups WHERE id = ?", id)

		return err
	})

}

func (g *GormDB) FindGroupsAndRecipientsByUserID(ctx context.Context, userID uint) ([]models.Group, error) {
	var groups []models.Group
	res := g.db.Select("id", "name", "description", "recipients").Where("user_id = ?", userID).Preload("Recipients").Find(&groups)
	if res.Error != nil {
		return nil, res.Error
	}

	return groups, nil
}

func (g *GormDB) CreateGroup(group *models.Group) error {

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

func (g *GormDB) CheckIfGroupExistsByID(ctx context.Context, groupID uint) (exists bool, err error) {
	count, err := gorm.G[models.Group](g.db).Where("id = ?", groupID).Count(ctx, "id")
	return count > 0, err

}
