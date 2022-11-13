package rangeValues

import (
	"strconv"
	"strings"
)

func ParseIntContainer[T Integer](value string) *IntContainer[T] {
	parts := strings.Split(value, ":")
	if len(parts) < 1 {
		return nil
	}

	intValue, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil
	}

	return &IntContainer[T]{
		Value: T(intValue),
		Flags: parts[1:],
	}
}
