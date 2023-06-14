package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestParse tests Parse().
func TestParse(t *testing.T) {
	var timeString string
	var result time.Time
	var expectedResult time.Time
	var err error

	// Test full format.
	timeString = "2017-01-08T00:39:40-05:00"
	expectedResult, err = time.Parse("2006-01-02T15:04:05-07:00", "2017-01-08T00:39:40-05:00")
	assert.NoError(t, err, "Parsed time.")
	result, err = Parse(timeString)
	assert.NoError(t, err, "Parsed time.")
	assert.Equal(t, expectedResult, result, "Parsed time.")

	// Test partial format.
	timeString = "2017-01-08T00:39:40Z"
	expectedResult, err = time.Parse("2006-01-02T15:04:05Z", "2017-01-08T00:39:40Z")
	assert.NoError(t, err, "Parsed time.")
	result, err = Parse(timeString)
	assert.NoError(t, err, "Parsed time.")
	assert.Equal(t, expectedResult, result, "Parsed time.")

	// Test invalid format.
	timeString = "today"
	_, err = Parse(timeString)
	assert.Error(t, err, "Invalid format.")
}

// TestSleepIncremental tests SleepIncremental().
func TestSleepIncremental(t *testing.T) {
	// Run five iterations.
	fiveSeconds := time.Now().Add(5 * time.Second)
	for increments := 0; increments < 5; increments++ {
		SleepIncremental(increments)
	}
	assert.False(t, fiveSeconds.After(time.Now()), "Woke too soon.")

	// Check an invalid iteration.
	SleepIncremental(-999)

	// Check the maximum iteration.
	oneMinute := time.Now().Add(time.Minute)
	SleepIncremental(999)
	assert.False(t, oneMinute.After(time.Now()), "Woke too soon.")
}

// TestSleepUntil tests SleepUntil().
func TestSleepUntil(t *testing.T) {
	fiveSeconds := time.Now().Add(5 * time.Second)
	SleepUntil(fiveSeconds)
	assert.False(t, fiveSeconds.After(time.Now()), "Woke too soon.")
}
