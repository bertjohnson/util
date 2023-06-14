package funcs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bertjohnson/util/time"
	"github.com/stretchr/testify/assert"
)

// TestRetryN tests RetryN().
func TestRetryN(t *testing.T) {
	var attemptsDone int
	var attemptsPlanned int
	var requiredAttemptsForSuccess int
	var err error

	// Test with an invalid number of retries.
	attemptsDone = 0
	attemptsPlanned = 0
	requiredAttemptsForSuccess = 1
	err = RetryN(func() (bool, error) {
		attemptsDone++
		if attemptsDone < requiredAttemptsForSuccess {
			return false, errors.New(fmt.Sprint("insufficient number of iterations: ", attemptsDone))
		}

		return true, nil
	}, attemptsPlanned)
	assert.NoError(t, err, "Invalid number of retries.")

	// Test with an error.
	attemptsDone = 0
	attemptsPlanned = 0
	requiredAttemptsForSuccess = 1
	err = RetryN(func() (bool, error) {
		return true, errors.New("sample error")
	}, attemptsPlanned)
	assert.Error(t, err, "Expected error.")

	// Test with immediate retries.
	attemptsDone = 0
	attemptsPlanned = 5
	requiredAttemptsForSuccess = 3
	err = RetryN(func() (bool, error) {
		attemptsDone++
		if attemptsDone < requiredAttemptsForSuccess {
			return false, errors.New(fmt.Sprint("insufficient number of iterations: ", attemptsDone))
		}

		return true, nil
	}, attemptsPlanned)
	assert.NoError(t, err, "Immediate retries.")

	// Test with jittered retries.
	attemptsDone = 0
	attemptsPlanned = 5
	requiredAttemptsForSuccess = 3
	err = RetryN(func() (bool, error) {
		attemptsDone++
		if attemptsDone < requiredAttemptsForSuccess {
			// Sleep with jitter.
			time.SleepIncremental(attemptsDone)

			return false, errors.New(fmt.Sprint("insufficient number of iterations: ", attemptsDone))
		}

		return true, nil
	}, attemptsPlanned)
	assert.NoError(t, err, "Jittered retries.")
}

// TestRetryUnlimited tests RetryUnlimited().
func TestRetryUnlimited(t *testing.T) {
	var attemptsDone int
	var requiredAttemptsForSuccess int
	var maximumAttempts int
	var err error

	// Test with an error.
	attemptsDone = 0
	requiredAttemptsForSuccess = 5
	maximumAttempts = 10
	err = RetryUnlimited(func() (bool, error) {
		return true, errors.New("sample error")
	})
	assert.Error(t, err, "Expected error.")

	// Test with immediate retries.
	attemptsDone = 0
	requiredAttemptsForSuccess = 5
	maximumAttempts = 10
	err = RetryUnlimited(func() (bool, error) {
		attemptsDone++
		if attemptsDone < requiredAttemptsForSuccess {
			return false, errors.New(fmt.Sprint("insufficient number of iterations: ", attemptsDone))
		}
		if attemptsDone >= maximumAttempts {
			return false, errors.New(fmt.Sprint("maximum number of iterations reached (safeguard): ", attemptsDone))
		}

		return true, nil
	})
	assert.NoError(t, err, "Immediate retries.")

	// Test with jittered retries.
	attemptsDone = 0
	requiredAttemptsForSuccess = 5
	maximumAttempts = 10
	err = RetryUnlimited(func() (bool, error) {
		attemptsDone++
		if attemptsDone < requiredAttemptsForSuccess {
			// Sleep with jitter.
			time.SleepIncremental(attemptsDone)

			return false, errors.New(fmt.Sprint("insufficient number of iterations: ", attemptsDone))
		}
		if attemptsDone >= maximumAttempts {
			return true, errors.New(fmt.Sprint("maximum number of iterations reached (safeguard): ", attemptsDone))
		}

		return true, nil
	})
	assert.NoError(t, err, "Jittered retries.")
}
