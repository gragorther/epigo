package tasks

import (
	"context"
	"encoding/json"
	"log"

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
	payload, err := json.Marshal(emailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRecurringEmail, payload), nil
}

func HandleRecurringEmailTask(ctx context.Context, t *asynq.Task) error {
	var p emailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Send Welcome Email to User %d", p.UserID)
	return nil
}
