package policycondition

import (
	"strconv"
	"strings"
)

func ParseAge(value string) *int {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil
	}

	value = strings.TrimSuffix(value, "세")

	age, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}

	return &age
}
