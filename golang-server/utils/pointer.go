package utils

import (
	"strconv"
)

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

func SetIntValue(value string, param **int) {
	if value != "" {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			*param = &intValue
		}
	}
}

func SetFloatValue(value string, param **float64) {
	if value != "" {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			*param = &floatValue
		}
	}
}

func SetBoolValue(value string, param **bool) {
	if value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			*param = &boolValue
		}
	}
}

func SetStringValue(value string, param **string) {
	if value != "" {
		*param = &value
	}
}
