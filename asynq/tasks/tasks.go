package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gragorther/epigo/database/db"
	"github.com/gragorther/epigo/email"
	"github.com/gragorther/epigo/tokens"
	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeRecurringEmail = "email:recurring"
)

type TaskEnqueuer func(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)

func EnqueueTask(client *asynq.Client) TaskEnqueuer {
	return func(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
		return client.Enqueue(task, opts...)
	}
}

// Task payload for any email related tasks.
type RecurringEmailTaskPayload struct {
	// ID for the email recipient.
	UserID       uint
	Name         string
	Email        string
	ExpiresAfter time.Duration
}
type verificationEmailSender interface {
	SendVerificationEmail(ctx context.Context, user email.User, registrationLink string) error
}

func NewRecurringEmailTask(userID uint, name string, email string, expiresAfter time.Duration) (*asynq.Task, error) {
	payload, err := sonic.Marshal(RecurringEmailTaskPayload{Name: name, Email: email, ExpiresAfter: expiresAfter, UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRecurringEmail, payload), nil
}

// verificationURL is the URL that takes a token parameter, e.g. https://afterwill.life/user/life/verify?token=loremipsumdolorsitamet
func HandleRecurringEmailTask(emailService interface {
	SendUserLifeStatusEmail(ctx context.Context, user email.LifeStatusUser, verificationURL string) error
}, db interface {
	IncrementUserSentEmailsCount(ctx context.Context, userID uint) error
	GetUserSentEmails(context.Context, uint) (db.UserSentEmails, error)
}, createUserLifeStatusToken tokens.CreateUserLifeStatusFunc, verificationURL string,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p RecurringEmailTaskPayload
		if err := sonic.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}

		token, err := createUserLifeStatusToken(p.UserID, time.Now().Add(p.ExpiresAfter))
		if err != nil {
			return err
		}

		if err := emailService.SendUserLifeStatusEmail(ctx, email.LifeStatusUser{Name: p.Name, Email: p.Email}, fmt.Sprintf("%s?%s", verificationURL, token)); err != nil {
			return err
		}

		return db.IncrementUserSentEmailsCount(ctx, p.UserID)
	}
}
