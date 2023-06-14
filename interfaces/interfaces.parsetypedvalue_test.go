package interfaces

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestParseTypedValue tests ParseTypedValue().
func TestParseTypedValue(t *testing.T) {
	var val, expectedValue reflect.Value
	var err error

	// Positive Integer.
	val = ParseTypedValue("+1 234,567,890")
	expectedValue = reflect.ValueOf(int64(1234567890))
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Negative Integer.
	val = ParseTypedValue("-1234567890")
	expectedValue = reflect.ValueOf(int64(-1234567890))
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// German Integer.
	val = ParseTypedValue("1.234.567.890")
	expectedValue = reflect.ValueOf(int64(1234567890))
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Integer ending with sign.
	val = ParseTypedValue("1234567890-")
	expectedValue = reflect.ValueOf(int64(-1234567890))
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Float.
	val = ParseTypedValue("12,345,678.90")
	expectedValue = reflect.ValueOf(12345678.90)
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// German Float.
	val = ParseTypedValue("12.345.678,90")
	expectedValue = reflect.ValueOf(12345678.90)
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Complex.
	val = ParseTypedValue("12,345.678 + 9i")
	expectedValue = reflect.ValueOf(complex(12345.678, 9))
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Leading Zero to String.
	val = ParseTypedValue("012345")
	expectedValue = reflect.ValueOf("012345")
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Leading Zero to Float.
	val = ParseTypedValue("0.12345")
	expectedValue = reflect.ValueOf(.12345)
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// String.
	val = ParseTypedValue("Last, First")
	expectedValue = reflect.ValueOf("Last, First")
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Boolean.
	val = ParseTypedValue("true")
	expectedValue = reflect.ValueOf(true)
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Current Time.
	timeVal := time.Now().Round(time.Second)
	val = ParseTypedValue(timeVal.Format(time.RFC1123Z))
	expectedValue = reflect.ValueOf(timeVal.UnixNano() / 1000000)
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Time.
	val = ParseTypedValue("1524866400000")
	timeVal, err = time.Parse("Monday, January 02, 2006 15:04:05 PM", "Friday, April 27, 2018 10:00:00 PM")
	expectedValue = reflect.ValueOf(timeVal.UnixNano() / 1000000)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue.Interface(), val.Interface())

	// Empty.
	val = ParseTypedValue("")
	expectedValue = (reflect.Value{})
	assert.Equal(t, expectedValue, val)
}
