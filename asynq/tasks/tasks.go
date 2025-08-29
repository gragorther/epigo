package tasks

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeRecurringEmail = "email:recurring"
)

type TaskEnqueuer func(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)

func EnqueueTask(client *asynq.Client) TaskEnqueuer {
	return func(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
		return client.Enqueue(task, opts...)
	}
}

// Task payload for any email related tasks.
type RecurringEmailTaskPayload struct {
	// ID for the email recipient.
	UserID uint
}

func NewRecurringEmailTask(id uint) (*asynq.Task, error) {
	payload, err := sonic.Marshal(RecurringEmailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRecurringEmail, payload), nil
}

func HandleRecurringEmailTask(ctx context.Context, t *asynq.Task) error {
	var p RecurringEmailTaskPayload
	if err := sonic.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	return nil
}
