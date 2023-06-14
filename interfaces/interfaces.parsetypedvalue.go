package interfaces

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/araddon/dateparse"
	utilstrings "github.com/bertjohnson/util/strings"
)

var (
	falseValue = reflect.ValueOf(false)
	trueValue  = reflect.ValueOf(true)
)

// ParseTypedValue attempts to parse the underlying type of a string.
func ParseTypedValue(dataString string) reflect.Value {
	// Handle empty string.
	if dataString == "" {
		return reflect.Value{}
	}

	// Track which data types are supported.
	complexSupported := true
	numberSupported := true

	// Attempt to handle German number formatting (12.345.678,90).
	originalDataString := dataString
	lastPeriodPos := strings.LastIndex(dataString, ".")
	lastCommaPos := strings.LastIndex(dataString, ",")
	if lastPeriodPos > -1 && lastCommaPos > -1 && lastCommaPos > lastPeriodPos {
		temp := strings.Replace(dataString, ",", "!", -1)
		temp = strings.Replace(temp, ".", ",", -1)
		dataString = strings.Replace(temp, "!", ".", -1)
	} else {
		if strings.Count(dataString, ".") > 1 {
			dataString = strings.Replace(dataString, ".", "", -1)
		}
	}

	// Handle cultures where '-' or '+' is the last character.
	dataStringLength := len(dataString)
	lastCharacter := dataString[dataStringLength-1 : dataStringLength]
	if lastCharacter == "-" || lastCharacter == "+" {
		dataString = lastCharacter + dataString[0:dataStringLength-1]
	}

	// Check each character.
	reader := bufio.NewReader(strings.NewReader(dataString))
	outputBuffer := bytes.Buffer{}
	decimalCount := 0
	hasI := false
	hasOperand := false
	leadingZeroes := false
	pos := 0
runeLoop:
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return reflect.Value{}
			}
		} else {
			// Ignore whitespace.
			if utilstrings.IsWhitespace(c) {
				continue runeLoop
			}

			// Parse characters.
			switch c {
			case '(', ')', '[', ']', '{', '}', ' ', '\t', '\r', '\n':
			case '-', '+':
				if c == '-' || pos > 0 {
					outputBuffer.WriteRune(c) // nolint
				}
				if pos > 0 {
					numberSupported = false
				} else {
					if hasOperand {
						complexSupported = false
						break runeLoop
					}
					hasOperand = true
				}
			case '.':
				if leadingZeroes {
					leadingZeroes = false
				}
				outputBuffer.WriteRune(c) // nolint
				if complexSupported && decimalCount == 2 {
					complexSupported = false
					numberSupported = false
					break runeLoop
				}
				if numberSupported && decimalCount == 1 {
					numberSupported = false
				}
				decimalCount++
			case ',':
			case '0':
				if outputBuffer.Len() == 0 {
					leadingZeroes = true
				}
				outputBuffer.WriteRune(c) // nolint
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				outputBuffer.WriteRune(c) // nolint
			case 'i':
				outputBuffer.WriteRune(c) // nolint
				numberSupported = false
				if hasI {
					complexSupported = false
					break runeLoop
				}
				hasI = true
			default:
				outputBuffer.WriteRune(c) // nolint
				complexSupported = false
				numberSupported = false

				// Check first five characters in case this starts with "true" or "false".
				if pos == 5 {
					break runeLoop
				}
			}
			pos++
		}
	}

	if pos == 0 {
		return reflect.Value{}
	}

	// Parse complex numbers.
	if complexSupported {
		parts := strings.Split(outputBuffer.String(), "+")
		partsLength := len(parts)
		if partsLength == 2 {
			if strings.HasSuffix(parts[0], "i") {
				if realFloat, err := strconv.ParseFloat(parts[1], 64); err == nil {
					if imaginaryFloat, err := strconv.ParseFloat(strings.Trim(parts[0], "i"), 64); err == nil {
						complexVal := complex(realFloat, imaginaryFloat)
						return reflect.ValueOf(complexVal)
					}
				}
			} else {
				if strings.HasSuffix(parts[1], "i") {
					if realFloat, err := strconv.ParseFloat(parts[0], 64); err == nil {
						if imaginaryFloat, err := strconv.ParseFloat(strings.Trim(parts[1], "i"), 64); err == nil {
							complexVal := complex(realFloat, imaginaryFloat)
							return reflect.ValueOf(complexVal)
						}
					}
				}
			}
		}

		parts = strings.Split(outputBuffer.String(), "-")
		partsLength = len(parts)
		if partsLength == 2 {
			if strings.HasSuffix(parts[0], "i") {
				if realFloat, err := strconv.ParseFloat(parts[1], 64); err == nil {
					if imaginaryFloat, err := strconv.ParseFloat(strings.Trim(parts[0], "i"), 64); err == nil {
						complexVal := complex(-1*realFloat, imaginaryFloat)
						return reflect.ValueOf(complexVal)
					}
				}
			} else {
				if strings.HasSuffix(parts[1], "i") {
					if realFloat, err := strconv.ParseFloat(parts[0], 64); err == nil {
						if imaginaryFloat, err := strconv.ParseFloat(strings.Trim(parts[1], "i"), 64); err == nil {
							complexVal := complex(realFloat, -1*imaginaryFloat)
							return reflect.ValueOf(complexVal)
						}
					}
				}
			}
		}
	}

	// Parse numbers.
	if leadingZeroes {
		numberSupported = false
	}
	if numberSupported {
		if decimalCount > 0 {
			if floatVal, err := strconv.ParseFloat(outputBuffer.String(), 64); err == nil {
				return reflect.ValueOf(floatVal)
			}
		}

		if intVal, err := strconv.ParseInt(outputBuffer.String(), 10, 64); err == nil {
			return reflect.ValueOf(intVal)
		}
	}

	// Check for booleans.
	if outputBuffer.Len() < 6 {
		switch strings.ToLower(outputBuffer.String()) {
		case "false":
			return falseValue
		case "true":
			return trueValue
		}
	}

	// Check for string dates/times.
	if len(originalDataString) < 32 {
		if timeVal, err := dateparse.ParseAny(originalDataString); err == nil {
			if timeVal.Year() != 0 {
				return reflect.ValueOf(timeVal.UnixNano() / 1000000)
			}
		}
	}

	// Default to a string.
	return reflect.ValueOf(originalDataString)
}
