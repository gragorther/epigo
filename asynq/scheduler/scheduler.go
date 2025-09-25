package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/database/db"
	"github.com/hibiken/asynq"
)

type configProviderDB interface {
	AllUserIntervalsAndSentEmails(ctx context.Context) (intervals []db.IntervalAndSentEmails, err error)
}

// takes a task Enqueuer in order to be able to schedule one-off tasks like
type configProvider struct {
	DB          configProviderDB
	EnqueueTask tasks.TaskEnqueueFunc
}

func Run(db configProviderDB, redisClientOpt asynq.RedisClientOpt) {
	client := asynq.NewClient(redisClientOpt)
	enqueuer := tasks.EnqueueTask(client)
	provider := &configProvider{DB: db, EnqueueTask: enqueuer}
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

func (p *configProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	ctx := context.Background()
	users, err := p.DB.AllUserIntervalsAndSentEmails(ctx)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		// No users - no tasks to schedule
		return nil, nil
	}
	var output []*asynq.PeriodicTaskConfig
	for _, user := range users {

		// if the user has no cron, skip them so asynq doesn't get mad at me
		if user.Cron == "" {
			continue
		}
		if user.SentEmails > user.MaxSentEmails+1 {
			continue
		}
		if user.SentEmails == user.MaxSentEmails+1 {
			deathTask, err := tasks.NewUserDeath(user.ID, user.Name, sonic.Marshal)
			if err != nil {
				return nil, err
			}
			p.EnqueueTask(deathTask)
			continue
		}
		task, err := tasks.NewRecurringEmailTask(user.ID, user.Name, user.Email, 24*time.Hour)
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
