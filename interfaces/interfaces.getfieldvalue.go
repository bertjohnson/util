package interfaces

import (
	"reflect"
	"strings"

	utilstrings "github.com/bertjohnson/util/strings"
)

// GetFieldValue returns the value of a named field.
func GetFieldValue(value interface{}, fieldName string) interface{} {
	reflectValue := reflect.ValueOf(value)
	for {
		if reflectValue.Kind() == reflect.Ptr {
			reflectValue = reflectValue.Elem()
		} else {
			break
		}
	}

	// Parse dot notation.
	periodPos := strings.Index(fieldName, ".")
	if periodPos > -1 {
		fieldValue := reflectValue.FieldByName(utilstrings.Capitalize(fieldName[0:periodPos]))
		if fieldValue.IsValid() {
			return GetFieldValue(fieldValue.Interface(), fieldName[periodPos+1:])
		}

		return nil
	}

	// Resolve maps.
	if reflectValue.Kind() == reflect.Map {
		fieldValue := reflectValue.MapIndex(reflect.ValueOf(fieldName))
		if fieldValue.IsValid() {
			return fieldValue.Interface()
		}

		return nil
	}

	// Resolve structs.
	fieldValue := reflectValue.FieldByName(utilstrings.Capitalize(fieldName))
	if fieldValue.IsValid() {
		return fieldValue.Interface()
	}

	return nil
}
