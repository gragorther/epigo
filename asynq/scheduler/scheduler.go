package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/database/gormdb"
	"github.com/hibiken/asynq"
)

type intervalGetter interface {
	GetUserIntervals(ctx context.Context) ([]gormdb.UserInterval, error)
}

type ConfigProvider struct {
	DB intervalGetter
}

type PeriodicTaskConfigContainer struct {
	Configs []*Config `json:"configs"`
}
type Config struct {
	Cronspec string `json:"cronspec"`
	TaskType string `json:"taskType"`
}

func Run(db intervalGetter, redisClientOpt asynq.RedisClientOpt) {
	provider := &ConfigProvider{DB: db}

	mgr, err := asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			RedisConnOpt:               redisClientOpt,
			PeriodicTaskConfigProvider: provider,    // this provider object is the interface to your config source
			SyncInterval:               time.Second, // this field specifies how often sync should happen
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
	users, err := p.DB.GetUserIntervals(context.TODO())
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		// No users - no tasks to schedule
		return []*asynq.PeriodicTaskConfig{}, nil
	}
	var output []*asynq.PeriodicTaskConfig
	for _, user := range users {

		// if the user has no cron, skip them so asynq doesn't get mad at me
		if user.Cron == "" {
			continue
		}
		task, err := tasks.NewRecurringEmailTask(user.ID)
		if err != nil {
			return nil, err
		}

		output = append(output, &asynq.PeriodicTaskConfig{
			Cronspec: user.Cron,
			Task:     task,
		})
	}

	return output, nil
}
