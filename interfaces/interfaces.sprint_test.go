package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSprint tests Sprint().
func TestSprint(t *testing.T) {
	var sprint, expectedSprint string

	// Test with object.
	sampleObject := SampleStruct{
		BoolVal:   true,
		StringVal: "object",
	}
	expectedSprint = "{true <nil> 0 0 <nil> object {<nil> map[] []} 0001-01-01 00:00:00 +0000 UTC <nil> 0}"
	sprint = Sprint(sampleObject)
	assert.Equal(t, expectedSprint, sprint)

	// Test with number.
	sampleNumber := 567
	expectedSprint = "567"
	sprint = Sprint(sampleNumber)
	assert.Equal(t, expectedSprint, sprint)

	// Test with string.
	sampleString := "Hello world"
	expectedSprint = "Hello world"
	sprint = Sprint(sampleString)
	assert.Equal(t, expectedSprint, sprint)

	// Test with nil.
	expectedSprint = ""
	sprint = Sprint(nil)
	assert.Equal(t, expectedSprint, sprint)
}
