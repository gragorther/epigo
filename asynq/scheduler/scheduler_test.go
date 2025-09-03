package scheduler_test

/*
type getIntervalsStub struct {
	err           error
	userIntervals []gormdb.UserInterval
}

func (g getIntervalsStub) GetUserIntervals(ctx context.Context) ([]gormdb.UserInterval, error) {
	return g.userIntervals, g.err
}
*/

/*
func TestGetConfigs(t *testing.T) {

	assert := assert.New(t)
	t.Run("with users", func(t *testing.T) {
		require := require.New(t)
		getIntervals := getIntervalsStub{err: nil, userIntervals: []gormdb.UserInterval{
			{ID: 1, Email: "gregor@gregtech.eu", Cron: "5 4 * * *"},
			{ID: 2, Email: "test@gregtech.eu", Cron: "5 4 * * *"},
		}}
		task1, err := tasks.NewRecurringEmailTask(1)
		assert.Nil(err, "expected no error when creating new recurring email task")
		task2, err := tasks.NewRecurringEmailTask(2)
		require.NoError(err)

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
*/
