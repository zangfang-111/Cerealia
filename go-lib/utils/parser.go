package utils

import (
	"strconv"
)

// PointerToString returns string value from string pointer
func PointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// UintToString returns string value from uint
func UintToString(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}
