package db

import (
	"context"

	"github.com/samber/lo"
)

type Recipient struct {
	Email string
}

func RecipientsToStringArray(r []Recipient) []string {
	return lo.Map(r, func(item Recipient, _ int) string {
		return item.Email
	})
}

func (d *DB) RecipientExistsByID(ctx context.Context, id uint) (exists bool, err error) {
	err = d.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM recipients WHERE id = $1)", id).Scan(&exists)
	return exists, err
}
