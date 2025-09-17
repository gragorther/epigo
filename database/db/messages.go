package db

import (
	"context"

	_ "embed"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
)

//go:embed queries/can_user_edit_lastmessage.sql
var canUserEditLastMessageQuery string

func (d *DB) CanUserEditLastmessage(ctx context.Context, userID uint, messageID uint, groupIDs []uint) (authorized bool, err error) {
	err = d.db.QueryRow(ctx, canUserEditLastMessageQuery, messageID, userID, groupIDs).Scan(&authorized)
	return
}

type CreateLastMessage struct {
	UserID   uint
	Title    string
	Content  null.String
	GroupIDs []uint
}

func (d *DB) CreateLastMessage(ctx context.Context, message CreateLastMessage) error {
	_, err := d.db.Exec(ctx, `WITH m AS (INSERT INTO last_messages (title, content, user_id) VALUES ($1, $2, $3) RETURNING id)
		INSERT INTO group_last_messages (last_message_id, group_id) SELECT m.id, UNNEST($4::int[]) FROM m`, message.Title, message.Content, message.UserID, message.GroupIDs)
	return err
}

type LastMessage struct {
	Title   string
	Content null.String
	ID      uint
}

func (d *DB) LastMessagesByUserID(ctx context.Context, userID uint) (lastMessages []LastMessage, err error) {
	if err := pgxscan.Select(ctx, d.db, &lastMessages, "SELECT title, description, id FROM last_messages WHERE user_id = $1", userID); err != nil {
		return nil, err
	}
	return lastMessages, err
}

type LastMessageAndRecipients struct {
	LastMessage LastMessage
	Recipients  []Recipient
}

//go:embed queries/last_messages_and_recipients.sql
var lastMessagesAndRecipientsQuery string

func (d *DB) LastMessagesAndRecipients(ctx context.Context, userID uint) (lastMessages []LastMessageAndRecipients, err error) {
	rows, err := d.db.Query(ctx, lastMessagesAndRecipientsQuery, userID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (LastMessageAndRecipients, error) {
		var lm LastMessageAndRecipients
		var recipientEmails []string
		if err := row.Scan(&lm.LastMessage.Title, &lm.LastMessage.Content, &recipientEmails); err != nil {
			return LastMessageAndRecipients{}, err
		}
		lm.Recipients = lo.Map(recipientEmails, func(item string, _ int) (recipient Recipient) {
			recipient.Email = item
			return
		})
		return lm, nil
	})
}

func (d *DB) UpdateLastMessageTitle(ctx context.Context, id uint, title string) error {
	_, err := d.db.Exec(ctx, "UPDATE last_messages SET title = $1 WHERE id = $2", title, id)
	return err
}

type UpdateLastMessage struct {
	Title    null.String
	Content  null.String
	GroupIDs []uint
}

func (d *DB) UpdateLastMessage(ctx context.Context, id uint, m UpdateLastMessage) error {
	_, err := d.db.Exec(ctx, "UPDATE last_messages SET title = COALESCE($1, name), content = COALESCE($2, description) WHERE id = $3", m.Title, m.Content, id)
	return err
}

func (g *DB) DeleteLastMessageByID(ctx context.Context, id uint) error {
	_, err := g.db.Exec(ctx, "DELETE FROM last_messages WHERE id = $1", id)
	return err
}

func (d *DB) LastMessageExistsByID(ctx context.Context, id uint) (exists bool, err error) {
	err = d.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM last_messages WHERE id = $1)", id).Scan(&exists)
	return exists, err
}

type lastMessage struct {
	Title       string
	Description null.String
}

func (d *DB) LastMessageByID(ctx context.Context, id uint) (lastMessage lastMessage, err error) {
	err = d.db.QueryRow(ctx, "SELECT title, description FROM last_messages WHERE id = ?", id).Scan(&lastMessage.Title, &lastMessage.Description)
	return lastMessage, err
}
