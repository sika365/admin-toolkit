package utils

import "reflect"

// DefaultIfZero returns defaultVal if val is the zero value for its type
func DefaultIfZero[T comparable](val, defaultVal T) T {
	if reflect.ValueOf(val).IsZero() {
		return defaultVal
	}
	return val
}
