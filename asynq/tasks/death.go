package tasks

import (
	"context"

	dbHandler "github.com/gragorther/epigo/database/db"
	"github.com/gragorther/epigo/email"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
)

const TypeUserDeath = "userDeath"

type userDeathPayload struct {
	UserID uint
	Name   string
}

func NewUserDeath(userID uint, name string, marshal MarshalFunc) (task *asynq.Task, err error) {
	payload, err := marshal(userDeathPayload{UserID: userID, Name: name})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeUserDeath, payload), nil
}

func HandleUserDeath(db interface {
	LastMessagesAndRecipients(ctx context.Context, userID uint) (lastMessages []dbHandler.LastMessageAndRecipients, err error)
}, emailService interface {
	SendUserDeathEmails(ctx context.Context, name string, emails []email.UserDeathEmailAndRecipients) error
}, unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload userDeathPayload
		if err := unmarshal(task.Payload(), &payload); err != nil {
			return err
		}
		lastMessages, err := db.LastMessagesAndRecipients(ctx, payload.UserID)
		if err != nil {
			return err
		}

		emailsAndRecipients := lo.Map(lastMessages, func(item dbHandler.LastMessageAndRecipients, _ int) email.UserDeathEmailAndRecipients {
			return email.UserDeathEmailAndRecipients{
				Title:           item.LastMessage.Title,
				Content:         item.LastMessage.Content.String,
				RecipientEmails: dbHandler.RecipientsToStringArray(item.Recipients),
			}
		})
		if err := emailService.SendUserDeathEmails(ctx, payload.Name, emailsAndRecipients); err != nil {
			return err
		}
		return nil
	}
}
