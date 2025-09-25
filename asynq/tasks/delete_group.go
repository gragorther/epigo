package tasks

import (
	"context"

	"github.com/gragorther/epigo/asynq/queues"
	"github.com/hibiken/asynq"
)

const TypeDeleteGroup = "deleteGroup"

func (q *queue) DeleteGroupByID(id uint) error {
	return q.createAndEnqueueTask(id, TypeDeleteGroup, asynq.Queue(queues.QueueLow))
}

func HandleDeleteGroupByID(db interface {
	DeleteGroupByID(ctx context.Context, id uint) error
}, unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var groupID uint
		if err := unmarshal(t.Payload(), &groupID); err != nil {
			return err
		}
		return db.DeleteGroupByID(ctx, groupID)
	}
}
