package workers

import (
	"context"
	"log"

	"github.com/bytedance/sonic"
	"github.com/gragorther/epigo/asynq/queues"
	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/database/db"
	"github.com/gragorther/epigo/email"
	"github.com/gragorther/epigo/tokens"
	"github.com/hibiken/asynq"
)

// starts the workers
func Run(ctx context.Context, redisClientOpt asynq.RedisClientOpt, db interface {
	CreateGroup(ctx context.Context, group db.CreateGroup) error
	CreateLastMessage(ctx context.Context, message db.CreateLastMessage) error
	DeleteLastMessageByID(ctx context.Context, id uint) error
	CreateUser(ctx context.Context, user db.CreateUserInput) error
	DeleteGroupByID(ctx context.Context, id uint) error
	UpdateLastMessage(ctx context.Context, id uint, group db.UpdateLastMessage) error
	SetUserMaxSentEmails(ctx context.Context, userID uint, maxSentEmails uint) error
	UpdateGroup(ctx context.Context, id uint, group db.UpdateGroup) error
	UpdateUserInterval(ctx context.Context, userID uint, cron string) error
	IncrementUserSentEmailsCount(ctx context.Context, userID uint) error
	GetUserSentEmails(context.Context, uint) (db.UserSentEmails, error)
	LastMessagesAndRecipients(ctx context.Context, userID uint) (lastMessages []db.LastMessageAndRecipients, err error)
}, jwtSecret []byte, emailService interface {
	SendUserLifeStatusEmail(ctx context.Context, user email.LifeStatusUser, verificationURL string) error
	SendUserDeathEmails(ctx context.Context, name string, emails []email.UserDeathEmailAndRecipients) error
	SendVerificationEmail(ctx context.Context, user email.User, registrationLink string) error
}, registrationRoute string, createVerificationEmailToken tokens.CreateEmailVerificationFunc, createUserLifeStatus tokens.CreateUserLifeStatusFunc, lifeVerificationURL string,
) {
	srv := asynq.NewServer(
		redisClientOpt,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 5,
			// Optionally specify multiple queues with different priority.
			Queues: queues.Queues,
			BaseContext: func() context.Context {
				return ctx
			},
			LogLevel: asynq.DebugLevel,

			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()

	unmarshal := sonic.Unmarshal

	handlerTypes := map[string]asynq.HandlerFunc{
		tasks.TypeCreateGroup:          tasks.HandleCreateGroup(db, unmarshal),
		tasks.TypeCreateLastMessage:    tasks.HandleCreateLastMessage(db, unmarshal),
		tasks.TypeUpdateUserInterval:   tasks.HandleUpdateUserInterval(db, unmarshal),
		tasks.TypeUpdateGroup:          tasks.HandleUpdateGroup(db, unmarshal),
		tasks.TypeSetUserMaxSentEmails: tasks.HandleSetUserMaxSentEmails(db, unmarshal),
		tasks.TypeRecurringEmail:       tasks.HandleRecurringEmail(emailService, db, unmarshal, createUserLifeStatus, lifeVerificationURL),
		tasks.TypeVerificationEmail:    tasks.HandleVerificationEmailTask(createVerificationEmailToken, unmarshal, emailService, registrationRoute),
		tasks.TypeUserDeath:            tasks.HandleUserDeath(db, emailService, unmarshal),
		tasks.TypeDeleteLastMessage:    tasks.HandleDeleteLastMessageByID(db, unmarshal),
		tasks.TypeCreateUser:           tasks.HandleCreateUser(db, unmarshal),
		tasks.TypeDeleteGroup:          tasks.HandleDeleteGroupByID(db, unmarshal),
		tasks.TypeUpdateLastMessage:    tasks.HandleUpdateLastMessage(db, unmarshal),
	}

	for typename, handlerFunc := range handlerTypes {
		mux.HandleFunc(typename, handlerFunc)
	}

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run asynq workers: %v", err.Error())
	}
}
