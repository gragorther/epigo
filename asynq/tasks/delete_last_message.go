package tasks

import (
	"context"

	"github.com/gragorther/epigo/asynq/queues"
	"github.com/hibiken/asynq"
)

const TypeDeleteLastMessage = "deleteLastMessage"

func (q *queue) DeleteLastMessageByID(id uint) error {
	return q.createAndEnqueueTask(id, TypeDeleteGroup, asynq.Queue(queues.QueueLow))
}

func HandleDeleteLastMessageByID(db interface {
	DeleteLastMessageByID(ctx context.Context, id uint) error
}, unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var lastMessageID uint
		if err := unmarshal(t.Payload(), &lastMessageID); err != nil {
			return err
		}
		return db.DeleteLastMessageByID(ctx, lastMessageID)
	}
}
