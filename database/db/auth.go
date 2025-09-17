package db

import (
	"context"
)

func (d *DB) UserAuthorizationForGroups(ctx context.Context, groupIDs []uint, userID uint) (match bool, err error) {
	var count int
	err = d.db.QueryRow(ctx, "SELECT COUNT(*) FROM groups WHERE id = ANY($1::int[]) AND user_id = $2", groupIDs, userID).Scan(&count)
	return len(groupIDs) == count, err
}

func (d *DB) UserAuthorizationForLastMessage(ctx context.Context, messageID uint, userID uint) (authorized bool, err error) {
	err = d.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM last_messages WHERE id = $1 AND user_id = $2)", messageID, userID).Scan(&authorized)
	return authorized, err
}
