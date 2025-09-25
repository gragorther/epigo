package tasks

import (
	"context"

	"github.com/gragorther/epigo/asynq/queues"
	dbHandler "github.com/gragorther/epigo/database/db"
	"github.com/hibiken/asynq"
)

const TypeUpdateGroup = "updateGroup"

type updateGroupPayload struct {
	ID    uint `json:"groupID"`
	Group dbHandler.UpdateGroup
}

func (q queue) createTask(payload any, typeName string, opts ...asynq.Option) (*asynq.Task, error) {
	marshaledPayload, err := q.marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(typeName, marshaledPayload, opts...), nil
}

func (q *queue) createAndEnqueueTask(payload any, typeName string, opts ...asynq.Option) error {
	task, err := q.createTask(payload, typeName, opts...)
	if err != nil {
		return err
	}
	_, err = q.enqueueTask(task, opts...)
	return err
}

func (q *queue) UpdateGroup(id uint, group dbHandler.UpdateGroup) error {
	return q.createAndEnqueueTask(updateGroupPayload{
		ID:    id,
		Group: group,
	}, TypeUpdateGroup, asynq.Queue(queues.QueueDefault))
}

func HandleUpdateGroup(
	db interface {
		UpdateGroup(ctx context.Context, id uint, group dbHandler.UpdateGroup) error
	},
	unmarshal UnmarshalFunc,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p updateGroupPayload
		if err := unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		return db.UpdateGroup(ctx, p.ID, p.Group)
	}
}
