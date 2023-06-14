// Package time provides functions for parsing time and sleeping.
package time

import (
	"math/rand"
	"strings"
	"time"
)

// increment defines sleep increments.
type increment struct {
	FixedAmount    time.Duration
	VariableAmount time.Duration
}

// Standard jitter increments.
var (
	increments = []increment{
		{FixedAmount: 250 * time.Millisecond, VariableAmount: 250 * time.Millisecond},
		{FixedAmount: 500 * time.Millisecond, VariableAmount: 500 * time.Millisecond},
		{FixedAmount: 1 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 2 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 4 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 8 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 16 * time.Second, VariableAmount: 1 * time.Second},
		{FixedAmount: 32 * time.Second, VariableAmount: 2 * time.Second},
		{FixedAmount: 64 * time.Second, VariableAmount: 2 * time.Second},
	}
)

// Parse parses an arbitrary time string, attempting to determine its layout.
func Parse(input string) (time.Time, error) {
	var format string
	if strings.HasSuffix(input, "Z") {
		format = "2006-01-02T15:04:05Z"
	} else {
		format = "2006-01-02T15:04:05-07:00"
	}

	return time.Parse(format, input)
}

// SleepIncremental sleeps using incremental backoff.
func SleepIncremental(increment int) {
	if increment < 0 {
		increment = 0
	}
	if increment > 8 {
		increment = 8
	}

	delay := increments[increment].FixedAmount + time.Duration(rand.Int31n(int32(increments[increment].VariableAmount)))
	time.Sleep(delay)
}

// SleepUntil sleeps until a specified time.
func SleepUntil(when time.Time) {
	now := time.Now()
	if now.Before(when) {
		time.Sleep(when.Sub(now))
	}
}
