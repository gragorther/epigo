package tasks

import (
	"context"

	"github.com/gragorther/epigo/asynq/queues"
	dbHandler "github.com/gragorther/epigo/database/db"
	"github.com/hibiken/asynq"
)

const TypeCreateLastMessage = "createLastMessage"

func (q *queue) CreateLastMessage(message dbHandler.CreateLastMessage) error {
	return q.createAndEnqueueTask(message, TypeCreateLastMessage, asynq.Queue(queues.QueueHigh))
}

func HandleCreateLastMessage(db interface {
	CreateLastMessage(ctx context.Context, message dbHandler.CreateLastMessage) error
}, unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p dbHandler.CreateLastMessage
		if err := unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		return db.CreateLastMessage(ctx, p)
	}
}
