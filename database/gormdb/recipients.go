package gormdb

import (
	"context"
	"fmt"

	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

func (g *GormDB) CheckIfRecipientExistsByID(ctx context.Context, id uint) (exists bool, err error) {
	count, err := gorm.G[models.Recipient](g.db).Where("id = ?", id).Count(ctx, "id")
	if err != nil {
		return false, fmt.Errorf("failed to check if recipient exists: %w", err)
	}
	return count == 1, nil
}
