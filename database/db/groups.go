package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	_ "embed"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/guregu/null/v6"
)

//go:embed queries/can_user_edit_group.sql
var canUserEditGroupQuery string

func (d *DB) CanUserEditGroup(ctx context.Context, userID uint, groupID uint, lastMessageIDs []uint) (authorized bool, err error) {
	err = d.db.QueryRow(ctx, canUserEditGroupQuery, groupID, userID, lastMessageIDs).Scan(&authorized)
	return
}

func (d *DB) DeleteGroupByID(ctx context.Context, id uint) error {
	_, err := d.db.Exec(ctx, "DELETE FROM groups WHERE id = $1", id)
	return err
}

type Group struct {
	Name        string
	Description null.String
	ID          uint
}

func (d *DB) GroupsByUserID(ctx context.Context, userID uint) (groups []Group, err error) {
	if err := pgxscan.Select(ctx, d.db, &groups, "SELECT name, description, id FROM groups WHERE user_id = $1", userID); err != nil {
		return nil, err
	}
	return groups, err
}

type CreateGroup struct {
	Name           string
	Description    null.String
	UserID         uint
	LastMessageIDs []uint
}

//go:embed queries/create_group.sql
var createGroupQuery string

func (d *DB) CreateGroup(ctx context.Context, group CreateGroup) error {
	_, err := d.db.Exec(ctx, createGroupQuery, group.Name, group.Description, group.UserID, group.LastMessageIDs)
	return err
}

type CreateGroupsForUser struct {
	Name           string
	Description    null.String
	LastMessageIDs []uint
}

// doesn't work yet
func (d *DB) CreateGroupsForUser(ctx context.Context, userID uint, groups []CreateGroupsForUser) error {
	if len(groups) == 0 {
		return nil
	}
	args := make([]any, 0, 1+len(groups)*3)
	args = append(args, userID) // $1

	valueLines := make([]string, 0, len(groups))
	placeholder := uint(2)
	for i, group := range groups {
		ord := i + 1
		namePos := placeholder
		descPos := placeholder + 1
		arrPos := placeholder + 2
		placeholder += 3

		valueLines = append(valueLines, fmt.Sprintf("(%d, $%d, $%d, $%d::int[])", ord, namePos, descPos, arrPos))

		// push args in same order as placeholders
		args = append(args, group.Name)
		args = append(args, group.Description) // null.String works with pgx
		args = append(args, group.LastMessageIDs)
	}

	sql := fmt.Sprintf(`
WITH data (ord, name, description, last_ids) AS (
  VALUES
    %s
),
ins AS (
  INSERT INTO groups (user_id, name, description)
  SELECT $1, name, description FROM data ORDER BY ord
  RETURNING id
),
ins_ord AS (
  SELECT id, row_number() OVER () AS ord FROM ins
)
INSERT INTO group_last_messages (group_id, last_message_id)
SELECT ins_ord.id, unnest(data.last_ids)
FROM ins_ord
JOIN data ON ins_ord.ord = data.ord;
`, strings.Join(valueLines, ",\n    "))

	d.db.Exec(ctx, sql, args...)
	return nil
}

//go:embed queries/create_group_returning_id.sql
var createGroupReturningIDQuery string

func (d *DB) CreateGroupReturningID(ctx context.Context, group CreateGroup) (groupID uint, err error) {
	err = d.db.QueryRow(ctx, createGroupReturningIDQuery, group.Name, group.Description, group.UserID, group.LastMessageIDs).Scan(&groupID)
	return
}

func (d *DB) UpdateGroupDescription(ctx context.Context, id uint, newDescription string) error {
	_, err := d.db.Exec(ctx, "UPDATE groups SET description = $1 WHERE id = $2", newDescription, id)
	return err
}

func (d *DB) UpdateGroupName(ctx context.Context, id uint, newName string) error {
	_, err := d.db.Exec(ctx, "UPDATE groups SET name = $1 WHERE id = $2", newName, id)
	return err
}

type UpdateGroup struct {
	Name        null.String
	Description null.String
}

var ErrAllFieldsEmpty = errors.New("all the fields in the struct are empty")

func (d *DB) UpdateGroup(ctx context.Context, id uint, group UpdateGroup) error {
	_, err := d.db.Exec(ctx, "UPDATE groups SET name = COALESCE($1, name), description = COALESCE($2, description) WHERE id = $3", group.Name, group.Description, id)
	return err
}

func (d *DB) GroupExistsByID(ctx context.Context, groupID uint) (exists bool, err error) {
	err = d.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM groups WHERE id = $1)", groupID).Scan(&exists)
	return exists, err
}
