package db

import (
	"context"
	"errors"

	_ "embed"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5"
)

func (d *DB) UpdateUserInterval(ctx context.Context, userID uint, cron string) error {
	_, err := d.db.Exec(ctx, "UPDATE users SET cron = $1 WHERE id = $2", cron, userID)
	return err
}

type UserInterval struct {
	ID    uint   `gorm:"primarykey"`
	Email string `json:"email" gorm:"unique"`
	Cron  string `json:"cron"`
	Name  string `json:"name"`
}

type UserSentEmails struct {
	SentEmails    uint `json:"sentEmails"`
	MaxSentEmails uint `json:"maxSentEmails"`
}

func (d *DB) GetUserSentEmails(ctx context.Context, userID uint) (sentEmails UserSentEmails, err error) {
	err = d.db.QueryRow(ctx, "SELECT sent_emails, max_sent_emails FROM users WHERE id = $1", userID).Scan(&sentEmails.SentEmails, &sentEmails.MaxSentEmails)
	return sentEmails, err
}

func (d *DB) GetUserIntervals(ctx context.Context) (userIntervals []UserInterval, err error) {
	err = pgxscan.Select(ctx, d.db, &userIntervals, "SELECT id, email, cron, name FROM users")
	return userIntervals, err
}

type IntervalAndSentEmails struct {
	UserSentEmails
	UserInterval
}

func (d *DB) AllUserIntervalsAndSentEmails(ctx context.Context) (intervals []IntervalAndSentEmails, err error) {
	rows, err := d.db.Query(ctx, "SELECT sent_emails, max_sent_emails, id, email, cron, name FROM users")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (interval IntervalAndSentEmails, err error) {
		err = row.Scan(&interval.SentEmails, &interval.MaxSentEmails, &interval.ID, &interval.Email, &interval.Cron, &interval.Name)
		return
	})
}

// true if user exists, false if they don't exist
func (d *DB) CheckIfUserExistsByUsernameAndEmail(ctx context.Context, username string, email string) (exists bool, err error) {
	err = d.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)", username, email).Scan(&exists)
	return exists, err
}

func (d *DB) CheckIfUserExistsByUsername(ctx context.Context, username string) (exists bool, err error) {
	err = d.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	return exists, err
}

func (d *DB) CheckIfUserExistsByEmail(ctx context.Context, email string) (exists bool, err error) {
	err = d.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}

type CreateUserInput struct {
	Username     string
	Email        string
	PasswordHash string
	Name         null.String
}

func (d *DB) CreateUser(ctx context.Context, user CreateUserInput) error {
	_, err := d.db.Exec(ctx, "INSERT INTO users (username, email, name, password_hash) VALUES ($1, $2, $3, $4)", user.Username, user.Email, user.Name, user.PasswordHash)
	return err
}

func (d *DB) CreateUserReturningID(ctx context.Context, user CreateUserInput) (userID uint, err error) {
	err = d.db.QueryRow(ctx, "INSERT INTO users (username, email, name, password_hash) VALUES ($1, $2, $3, $4,) RETURNING id", user.Username, user.Email, user.Name, user.PasswordHash).Scan(&userID)
	return
}

type UserIDAndPasswordHash struct {
	ID           uint
	PasswordHash string
}

func (d *DB) UserIDAndPasswordHashByUsername(ctx context.Context, username string) (user UserIDAndPasswordHash, err error) {
	err = d.db.QueryRow(ctx, "SELECT password_hash, id FROM users WHERE username = $1", username).Scan(&user.PasswordHash, &user.ID)
	return user, err
}

func (d *DB) UserExistsByID(ctx context.Context, ID uint) (exists bool, err error) {
	err = d.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", ID).Scan(&exists)
	return exists, err
}

type User struct {
	Username string
	Name     null.String
	Email    string
}

func (d *DB) UserByID(ctx context.Context, ID uint) (user User, err error) {
	err = d.db.QueryRow(ctx, "SELECT username, name, email FROM users WHERE id = $1", ID).Scan(&user.Username, &user.Name, &user.Email)
	return user, err
}

type UserWithCron struct {
	User
	Cron null.String
}

func (d *DB) UserWithCronByID(ctx context.Context, userID uint) (user UserWithCron, err error) {
	err = d.db.QueryRow(ctx, "SELECT username, name, email, cron FROM users WHERE id = $1", userID).Scan(&user.Username, &user.Name, &user.Email, &user.Cron)
	return
}

func (d *DB) DeleteUser(ctx context.Context, ID uint) error {
	_, err := d.db.Exec(ctx, "DELETE FROM users WHERE id = $1", ID)
	return err
}

func (d *DB) SetUserMaxSentEmails(ctx context.Context, userID uint, maxSentEmails uint) error {
	_, err := d.db.Exec(ctx, "UPDATE users SET max_sent_emails = $1 WHERE id = $2", maxSentEmails, userID)
	return err
}

func (d *DB) SetUserName(ctx context.Context, userID uint, name string) error {
	_, err := d.db.Exec(ctx, "UPDATE users SET name = $1 WHERE id = $2", name, userID)
	return err
}

var ErrNoRowsAffected error = errors.New("no rows affected")

func (d *DB) IncrementUserSentEmailsCount(ctx context.Context, userID uint) error {
	_, err := d.db.Exec(ctx, "UPDATE users SET sent_emails = sent_emails + 1 WHERE id = $1", userID)
	return err
}
