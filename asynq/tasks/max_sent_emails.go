package tasks

import (
	"context"

	"github.com/hibiken/asynq"
)

const TypeSetUserMaxSentEmails = "setUserMaxSentEmails"

type setUserMaxSentEmailsPayload struct {
	UserID        uint
	MaxSentEmails uint
}

func (q *queue) SetUserMaxSentEmails(userID uint, maxSentEmails uint) error {
	return q.createAndEnqueueTask(setUserMaxSentEmailsPayload{UserID: userID, MaxSentEmails: maxSentEmails}, TypeSetUserMaxSentEmails)
}

func HandleSetUserMaxSentEmails(
	db interface {
		SetUserMaxSentEmails(ctx context.Context, userID uint, maxSentEmails uint) error
	}, unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload setUserMaxSentEmailsPayload
		if err := unmarshal(t.Payload(), &payload); err != nil {
			return err
		}
		return db.SetUserMaxSentEmails(ctx, payload.UserID, payload.MaxSentEmails)
	}
}
