package tasks

import (
	"context"

	"github.com/bytedance/sonic"
	dbHandler "github.com/gragorther/epigo/database/db"
	"github.com/gragorther/epigo/email"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
)

const TypeUserDeath = "userDeath"

type UserDeathPayload struct {
	UserID uint
	Name   string
}

func NewUserDeathTask(userID uint, name string) (task *asynq.Task, err error) {
	payload, err := sonic.Marshal(UserDeathPayload{UserID: userID, Name: name})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeUserDeath, payload), nil
}

func HandleUserDeathTask(db interface {
	LastMessagesAndRecipients(ctx context.Context, userID uint) (lastMessages []dbHandler.LastMessageAndRecipients, err error)
}, emailService interface {
	SendUserDeathEmails(ctx context.Context, name string, emails []email.UserDeathEmailAndRecipients) error
},
) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload UserDeathPayload
		if err := sonic.Unmarshal(task.Payload(), &payload); err != nil {
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
