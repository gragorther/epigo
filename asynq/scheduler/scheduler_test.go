package scheduler_test

import (
	"testing"

	"github.com/gragorther/epigo/asynq/scheduler"
	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/database/gormdb"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

type getIntervalsStub struct {
	err           error
	userIntervals []gormdb.UserInterval
}

func (g getIntervalsStub) GetUserIntervals() ([]gormdb.UserInterval, error) {
	return g.userIntervals, g.err
}

func TestGetConfigs(t *testing.T) {

	assert := assert.New(t)
	t.Run("with users", func(t *testing.T) {
		getIntervals := getIntervalsStub{err: nil, userIntervals: []gormdb.UserInterval{
			{ID: 1, Email: "gregor@gregtech.eu", EmailCron: "5 4 * * *"},
			{ID: 2, Email: "test@gregtech.eu", EmailCron: "5 4 * * *"},
		}}
		task1, err := tasks.NewRecurringEmailTask(1)
		assert.Nil(err, "expected no error when creating new recurring email task")
		task2, err := tasks.NewRecurringEmailTask(2)

		configprovider := scheduler.ConfigProvider{DB: getIntervals}
		got, err := configprovider.GetConfigs()

		assert.Nil(err, "expected no error")
		want := []*asynq.PeriodicTaskConfig{
			{Cronspec: "5 4 * * *", Task: task1},
			{Cronspec: "5 4 * * *", Task: task2},
		}

		assert.Equal(want, got)
	})

	// this tests if the function handles zero configs gracefully
	t.Run("no users", func(t *testing.T) {
		getIntervals := getIntervalsStub{err: nil, userIntervals: []gormdb.UserInterval{}}
		want := []*asynq.PeriodicTaskConfig{}

		configprovider := scheduler.ConfigProvider{DB: getIntervals}

		got, err := configprovider.GetConfigs()
		assert.Nil(err, "expected no error")

		assert.Equal(want, got, "expected zero configs")
	})
}
