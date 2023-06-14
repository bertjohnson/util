package interfaces

import (
	"errors"
	"reflect"
)

// Inherit inherits values from each ancestor if the output values are their defaults.
func Inherit(output interface{}, ancestors []interface{}) error {
	return InheritWithTag(output, ancestors, "")
}

// InheritWithTag inherits values from each ancestor if the output values are their defaults, only if the tag exists.
func InheritWithTag(output interface{}, ancestors []interface{}, tagName string) error {
	// If there are no ancestors, return.
	if len(ancestors) == 0 {
		return nil
	}

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

	// Unpack each ancestor.
	var ancestorValues []reflect.Value
	for _, ancestor := range ancestors {
		ancestorValue, ok := ancestor.(reflect.Value)
		if !ok {
			ancestorValue = reflect.ValueOf(ancestor)
			if ancestorValue.Kind() == reflect.Ptr {
				ancestorValue = ancestorValue.Elem()
			}
		}
		ancestorValue = reflect.Indirect(ancestorValue)

		// If the passed in object is an interface, unpack its struct.
		if ancestorValue.Kind() == reflect.Interface {
			ancestorValue = ancestorValue.Elem()
		}

		ancestorValues = append(ancestorValues, ancestorValue)
	}

	// Loop through each field.
	fieldsLength := outputType.NumField()
fieldLoop:
	for i := 0; i < fieldsLength; i++ {
		field := outputType.Field(i)
		outputFieldValue := outputValue.Field(i)

		// Only process exported fields.
		if outputFieldValue.CanSet() {
			// Check for tag match.
			inheritBool := false
			boolStrategy := "or"
			if tagName != "" {
				_, ok := field.Tag.Lookup(tagName)
				if !ok {
					continue fieldLoop
				}
				boolStrategy, ok = field.Tag.Lookup(tagName + "bool")
				if ok {
					inheritBool = true
				}
			}

			fieldKind := field.Type.Kind()
			switch fieldKind {
			case reflect.Array, reflect.Slice:
				// Copy all array values, regardless of whether output was empty.
				indirect := reflect.Indirect(outputFieldValue)
				for _, ancestorValue := range ancestorValues {
					ancestorFieldValue := ancestorValue.FieldByName(field.Name)
					if ancestorFieldValue.Kind() == fieldKind {
						ancestorFieldValueLength := ancestorFieldValue.Len()
						for j := 0; j < ancestorFieldValueLength; j++ {
							outputFieldValue.Set(reflect.Append(indirect, ancestorFieldValue.Index(j)))
						}
					}
				}
			case reflect.Bool:
				if inheritBool {
					// Implement boolean strategy.
					boolVal := outputFieldValue.Bool()
					if (boolVal && boolStrategy == "or") || (!boolVal && boolStrategy == "and") {
						continue
					}

				boolAncestorsLoop:
					// Process ancestor values.
					for _, ancestorValue := range ancestorValues {
						ancestorFieldValue := ancestorValue.FieldByName(field.Name)
						if ancestorFieldValue.Kind() == reflect.Bool {
							ancestorBool := ancestorFieldValue.Bool()
							if boolStrategy == "and" {
								if !ancestorBool {
									outputFieldValue.SetBool(false)
									break boolAncestorsLoop
								}
							} else {
								if ancestorBool {
									outputFieldValue.SetBool(true)
									break boolAncestorsLoop
								}
							}
						}
					}
				}
			case reflect.Chan, reflect.Func:
				// Ignore channels and functions.
				continue
			case reflect.Complex128, reflect.Complex64:
				if outputFieldValue.Complex() == 0 {
				complex128AncestorsLoop:
					for _, ancestorValue := range ancestorValues {
						ancestorFieldValue := ancestorValue.FieldByName(field.Name)
						if ancestorFieldValue.Complex() != 0 {
							outputFieldValue.SetComplex(ancestorFieldValue.Complex())
							break complex128AncestorsLoop
						}
					}
				}
			case reflect.Float32, reflect.Float64:
				if outputFieldValue.Float() == 0 {
				float32AncestorsLoop:
					for _, ancestorValue := range ancestorValues {
						ancestorFieldValue := ancestorValue.FieldByName(field.Name)
						if ancestorFieldValue.Float() != 0 {
							outputFieldValue.SetFloat(ancestorFieldValue.Float())
							break float32AncestorsLoop
						}
					}
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if outputFieldValue.Int() == 0 {
				intAncestorsLoop:
					for _, ancestorValue := range ancestorValues {
						ancestorFieldValue := reflect.Indirect(ancestorValue.FieldByName(field.Name))
						if ancestorFieldValue.Int() != 0 {
							outputFieldValue.SetInt(ancestorFieldValue.Int())
							break intAncestorsLoop
						}
					}
				}
			case reflect.Map:
				// If the map pointer is nil, create it.
				if outputFieldValue.IsNil() {
					outputFieldValue.Set(reflect.MakeMap(outputFieldValue.Type()))
				}

				// Copy all map values, regardless of whether output was empty.
				for _, ancestorValue := range ancestorValues {
					ancestorFieldValue := ancestorValue.FieldByName(field.Name)
					if ancestorFieldValue.Kind() == reflect.Map {
						for _, key := range ancestorFieldValue.MapKeys() {
							outputFieldValue.SetMapIndex(key, ancestorFieldValue.MapIndex(key))
						}
					}
				}
			case reflect.Ptr:
				if outputFieldValue.IsNil() {
				pointerAncestorsLoop:
					for _, ancestorValue := range ancestorValues {
						ancestorFieldValue := ancestorValue.FieldByName(field.Name)
						if ancestorFieldValue.Kind() == reflect.Ptr {
							if !ancestorFieldValue.IsNil() {
								outputFieldPointer := reflect.New(ancestorFieldValue.Elem().Type())
								outputFieldPointer.Elem().Set(ancestorFieldValue.Elem())
								outputFieldValue.Set(outputFieldPointer)
								break pointerAncestorsLoop
							}
						}
					}
				}
			case reflect.String:
				if outputFieldValue.String() == "" {
				stringAncestorsLoop:
					for _, ancestorValue := range ancestorValues {
						ancestorFieldValue := ancestorValue.FieldByName(field.Name)
						if ancestorFieldValue.Kind() == reflect.String {
							if ancestorFieldValue.Len() > 0 {
								outputFieldValue.SetString(ancestorFieldValue.String())
								break stringAncestorsLoop
							}
						}
					}
				}
			case reflect.Struct:
				switch field.Type.Name() {
				case timeString:
					// Special case to handle times.
				timeAncestorsLoop:
					for _, ancestorValue := range ancestorValues {
						ancestorFieldValue := ancestorValue.FieldByName(field.Name)
						if ancestorFieldValue.Kind() == reflect.Struct {
							outputFieldValue.Set(ancestorFieldValue)
							break timeAncestorsLoop
						}
					}
				default:
					// Handle sub-structure recursively.
					ancestorFieldValues := []interface{}{}
					for _, ancestorValue := range ancestorValues {
						ancestorFieldValue := ancestorValue.FieldByName(field.Name)
						if ancestorFieldValue.Kind() == reflect.Struct {
							ancestorFieldValues = append(ancestorFieldValues, ancestorFieldValue)
						}
					}
					if len(ancestorFieldValues) > 0 {
						err := InheritWithTag(outputFieldValue, ancestorFieldValues, tagName)
						if err != nil {
							return err
						}
					}
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if outputFieldValue.Uint() == 0 {
				uintAncestorsLoop:
					for _, ancestorValue := range ancestorValues {
						ancestorFieldValue := ancestorValue.FieldByName(field.Name)
						if ancestorFieldValue.Uint() != 0 {
							outputFieldValue.SetUint(ancestorFieldValue.Uint())
							break uintAncestorsLoop
						}
					}
				}
			}
		}
	}

	return nil
}
