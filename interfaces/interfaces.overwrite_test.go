package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOverwrite tests Overwrite().
func TestOverwrite(t *testing.T) {
	// Test full inheritance.
	parent := SampleStruct{
		BoolVal:   true,
		FloatVal:  82471.1419,
		IntVal:    12,
		StringVal: "parent",
		SubstructVal: SampleSubstruct{
			MapVal: map[string]string{
				"key1": "val1",
			},
		},
		UintVal: 99,
	}
	child := SampleStruct{}
	Overwrite(&child, []interface{}{&parent})
	assert.NotEqual(t, parent.BoolVal, child.BoolVal)
	assert.Equal(t, parent.FloatVal, child.FloatVal)
	assert.Equal(t, parent.IntVal, child.IntVal)
	assert.Equal(t, parent.StringVal, child.StringVal)
	assert.Equal(t, parent.SubstructVal, child.SubstructVal)
	assert.Equal(t, parent.UintVal, child.UintVal)

	// Test partial inheritance.
	child = SampleStruct{
		StringVal: "child",
		SubstructVal: SampleSubstruct{
			MapVal: map[string]string{
				"key2": "val2",
			},
		},
		UintVal: 250,
	}
	Overwrite(&child, []interface{}{&parent})
	assert.NotEqual(t, parent.BoolVal, child.BoolVal)
	assert.Equal(t, parent.FloatVal, child.FloatVal)
	assert.Equal(t, parent.IntVal, child.IntVal)
	assert.Equal(t, parent.StringVal, child.StringVal)
	assert.NotEqual(t, parent.SubstructVal, child.SubstructVal)
	assert.Equal(t, parent.UintVal, child.UintVal)
}

// TestOverwriteWithTag tests OverwriteWithTag().
func TestOverwriteWithTag(t *testing.T) {
	// Test full inheritance.
	parent := SampleStruct{
		BoolVal:   true,
		FloatVal:  82471.1419,
		IntVal:    12,
		StringVal: "parent",
		SubstructVal: SampleSubstruct{
			MapVal: map[string]string{
				"key1": "val1",
			},
		},
	}
	child := SampleStruct{}
	OverwriteWithTag(&child, []interface{}{&parent}, "inherit")
	assert.NotEqual(t, parent.BoolVal, child.BoolVal)
	assert.Equal(t, parent.FloatVal, child.FloatVal)
	assert.NotEqual(t, parent.IntVal, child.IntVal)
	assert.Equal(t, parent.StringVal, child.StringVal)
	assert.Equal(t, parent.SubstructVal, child.SubstructVal)
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
	OverwriteWithTag(&child, []interface{}{&parent}, "inherit")
	assert.NotEqual(t, parent.BoolVal, child.BoolVal)
	assert.Equal(t, parent.FloatVal, child.FloatVal)
	assert.NotEqual(t, parent.IntVal, child.IntVal)
	assert.Equal(t, parent.StringVal, child.StringVal)
	assert.NotEqual(t, parent.SubstructVal, child.SubstructVal)
	assert.NotEqual(t, parent.UintVal, child.UintVal)
}
