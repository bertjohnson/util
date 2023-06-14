package interfaces

import (
	"errors"
	"os"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// SetEnvFieldValues sets field values based on environment variables.
func SetEnvFieldValues(output interface{}) error {
	// If the passed in output isn't a reflection value, unpack it.
	var outputValue reflect.Value
	var ok bool
	if outputValue, ok = output.(reflect.Value); !ok {
		outputValue = reflect.ValueOf(output)
		if outputValue.Kind() != reflect.Ptr {
			return errors.New("pointer required")
		}
		outputValue = outputValue.Elem()
	}
	outputType := outputValue.Type()

	// Loop through each field.
	fieldsLength := outputType.NumField()
fieldLoop:
	for i := 0; i < fieldsLength; i++ {
		field := outputType.Field(i)
		outputField := outputValue.Field(i)

		// Check for env tag.
		envVarName, ok := field.Tag.Lookup("env")
		if !ok {
			continue fieldLoop
		}

		// Check for env variable.
		envVar := os.Getenv(envVarName)
		if envVar == "" {
			continue fieldLoop
		}
		switch field.Type.Kind() {
		case reflect.String:
			if strings.Index(envVar, `"`) < 0 {
				envVar = `"` + envVar + `"`
			}
		}

		// Unmarshal value.
		outputFieldValue := reflect.New(outputField.Type())
		err := jsoniter.Unmarshal([]byte(envVar), outputFieldValue.Interface())
		if err != nil {
			return err
		}
		outputField.Set(outputFieldValue.Elem())
	}

	return nil
}
