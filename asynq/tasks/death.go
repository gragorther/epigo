package tasks

import (
	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

const TypeUserDeath = "userDeath"

type UserDeathPayload struct {
	UserID uint
}

func NewUserDeathTask(userID uint) (task *asynq.Task, err error) {
	payload, err := sonic.Marshal(UserDeathPayload{UserID: userID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeUserDeath, payload), nil
}
