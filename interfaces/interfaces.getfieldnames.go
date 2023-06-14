package interfaces

import (
	"reflect"
	"strings"
)

// GetAPIFieldNames returns the names of all fields of an object with an "api" tag.
func GetAPIFieldNames(input interface{}) []string {
	fieldNames := []string{}
	getTagFieldNamesHelper(&input, "api", "", &fieldNames)
	return fieldNames
}

// GetFieldNames returns the names of all fields of an object.
func GetFieldNames(input interface{}) []string {
	fieldNames := []string{}
	getTagFieldNamesHelper(&input, "", "", &fieldNames)
	return fieldNames
}

// GetTagFieldNames returns the names of all fields of an object.
func GetTagFieldNames(input interface{}, tag string) []string {
	fieldNames := []string{}
	getTagFieldNamesHelper(&input, tag, "", &fieldNames)
	return fieldNames
}

// getTagFieldNamesHelper returns the names of all fields of an object.
func getTagFieldNamesHelper(input *interface{}, tag string, prefix string, fieldNames *[]string) {
	inputValue := reflect.ValueOf(*input)
	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}
	inputValueType := reflect.TypeOf(*input)
	if inputValueType.Kind() == reflect.Ptr {
		inputValueType = inputValueType.Elem()
	}
	inputValueFieldsLength := inputValue.NumField()
	var fieldTag string
	var ok bool
	for i := 0; i < inputValueFieldsLength; i++ {
		fieldValue := inputValue.Field(i)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}
		fieldType := inputValueType.Field(i)
		if len(tag) > 0 {
			fieldTag, ok = fieldType.Tag.Lookup(tag)
			if ok {
				commaPos := strings.Index(fieldTag, ",")
				if commaPos > -1 {
					fieldTag = fieldTag[0:commaPos]
				}
				if len(prefix) > 0 {
					fieldTag = prefix + "." + fieldTag
				}
				if len(fieldTag) > 0 {
					*fieldNames = append(*fieldNames, fieldTag)
				}
			}
		} else {
			fieldTag = fieldType.Name
			if len(prefix) > 0 {
				fieldTag = prefix + "." + fieldTag
			}
			*fieldNames = append(*fieldNames, fieldTag)
		}

		if fieldValue.Kind() == reflect.Struct && fieldValue.CanInterface() {
			fieldValueObject := fieldValue.Interface()
			getTagFieldNamesHelper(&fieldValueObject, tag, fieldTag, fieldNames)
		}
	}
}
