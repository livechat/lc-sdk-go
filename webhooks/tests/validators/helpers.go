package webhooks_validators

import (
	"fmt"
)

func PropEq(propertyName string, actual, expected interface{}, validationAccumulator *string) {
	if actual != expected {
		*validationAccumulator += fmt.Sprintf("%s mismatch, actual: %s, expected: %s\n", propertyName, actual, expected)
	}
}