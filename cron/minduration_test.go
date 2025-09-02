package cron_test

import (
	"testing"
	"time"

	"github.com/gragorther/epigo/cron"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinDurationBetweenCronTicks(t *testing.T) {
	table := map[string]struct {
		Expr         string
		Iterations   uint
		WantDuration time.Duration
	}{
		"normal time": {
			Expr:         "* * * * *",
			WantDuration: time.Minute,
		},
		"every 2nd minute": {
			Expr:         "*/2 * * * *",
			WantDuration: time.Minute * 2,
		},
		"every month": {
			Expr:         "0 0 1 * *",
			WantDuration: time.Hour * 24 * 28,
		},
		"every minute on monday": {
			Expr:         "* * * * 1",
			WantDuration: time.Minute,
		},
		"every day at 2:02": {
			Expr:         "2 2 * * *",
			WantDuration: time.Hour * 24,
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			gotDuration, err := cron.MinDurationBetweenCronTicks(test.Expr, test.Iterations)
			require.NoError(err)
			assert.Equal(test.WantDuration, gotDuration, "durations should match")
		})
	}
}

func BenchmarkMinDurationBetweenCronTicks(b *testing.B) {
	for b.Loop() {
		_, err := cron.MinDurationBetweenCronTicks("*/2 * * * *", 0)
		if err != nil {
			b.Fatalf("got unexpected error: %v", err)
		}
	}
}
