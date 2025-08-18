package gormdb

import (
	"context"

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
	// the preload needs a builder function, we don't need any additional arguments so I just made it return nil
	return gorm.G[models.Group](g.db).Where("user_id = ?", userID).Preload("Recipients", func(db gorm.PreloadBuilder) error {
		return nil
	}).Find(ctx)

}
func (g *GormDB) FindGroupsAndLastMessagesByUserID(ctx context.Context, userID uint) ([]models.Group, error) {
	return gorm.G[models.Group](g.db).Where("user_id = ?", userID).Preload("LastMessages", func(db gorm.PreloadBuilder) error {
		return nil
	}).Find(ctx)
}
func (g *GormDB) FindGroupsAndLastMessagesAndRecipientsByUserID(ctx context.Context, userID uint) ([]models.Group, error) {
	return gorm.G[models.Group](g.db).Where("user_id = ?", userID).Preload("Recipients", func(db gorm.PreloadBuilder) error { return nil }).
		Preload("LastMessages", func(db gorm.PreloadBuilder) error { return nil }).Find(ctx)
}
func (g *GormDB) CreateGroup(group *models.Group) error {

	err := g.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(group).Error
		return err
	})
	return err
}
func (g *GormDB) CreateGroups(ctx context.Context, groups *[]models.Group) error {
	return gorm.G[models.Group](g.db).CreateInBatches(ctx, groups, 500)
}

func (g *GormDB) UpdateGroup(ctx context.Context, group models.Group) error {
	// make sure gorm doesn't create additional recipients instead of replacing them
	groupWithoutRecipients := group
	groupWithoutRecipients.Recipients = nil
	groupWithoutRecipients.LastMessages = nil

	return g.db.Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[models.Group](tx).Where("id = ?", group.ID).Updates(ctx, groupWithoutRecipients)
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Model(&group).Association("Recipients").Replace(group.Recipients)
		if err != nil {
			return err
		}
		err = tx.WithContext(ctx).Model(&group).Association("LastMessages").Replace(group.LastMessages)
		return err
	})

}

func (g *GormDB) CheckIfGroupExistsByID(ctx context.Context, groupID uint) (exists bool, err error) {
	count, err := gorm.G[models.Group](g.db).Where("id = ?", groupID).Count(ctx, "id")
	return count > 0, err

}
