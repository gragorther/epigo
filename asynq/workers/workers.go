package workers

import (
	"context"
	"log"

	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/email"
	"github.com/gragorther/epigo/tokens"
	"github.com/hibiken/asynq"
)

// starts the workers
func Run(ctx context.Context, redisClientOpt asynq.RedisClientOpt, jwtSecret []byte, emailService *email.EmailService, registrationRoute string, createVerificationEmailToken tokens.CreateEmailVerificationFunc, createUserLifeStatus tokens.CreateUserLifeStatusFunc, lifeVerificationURL string) {
	srv := asynq.NewServer(
		redisClientOpt,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 5,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			BaseContext: func() context.Context {
				return ctx
			},

			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeRecurringEmail, tasks.HandleRecurringEmailTask(emailService, createUserLifeStatus, lifeVerificationURL))
	mux.HandleFunc(tasks.TypeVerificationEmail, tasks.HandleVerificationEmailTask(createVerificationEmailToken, emailService, registrationRoute))
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run asynq workers: %v", err.Error())
	}
}
