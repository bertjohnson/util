package interfaces

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SampleEnvStruct contains fields to test SetEnvFieldValues.
type SampleEnvStruct struct {
	BoolVal1     *bool  `default:"false" env:"SAMPLE_BOOLVAL1"`
	BoolVal2     *bool  `default:"true" env:"SAMPLE_BOOLVAL2"`
	IntVal       int    `default:"56" env:"SAMPLE_INTVAL"`
	SetStringVal string `default:"test" env:"SAMPLE_SETSTRINGVAL"`
	StringVal    string `default:"test" env:"SAMPLE_STRINGVAL"`
}

// TestSetEnvFieldValues tests SetEnvFieldValues().
func TestSetEnvFieldValues(t *testing.T) {
	sample := &SampleEnvStruct{
		SetStringVal: "already set",
	}
	trueVal := true
	falseVal := false
	expectedSample := &SampleEnvStruct{
		BoolVal1:     &trueVal,
		BoolVal2:     &falseVal,
		IntVal:       60606,
		SetStringVal: "replaced",
		StringVal:    "lalala",
	}
	os.Setenv("SAMPLE_BOOLVAL1", "true")
	os.Setenv("SAMPLE_BOOLVAL2", "false")
	os.Setenv("SAMPLE_INTVAL", "60606")
	os.Setenv("SAMPLE_SETSTRINGVAL", "replaced")
	os.Setenv("SAMPLE_STRINGVAL", "lalala")
	err := SetEnvFieldValues(sample)
	assert.NoError(t, err)
	assert.Equal(t, expectedSample, sample)
}
