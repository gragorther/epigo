package gormdb

import (
	"context"

	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

func (g *GormDB) CheckUserAuthorizationForGroups(ctx context.Context, groupIDs []uint, userID uint) (bool, error) {

	count, err := gorm.G[models.Group](g.db).Where("user_id = ? AND id IN ?", userID, groupIDs).Count(ctx, "id")
	return int(count) == len(groupIDs), err
}
func (g *GormDB) CheckUserAuthorizationForLastMessage(ctx context.Context, messageID uint, userID uint) (bool, error) {
	count, err := gorm.G[models.LastMessage](g.db).Where("id = ? AND user_id = ?", messageID, userID).Count(ctx, "id")
	return count == 1, err
}
