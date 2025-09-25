package tasks

import (
	"context"

	"github.com/gragorther/epigo/asynq/queues"
	dbHandler "github.com/gragorther/epigo/database/db"
	"github.com/hibiken/asynq"
)

const TypeCreateUser = "createUser"

func (q *queue) CreateUser(user dbHandler.CreateUserInput) error {
	return q.createAndEnqueueTask(user, TypeCreateUser, asynq.Queue(queues.QueueCritical))
}

func HandleCreateUser(db interface {
	CreateUser(ctx context.Context, user dbHandler.CreateUserInput) error
}, unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p dbHandler.CreateUserInput
		if err := unmarshal(t.Payload(), &p); err != nil {
			return err
		}

		return db.CreateUser(ctx, p)
	}
}
