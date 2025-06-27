package scheduler

import (
	"log"
	"time"

	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/db"
	"github.com/gragorther/epigo/tasks"
	"github.com/hibiken/asynq"
)

type ConfigProvider struct {
	db *db.DBHandler
}

type PeriodicTaskConfigContainer struct {
	Configs []*Config `json:"configs"`
}
type Config struct {
	Cronspec string `json:"cronspec"`
	TaskType string `json:"task_type"`
}

func Run(db *db.DBHandler, redisAddress string) {
	provider := &ConfigProvider{db: db}

	mgr, err := asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			RedisConnOpt:               asynq.RedisClientOpt{Addr: redisAddress},
			PeriodicTaskConfigProvider: provider,        // this provider object is the interface to your config source
			SyncInterval:               2 * time.Minute, // this field specifies how often sync should happen
		})
	if err != nil {
		log.Fatal(err)
	}

	if err := mgr.Run(); err != nil {
		log.Fatal(err)
	}
	defer mgr.Shutdown()
}

func (p *ConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	users, err := p.db.GetUserIntervals()
	if err != nil {
		return nil, err
	}
	if len(users) > 1 {
		return nil, apperrors.ErrNoUsers
	}
	var output []*asynq.PeriodicTaskConfig
	for _, user := range users {
		task, err := tasks.NewRecurringEmailTask(user.ID)
		if err != nil {
			return nil, err
		}
		output = append(output, &asynq.PeriodicTaskConfig{
			Cronspec: user.EmailCron,
			Task:     task,
		})
	}
	return output, nil
}
