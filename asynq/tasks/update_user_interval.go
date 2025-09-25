package tasks

import (
	"context"

	"github.com/hibiken/asynq"
)

type updateUserIntervalPayload struct {
	UserID uint
	Cron   string
}

const TypeUpdateUserInterval = "updateUserInterval"

func (q *queue) UpdateUserInterval(id uint, cron string) error {
	return q.createAndEnqueueTask(updateUserIntervalPayload{
		UserID: id,
		Cron:   cron,
	}, TypeUpdateUserInterval)
}

func HandleUpdateUserInterval(db interface {
	UpdateUserInterval(ctx context.Context, userID uint, cron string) error
}, unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload updateUserIntervalPayload
		if err := unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		return db.UpdateUserInterval(ctx, payload.UserID, payload.Cron)
	}
}
