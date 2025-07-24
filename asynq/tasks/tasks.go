package tasks

import (
	"context"
	"log"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeRecurringEmail = "email:recurring"
)

// Task payload for any email related tasks.
type emailTaskPayload struct {
	// ID for the email recipient.
	UserID uint
}

func NewRecurringEmailTask(id uint) (*asynq.Task, error) {
	payload, err := sonic.Marshal(emailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRecurringEmail, payload), nil
}

func HandleRecurringEmailTask(ctx context.Context, t *asynq.Task) error {
	var p emailTaskPayload
	if err := sonic.UnmarshalString(string(t.Payload()), &p); err != nil {
		return err
	}
	log.Printf("Sent Email to User %d", p.UserID)
	return nil
}
