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

type queue struct {
	enqueueTask TaskEnqueueFunc
	marshal     MarshalFunc
}

func NewQueue(
	enqueueTask TaskEnqueueFunc,
	marshal MarshalFunc,
) *queue {
	return &queue{
		enqueueTask: enqueueTask,
		marshal:     marshal,
	}
}

type (
	UnmarshalFunc func(data []byte, v any) error
	MarshalFunc   func(val any) ([]byte, error)
)

// A list of task types.
const (
	TypeRecurringEmail = "email:recurring"
)

type TaskEnqueueFunc func(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)

type Enqueuer interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

func EnqueueTask(client Enqueuer) TaskEnqueueFunc {
	return func(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
		return client.Enqueue(task, opts...)
	}
}

// Task payload for any email related tasks.
type recurringEmailTaskPayload struct {
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
	payload, err := sonic.Marshal(recurringEmailTaskPayload{Name: name, Email: email, ExpiresAfter: expiresAfter, UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRecurringEmail, payload), nil
}

// verificationURL is the URL that takes a token parameter, e.g. https://afterwill.life/user/life/verify?token=loremipsumdolorsitamet
func HandleRecurringEmail(emailService interface {
	SendUserLifeStatusEmail(ctx context.Context, user email.LifeStatusUser, verificationURL string) error
}, db interface {
	IncrementUserSentEmailsCount(ctx context.Context, userID uint) error
	GetUserSentEmails(context.Context, uint) (db.UserSentEmails, error)
}, unmarshal UnmarshalFunc, createUserLifeStatusToken tokens.CreateUserLifeStatusFunc, verificationURL string,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p recurringEmailTaskPayload
		if err := unmarshal(t.Payload(), &p); err != nil {
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
