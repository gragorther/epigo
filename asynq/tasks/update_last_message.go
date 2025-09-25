package tasks

import (
	"context"

	dbHandler "github.com/gragorther/epigo/database/db"
	"github.com/hibiken/asynq"
)

const TypeUpdateLastMessage = "updateLastMessage"

type updateLastMessagePayload struct {
	LastMessage dbHandler.UpdateLastMessage
	ID          uint
}

func (q *queue) UpdateLastMessage(id uint, m dbHandler.UpdateLastMessage) error {
	return q.createAndEnqueueTask(m, TypeUpdateLastMessage)
}

func HandleUpdateLastMessage(
	db interface {
		UpdateLastMessage(ctx context.Context, id uint, group dbHandler.UpdateLastMessage) error
	},
	unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p updateLastMessagePayload
		if err := unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		return db.UpdateLastMessage(ctx, p.ID, p.LastMessage)
	}
}
