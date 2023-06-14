// Package funcs provides meta-functions for method execution.
package funcs

import (
	"time"

	timeUtils "github.com/bertjohnson/util/time"
)

// RetryN calls method up to retries times with exponential backoff until no error is returned.
func RetryN(method func() (bool, error), retries int) error {
	// Ensure there is a retry.
	if retries < 1 {
		retries = 1
	}

	var isFinal bool
	var err error
	for i := 0; i < retries; i++ {
		isFinal, err = method()
		if isFinal {
			return err
		}
		if err == nil {
			return nil
		}

		timeUtils.SleepIncremental(i)
	}

	return err
}

// RetryUnlimited calls method as many times as required with constant backoff until no error is returned.
func RetryUnlimited(method func() (bool, error)) error {
	var isFinal bool
	var err error
	attempts := 0
	delay := time.Second
	for {
		isFinal, err = method()
		if isFinal {
			return err
		}
		if err == nil {
			return nil
		}

		attempts++
		time.Sleep(delay)
	}
}
