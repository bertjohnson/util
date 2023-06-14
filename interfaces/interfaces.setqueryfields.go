package interfaces

import (
	"bytes"
	"reflect"
	"strings"
	"time"
)

// SetQueryFields maps string versions of each field into queryValues and concatenates all values into allValuesBuffer.
func SetQueryFields(input interface{}, queryValues map[string]string, allValues *bytes.Buffer) error {
	stringsMap := make(map[string]struct{})
	err := setQueryFieldsHelper(input, queryValues, stringsMap)
	if err != nil {
		return err
	}

	// Concatenate space-separated string values.
	hasValue := false
	for val := range stringsMap {
		if hasValue {
			_, err = allValues.WriteRune(' ')
			if err != nil {
				return err
			}
		} else {
			hasValue = true
		}

		_, err = allValues.WriteString(val)
		if err != nil {
			return err
		}
	}

	return nil
}

// setQueryFieldsHelper maps string versions of each field into queryValues.
func setQueryFieldsHelper(input interface{}, queryValues map[string]string, stringsMap map[string]struct{}) error {
	// Reflect value and iterate through properties.
	reflectType := reflect.TypeOf(input)
	reflectValue := reflect.ValueOf(input)
	if reflectValue.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflectValue.Elem()
	}
	numFields := reflectValue.NumField()
	var err error
	var queryValue string
	if numFields > 0 {
		for i := 0; i < numFields; i++ {
			// Get field type.
			field := reflectValue.Field(i)
			if !field.CanInterface() {
				continue
			}
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}

			// Only process API fields.
			fieldID, ok := reflectType.Field(i).Tag.Lookup("api")
			if !ok {
				continue
			}
			commaPos := strings.Index(fieldID, ",")
			if commaPos > -1 {
				fieldID = fieldID[0:commaPos]
			}

			// Handle shadow fields based on field type.
			switch field.Kind() {
			case reflect.Bool, reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64, reflect.Int, reflect.Interface, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.String, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				queryValue = Sprint(field.Interface())
				if queryValue != "" {
					queryValues[fieldID] = queryValue

					stringsMap[queryValue] = struct{}{}
				}
			case reflect.Map:
				// Only process maps of type map[string]*.
				if field.Type().Key().Kind() == reflect.String {
					for _, key := range field.MapKeys() {
						keyValue := field.MapIndex(key).Interface()
						switch keyValue.(type) {
						case time.Time:
							timeVal := keyValue.(time.Time)
							queryValue = timeVal.UTC().Format("Monday January 02 15:04:05 2006 -0700")
						default:
							queryValue = Sprint(keyValue)
						}

						if queryValue != "" {
							queryValues[key.String()] = queryValue

							stringsMap[queryValue] = struct{}{}
						}
					}
				}
			case reflect.Struct:
				if field.Type().Name() == timeString {
					timeVal := field.Interface().(time.Time)
					queryValue = timeVal.UTC().Format("Monday January 02 15:04:05 2006 -0700")

					if queryValue != "" {
						queryValues[fieldID] = queryValue

						stringsMap[queryValue] = struct{}{}
					}
				} else {
					err = setQueryFieldsHelper(field.Interface(), queryValues, stringsMap)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
