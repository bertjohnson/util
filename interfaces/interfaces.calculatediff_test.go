package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCalculateDiff tests CalculateDiff().
func TestCalculateDiff(t *testing.T) {
	falseVal := false
	trueVal := true
	before := SampleStruct{
		BoolVal:        true,
		BoolPointerVal: &falseVal,
		FloatVal:       9.99,
		IntVal:         10,
		StringVal:      "before value",
		SubstructVal: SampleSubstruct{
			MapVal: map[string]string{
				"A": "1",
				"B": "2",
			},
		},
		UintVal: 87,
	}
	after := SampleStruct{
		BoolVal:        false,
		BoolPointerVal: &trueVal,
		FloatVal:       9.99,
		IntVal:         25,
		StringVal:      "after value",
		SubstructVal: SampleSubstruct{
			MapVal: map[string]string{
				"C": "3",
				"B": "2",
			},
		},
	}

	diffValues := CalculateDiff(before, after)
	expectedDiffValues := map[string]string{
		"boolVal":        "false",
		"boolPointerVal": "true",
		"intVal":         "25",
		"-mapVal_A":      "1",
		"mapVal_C":       "3",
		"stringVal":      "after value",
		"-uintVal":       "87",
	}
	assert.Equal(t, expectedDiffValues, diffValues)
}
