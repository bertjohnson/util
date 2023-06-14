package interfaces

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGetFieldValue tests GetFieldValue().
func TestGetFieldValue(t *testing.T) {
	trueVal := true
	now := time.Now()
	input := SampleStruct{
		BoolVal:        true,
		BoolPointerVal: &trueVal,
		FloatVal:       82471.1419,
		IntVal:         12,
		StringVal:      "parent",
		SubstructVal: SampleSubstruct{
			MapVal: map[string]string{
				"key1": "val1",
			},
		},
		TimeVal:        time.Now(),
		TimePointerVal: &now,
	}

	stringVal := GetFieldValue(input, "stringVal")
	expectedStringVal := "parent"
	assert.Equal(t, expectedStringVal, stringVal)

	key1Value := GetFieldValue(input.SubstructVal, "mapVal.key1")
	expectedKey1Value := "val1"
	assert.Equal(t, expectedKey1Value, key1Value)
}
