package tasks

import (
	"context"

	"github.com/gragorther/epigo/asynq/queues"
	"github.com/gragorther/epigo/database/db"
	dbHandler "github.com/gragorther/epigo/database/db"
	"github.com/hibiken/asynq"
)

const TypeCreateGroup = "createGroup"

func (q *queue) CreateGroup(group dbHandler.CreateGroup) error {
	return q.createAndEnqueueTask(group, TypeCreateGroup, asynq.Queue(queues.QueueHigh))
}

func HandleCreateGroup(db interface {
	CreateGroup(ctx context.Context, group db.CreateGroup) error
}, unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var group dbHandler.CreateGroup
		if err := unmarshal(t.Payload(), &group); err != nil {
			return err
		}
		return db.CreateGroup(ctx, group)
	}
}
