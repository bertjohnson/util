package interfaces

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// SampleStruct contains fields to test inheritance.
type SampleStruct struct {
	BoolVal        bool            `api:"boolVal,omitempty" inherit:"dynamic" testtag:"boolVal"`
	BoolPointerVal *bool           `api:"boolPointerVal" inherit:"dynamic"`
	FloatVal       float64         `api:"floatVal" inherit:"dynamic"`
	IntVal         int             `api:"intVal"`
	InterfaceVal   interface{}     `api:"interfaceVal"`
	StringVal      string          `api:"stringVal" inherit:"dynamic" testtag:"stringVal"`
	SubstructVal   SampleSubstruct `api:",inline" inherit:"dynamic"`
	TimeVal        time.Time       `api:"timeVal,omitempty" inherit:"dynamic"`
	TimePointerVal *time.Time      `api:"timePointerVal" inherit:"dynamic"`
	UintVal        uint            `api:"uintVal" testtag:"uintVal"`
}

// SampleSubstruct contains fields to test inheritance.
type SampleSubstruct struct {
	FuncCal  func()
	MapVal   map[string]string `api:"mapVal" inherit:"dynamic"`
	SliceVal []int             `api:"sliceVal"`
}

// TestInherit tests Inherit().
func TestInherit(t *testing.T) {
	// Test full inheritance.
	trueVal := true
	now := time.Now()
	parent := SampleStruct{
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
	child := SampleStruct{}
	Inherit(&child, []interface{}{&parent})
	assert.NotEqual(t, parent.BoolVal, child.BoolVal)
	assert.Equal(t, parent.BoolPointerVal, child.BoolPointerVal)
	assert.Equal(t, parent.FloatVal, child.FloatVal)
	assert.Equal(t, parent.IntVal, child.IntVal)
	assert.Equal(t, parent.StringVal, child.StringVal)
	assert.Equal(t, parent.SubstructVal, child.SubstructVal)
	assert.Equal(t, parent.TimeVal, child.TimeVal)
	assert.Equal(t, parent.TimePointerVal, child.TimePointerVal)
	assert.Equal(t, parent.UintVal, child.UintVal)

	// Test partial inheritance.
	child = SampleStruct{
		StringVal: "child",
		SubstructVal: SampleSubstruct{
			MapVal: map[string]string{
				"key2": "val2",
			},
		},
		UintVal: 98,
	}
	Inherit(&child, []interface{}{&parent})
	assert.NotEqual(t, parent.BoolVal, child.BoolVal)
	assert.Equal(t, parent.BoolPointerVal, child.BoolPointerVal)
	assert.Equal(t, parent.FloatVal, child.FloatVal)
	assert.Equal(t, parent.IntVal, child.IntVal)
	assert.NotEqual(t, parent.StringVal, child.StringVal)
	assert.NotEqual(t, parent.SubstructVal, child.SubstructVal)
	assert.Equal(t, parent.TimeVal, child.TimeVal)
	assert.Equal(t, parent.TimePointerVal, child.TimePointerVal)
	assert.NotEqual(t, parent.UintVal, child.UintVal)
}

// TestInheritWithTag tests InheritWithTag().
func TestInheritWithTag(t *testing.T) {
	// Test full inheritance.
	trueVal := true
	now := time.Now()
	parent := SampleStruct{
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
	child := SampleStruct{}
	InheritWithTag(&child, []interface{}{&parent}, "inherit")
	assert.NotEqual(t, parent.BoolVal, child.BoolVal)
	assert.Equal(t, parent.BoolPointerVal, child.BoolPointerVal)
	assert.Equal(t, parent.FloatVal, child.FloatVal)
	assert.NotEqual(t, parent.IntVal, child.IntVal)
	assert.Equal(t, parent.StringVal, child.StringVal)
	assert.Equal(t, parent.SubstructVal, child.SubstructVal)
	assert.Equal(t, parent.TimeVal, child.TimeVal)
	assert.Equal(t, parent.TimePointerVal, child.TimePointerVal)
	assert.Equal(t, parent.UintVal, child.UintVal)

	// Test partial inheritance.
	child = SampleStruct{
		StringVal: "child",
		SubstructVal: SampleSubstruct{
			MapVal: map[string]string{
				"key2": "val2",
			},
		},
		UintVal: 98,
	}
	InheritWithTag(&child, []interface{}{&parent}, "inherit")
	assert.NotEqual(t, parent.BoolVal, child.BoolVal)
	assert.Equal(t, parent.BoolPointerVal, child.BoolPointerVal)
	assert.Equal(t, parent.FloatVal, child.FloatVal)
	assert.NotEqual(t, parent.IntVal, child.IntVal)
	assert.NotEqual(t, parent.StringVal, child.StringVal)
	assert.NotEqual(t, parent.SubstructVal, child.SubstructVal)
	assert.Equal(t, parent.TimeVal, child.TimeVal)
	assert.Equal(t, parent.TimePointerVal, child.TimePointerVal)
	assert.NotEqual(t, parent.UintVal, child.UintVal)
}
