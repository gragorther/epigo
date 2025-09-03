package cron

import (
	"errors"
	"time"

	"github.com/aptible/supercronic/cronexpr"
)

var ErrInvalidCronExpr = errors.New("invalid cron")

// takes a cron, checks the duration between its ticks and returns the smallest one
//
// defaults to 1000 iterations if the parameter you input is 0
func MinDurationBetweenCronTicks(expr string, iterations uint) (time.Duration, error) {
	expression, err := cronexpr.Parse(expr)
	if err != nil {
		return 0, err
	}

	if iterations == 0 {
		iterations = 1000
	}
	var minDuration time.Duration
	startTick := expression.Next(time.Date(2025, time.January, 0, 0, 0, 0, 0, time.UTC))
	ticks := expression.NextN(startTick, iterations)
	for i, tick := range ticks {
		var durationBetweenCurrentAndNextTick time.Duration

		// if we're at the end of the loop (so it doesn't panic because the array is out of bounds)
		if i == len(ticks)-1 {
			return minDuration, nil
		} else {
			durationBetweenCurrentAndNextTick = ticks[i+1].Sub(tick)
		}

		if durationBetweenCurrentAndNextTick < minDuration || i == 0 {
			minDuration = durationBetweenCurrentAndNextTick
		}
	}
	return minDuration, nil

}
