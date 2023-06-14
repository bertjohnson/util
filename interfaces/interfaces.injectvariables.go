package interfaces

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// InjectVariables injects input variables into a formatted string.
func InjectVariables(input string, variables map[string]interface{}) (interface{}, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var c, lastC rune
	var err error
	var hasVariable, inVariable bool
	var outputBuffer, variableBuffer bytes.Buffer
	for {
		if c, _, err = reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		} else {
			if inVariable {
				if c == '}' {
					inVariable = false

					inputVariableName := variableBuffer.String()
					variableValue, ok := variables[inputVariableName]
					if !ok {
						return "", errors.New("input variable (" + inputVariableName + ") not found")
					}

					_, err = outputBuffer.WriteString(fmt.Sprint(variableValue))
					if err != nil {
						return "", err
					}

					variableBuffer.Reset()
				} else {
					_, err = variableBuffer.WriteRune(c)
					if err != nil {
						return "", err
					}
				}
			} else {
				if c == '{' && lastC == '$' {
					inVariable = true
					hasVariable = true
					outputBuffer.Truncate(outputBuffer.Len() - 1)
				} else {
					_, err = outputBuffer.WriteRune(c)
					if err != nil {
						return "", err
					}
				}
			}
			lastC = c
		}
	}
	if hasVariable {
		return outputBuffer.String(), nil
	}

	// Individual variables not found. Treat entire input as a variable.
	output, ok := variables[input]
	if !ok {
		return "", errors.New("input variable (" + input + ") not found")
	}
	return output, nil
}
