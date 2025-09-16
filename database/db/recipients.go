package db

import (
	"context"
)

type Recipient struct {
	Email string
}

func (d *DB) RecipientExistsByID(ctx context.Context, id uint) (exists bool, err error) {
	err = d.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM recipients WHERE id = $1)", id).Scan(&exists)
	return exists, err
}
