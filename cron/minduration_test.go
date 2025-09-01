package cron_test

import (
	"testing"
	"time"

	"github.com/gragorther/epigo/cron"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinDurationBetweenCronTicks(t *testing.T) {
	startTime := time.Date(2025, time.January, 2, 1, 0, 0, 0, time.UTC)
	table := map[string]struct {
		Expr         string
		Start        time.Time
		Iterations   uint
		WantDuration time.Duration
	}{
		"normal time": {
			Expr:         "* * * * *",
			WantDuration: time.Minute,
			Start:        startTime,
		},
		"every 2nd minute": {
			Expr:         "*/2 * * * *",
			WantDuration: time.Minute * 2,
			Start:        startTime,
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			gotDuration, err := cron.MinDurationBetweenCronTicks(test.Expr, test.Start, test.Iterations)
			require.NoError(err)
			assert.Equal(test.WantDuration, gotDuration)
		})
	}
}
