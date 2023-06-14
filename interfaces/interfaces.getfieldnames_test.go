package interfaces

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestAPIGetFieldNames tests GetAPIFieldNames().
func TestAPIGetFieldNames(t *testing.T) {
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

	fieldNames := GetAPIFieldNames(input)
	expectedFieldNames := []string{"boolVal", "boolPointerVal", "floatVal", "intVal", "interfaceVal", "stringVal", "mapVal", "sliceVal", "timeVal", "timePointerVal", "uintVal"}
	assert.Equal(t, expectedFieldNames, fieldNames)
}

// TestGetFieldNames tests GetFieldNames().
func TestGetFieldNames(t *testing.T) {
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

	fieldNames := GetFieldNames(input)
	expectedFieldNames := []string{"BoolVal", "BoolPointerVal", "FloatVal", "IntVal", "InterfaceVal", "StringVal", "SubstructVal", "SubstructVal.FuncCal", "SubstructVal.MapVal", "SubstructVal.SliceVal", "TimeVal", "TimeVal.wall", "TimeVal.ext", "TimeVal.loc", "TimePointerVal", "TimePointerVal.wall", "TimePointerVal.ext", "TimePointerVal.loc", "UintVal"}
	assert.Equal(t, expectedFieldNames, fieldNames)
}

// TestGetTagFieldNames tests GetTagFieldNames().
func TestGetTagFieldNames(t *testing.T) {
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

	fieldNames := GetTagFieldNames(input, "testtag")
	expectedFieldNames := []string{"boolVal", "stringVal", "uintVal"}
	assert.Equal(t, expectedFieldNames, fieldNames)
}
