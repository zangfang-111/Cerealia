package validation

import (
	"github.com/go-openapi/validate"
	"github.com/robert-zaremba/errstack"
)

// Required validates that the src value is not empty
func Required(data interface{}, errp errstack.Putter) {
	if validate.Required("path", "dummy", data) != nil {
		errp.Put("is required")
	}
}
