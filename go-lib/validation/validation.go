// Package validation contains validation functions
package validation

import (
	"time"

	"github.com/go-openapi/validate"
	"github.com/robert-zaremba/errstack"
)

const required = "validation.required"
const insufficientLength = "validation.insufficient-length"
const indexOutOfBounds = "validation.index-out-of-bounds"
const badFormat = "validation.bad-format"
const doesNotMatch = "validation.does-not-match"

type mapping struct {
	Key   string
	Value string
}

// Builder is a validation builder object
type Builder struct {
	Accumulated []mapping
}

// Append appends validation output to the builder
func (vb *Builder) Append(fieldName string, value string) {
	vb.Accumulated = append(vb.Accumulated, mapping{fieldName, value})
}

// Required validates and appends messages to vb
func (vb *Builder) Required(fieldName string, data interface{}) {
	v := validate.Required(fieldName, "dummy", data)
	if v != nil {
		vb.Append(fieldName, required)
	}
}

// MinLength validates and checks if the data has minimum string length
func (vb *Builder) MinLength(fieldName, data string, n int) {
	if n <= 0 {
		return
	}
	length := len(data)
	if length == 0 {
		vb.Append(fieldName, required)
		return
	}
	if length < n {
		vb.Append(fieldName, insufficientLength)
	}
}

// Unique validates for uniqueness
// only adds a message if validation doesn't pass
func (vb *Builder) Unique(fieldName1, fieldName2 string, data1, data2 interface{}, message string) {
	v := validate.UniqueItems("dummy", "dummy", []interface{}{data1, data2})
	if v != nil {
		vb.Append(fieldName1, message)
		vb.Append(fieldName2, message)
	}
}

// IndexLessThan validates index bounds
func (vb *Builder) IndexLessThan(fieldName string, idx, maxValue int) {
	if idx >= maxValue || idx < 0 {
		vb.Append(fieldName, indexOutOfBounds)
	}
}

// Time checks if time value is well formed
// Example: "2006-01-02T15:04:05+07:00"
func (vb *Builder) Time(fieldName string, str string) {
	_, err := time.Parse(time.RFC3339, str)
	if err != nil {
		vb.Append(fieldName, badFormat)
	}
}

// Match checks if values match
func (vb *Builder) Match(fieldName string, str1, str2 interface{}) {
	if str1 != str2 {
		vb.Append(fieldName, doesNotMatch)
	}
}

// IsEmpty returns true when this builder is empty
func (vb *Builder) IsEmpty() bool {
	return len(vb.Accumulated) == 0
}

// ToErrstackBuilder converts ValidatorBuilder to errstack.Builder
// prevents deep coupling with errstack.Builder
func (vb *Builder) ToErrstackBuilder() errstack.Builder {
	errb := errstack.NewBuilder()
	for _, mapping := range vb.Accumulated {
		errb.Put(mapping.Key, mapping.Value)
	}
	return errb
}
