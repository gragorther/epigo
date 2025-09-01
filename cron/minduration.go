package cron

import (
	"errors"
	"time"

	"github.com/adhocore/gronx"
)

var ErrInvalidCronExpr = errors.New("invalid cron")

// takes a cron, checks the duration between its ticks and returns the smallest one
//
// defaults to 100 iterations if the parameter you input is 0
func MinDurationBetweenCronTicks(expr string, start time.Time, iterations uint) (time.Duration, error) {

	// validate
	if !gronx.IsValid(expr) {
		return 0, ErrInvalidCronExpr
	}
	if iterations == 0 {
		iterations = 1000
	}
	var minDuration time.Duration
	var previousTick time.Time
	previousTick = start
	for i := range iterations {
		nextTick, err := gronx.NextTickAfter(expr, previousTick, false)
		if err != nil {
			return 0, err
		}

		durationBetweenPreviousTickAndCurrentTick := nextTick.Sub(previousTick)
		if durationBetweenPreviousTickAndCurrentTick < minDuration || i == 0 {
			minDuration = durationBetweenPreviousTickAndCurrentTick
		}

		previousTick = nextTick

	}
	return minDuration, nil

}
