package interfaces

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// Accessor character to delineate parent/child relationships. Typically ".", but that is invalid for MongoDB, so we use "_".
	accessorCharacter = "_"
)

// CalculateDiff calculates the difference between two objects.
func CalculateDiff(before interface{}, after interface{}) map[string]string {
	diffFields := make(map[string]string)
	calculateDiffHelper(before, after, "", diffFields)

	return diffFields
}

// calculateDiffHelper calculates the difference between two objects.
func calculateDiffHelper(before interface{}, after interface{}, prefix string, diffFields map[string]string) {
	// Cast to reflected values.
	beforeValue := reflect.ValueOf(before)
	if beforeValue.Kind() == reflect.Ptr {
		beforeValue = beforeValue.Elem()
	}
	beforeType := reflect.TypeOf(before)
	if beforeType.Kind() == reflect.Ptr {
		beforeType = beforeType.Elem()
	}
	afterValue := reflect.ValueOf(after)
	if afterValue.Kind() == reflect.Ptr {
		afterValue = afterValue.Elem()
	}
	afterType := reflect.TypeOf(after)
	if afterType.Kind() == reflect.Ptr {
		afterType = afterType.Elem()
	}

	// Iterate through before fields.
	beforeFieldsCount := beforeValue.NumField()
	processedFields := make(map[string]struct{})
beforeFieldLoop:
	for i := 0; i < beforeFieldsCount; i++ {
		beforeField := beforeValue.Field(i)
	indirectBeforeLoop:
		for {
			if beforeField.Kind() == reflect.Ptr {
				if beforeField.IsNil() {
					continue beforeFieldLoop
				}
				beforeField = beforeField.Elem()
				continue indirectBeforeLoop
			}
			break
		}
		beforeFieldType := beforeType.Field(i)
		fieldName := beforeFieldType.Name
		fieldTag, ok := beforeFieldType.Tag.Lookup("api")
		if !ok {
			continue
		}
		commaPos := strings.Index(fieldTag, ",")
		if commaPos > -1 {
			fieldTag = fieldTag[0:commaPos]
		}
		if len(prefix) > 0 {
			fieldTag = prefix + accessorCharacter + fieldTag
		}
		processedFields[fieldTag] = struct{}{}

		// Handle deleted fields.
		afterField := afterValue.FieldByName(fieldName)
	indirectAfterLoop:
		for {
			if afterField.Kind() == reflect.Ptr {
				if afterField.IsNil() {
					continue beforeFieldLoop
				}
				afterField = afterField.Elem()
				continue indirectAfterLoop
			}
			break
		}
		beforeFieldInterface := beforeField.Interface()
		if fieldTag != "" {
			if !afterField.IsValid() {
				diffFields["-"+fieldTag] = fmt.Sprint(beforeFieldInterface)
				continue
			}
		}

		// Handle updated fields.
		afterFieldInterface := afterField.Interface()
		if !reflect.DeepEqual(beforeFieldInterface, afterFieldInterface) {
			beforeFieldKind := beforeField.Kind()
			if beforeFieldKind != reflect.Bool && reflect.DeepEqual(afterFieldInterface, reflect.Zero(beforeFieldType.Type).Interface()) {
				diffFields["-"+fieldTag] = fmt.Sprint(beforeFieldInterface)
				continue
			}

			// Recurse through structs.
			switch beforeFieldKind {
			case reflect.Map:
				if afterField.Kind() == reflect.Map {
					// Handle updated fields.
					processedMapKeys := make(map[interface{}]bool)
					for _, key := range beforeField.MapKeys() {
						processedMapKeys[key.Interface()] = true
						beforeMapValue := beforeField.MapIndex(key)
						afterMapValue := afterField.MapIndex(key)
						if afterMapValue.IsValid() {
							if !reflect.DeepEqual(beforeMapValue.Interface(), afterMapValue.Interface()) {
								diffFields[fieldTag+accessorCharacter+fmt.Sprint(key)] = fmt.Sprint(afterMapValue.Interface())
							}
						} else {
							diffFields["-"+fieldTag+accessorCharacter+fmt.Sprint(key)] = fmt.Sprint(beforeMapValue.Interface())
						}
					}

					// Handle deleted fields.
					for _, key := range afterField.MapKeys() {
						if _, ok := processedMapKeys[key.Interface()]; !ok {
							afterMapValue := afterField.MapIndex(key)
							diffFields[fieldTag+accessorCharacter+fmt.Sprint(key)] = fmt.Sprint(afterMapValue.Interface())
						}
					}
				}
			case reflect.Struct:
				calculateDiffHelper(beforeField.Interface(), afterField.Interface(), fieldTag, diffFields)
			default:
				diffFields[fieldTag] = fmt.Sprint(afterFieldInterface)
			}
		}
	}

	// Iterate through after fields.
	afterFieldsCount := afterValue.NumField()
afterFieldLoop:
	for i := 0; i < afterFieldsCount; i++ {
		afterField := afterValue.Field(i)
	indirectAfterDeleteLoop:
		for {
			if afterField.Kind() == reflect.Ptr {
				if afterField.IsNil() {
					continue afterFieldLoop
				}
				afterField = afterField.Elem()
				continue indirectAfterDeleteLoop
			}
			break
		}
		switch afterField.Kind() {
		case reflect.Func:
			continue afterFieldLoop
		}
		afterFieldType := afterType.Field(i)
		fieldTag, ok := afterFieldType.Tag.Lookup("api")
		if !ok {
			continue afterFieldLoop
		}
		commaPos := strings.Index(fieldTag, ",")
		if commaPos > -1 {
			fieldTag = fieldTag[0:commaPos]
		}
		if len(prefix) > 0 {
			fieldTag = prefix + accessorCharacter + fieldTag
		}
		processedFields[fieldTag] = struct{}{}

		// Handle added fields.
		if _, ok := processedFields[fieldTag]; !ok {
			if afterField.CanInterface() && !reflect.DeepEqual(afterField, reflect.Zero(afterFieldType.Type)) {
				afterFieldInterface := afterField.Interface()
				diffFields[fieldTag] = fmt.Sprint(afterFieldInterface)
			}
		}
	}
}
