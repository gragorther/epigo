package scheduler

import (
	"log"
	"time"

	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/types"
	"github.com/hibiken/asynq"
)

type intervalGetter interface {
	GetUserIntervals() ([]types.UserIntervalsOutput, error)
}

type ConfigProvider struct {
	DB intervalGetter
}

type PeriodicTaskConfigContainer struct {
	Configs []*Config `json:"configs"`
}
type Config struct {
	Cronspec string `json:"cronspec"`
	TaskType string `json:"task_type"`
}

func Run(db intervalGetter, redisAddress string) {
	provider := &ConfigProvider{DB: db}

	mgr, err := asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			RedisConnOpt:               asynq.RedisClientOpt{Addr: redisAddress},
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
	users, err := p.DB.GetUserIntervals()
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
		if user.EmailCron == "" {
			continue
		}
		task, err := tasks.NewRecurringEmailTask(user.ID)
		if err != nil {
			log.Print(err)

			return nil, err
		}

		output = append(output, &asynq.PeriodicTaskConfig{
			Cronspec: user.EmailCron,
			Task:     task,
		})
	}

	return output, nil
}
