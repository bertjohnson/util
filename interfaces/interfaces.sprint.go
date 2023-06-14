package interfaces

import (
	"fmt"
)

// Sprint returns fmt.Sprint() unless the input is nil.
func Sprint(input interface{}) string {
	if input == nil {
		return ""
	}

	return fmt.Sprint(input)
}
