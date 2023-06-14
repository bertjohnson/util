package interfaces

import (
	"bytes"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSetQueryFields tests SetQueryFields.
func TestSetQueryFields(t *testing.T) {
	sampleTime, err := time.Parse(time.RFC3339, "2018-08-27T11:19:44-05:00")
	assert.NoError(t, err)
	sample := SampleStruct{
		BoolVal:   true,
		FloatVal:  45.67,
		IntVal:    89,
		StringVal: "aaabbb",
		SubstructVal: SampleSubstruct{
			MapVal: map[string]string{
				"aa": "123",
				"bb": "456",
				"cc": "789",
			},
		},
		TimeVal:        sampleTime,
		TimePointerVal: &sampleTime,
	}
	queryFields := make(map[string]string)
	allFields := bytes.Buffer{}

	err = SetQueryFields(&sample, queryFields, &allFields)
	assert.NoError(t, err)
	allFieldParts := strings.Split(allFields.String(), " ")
	sort.Strings(allFieldParts)
	assert.Equal(t, "+0000 0 123 16:19:44 2018 27 45.67 456 789 89 August Monday aaabbb true", strings.Join(allFieldParts, " "))
}
