package validation

import (
	"github.com/robert-zaremba/errstack"
)

// NotEmpty validates if the given string is not empty
func NotEmpty(val string, errp errstack.Putter) {
	if val == "" {
		errp.Put("can't be empty")
	}
}

// Positive validates if the given uint is more than 0
func Positive(val uint, errp errstack.Putter) {
	if val <= 0 {
		errp.Put("should be positive")
	}
}
