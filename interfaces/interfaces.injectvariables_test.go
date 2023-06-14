package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInjectVariables tests InjectVariables().
func TestInjectVariables(t *testing.T) {
	var input, expectedOutput string
	var values map[string]interface{}
	var output interface{}
	var err error

	input = "test${value}"
	values = map[string]interface{}{
		"value": "success",
	}
	expectedOutput = "testsuccess"
	output, err = InjectVariables(input, values)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)

	input = "test${unmatched}"
	values = map[string]interface{}{}
	expectedOutput = ""
	output, err = InjectVariables(input, values)
	assert.Error(t, err)
	assert.Equal(t, expectedOutput, output)
}
