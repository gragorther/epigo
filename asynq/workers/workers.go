package workers

import (
	"log"

	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/email"
	"github.com/hibiken/asynq"
)

// starts the workers
func Run(redisClientOpt asynq.RedisClientOpt, jwtSecret []byte, emailService *email.EmailService, registrationRoute string) {
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
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeRecurringEmail, tasks.HandleRecurringEmailTask)
	mux.HandleFunc(tasks.TypeVerificationEmail, tasks.HandleVerificationEmailTask(jwtSecret, emailService, registrationRoute))
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run asynq workers: %v", err.Error())
	}
}
