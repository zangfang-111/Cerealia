package dal

import (
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

// notFound wraps errstack.E into not found error
type notFound struct {
	errstack.E
}

// NewNotFound returns new not found error
func NewNotFound(msg string) errstack.E {
	if msg == "" {
		msg = "object not found"
	}
	return notFound{errstack.NewReq(msg)}
}

// MaybeWrapAsNotFound tries to err into not found error if it makes sense,
// meaning: the error is driver.NotFound or it's a request error.
func MaybeWrapAsNotFound(err error) errstack.E {
	if err == nil {
		return nil
	}
	if driver.IsNotFound(err) || driver.IsNoMoreDocuments(err) {
		return notFound{errstack.WrapAsReq(err, "DB: object not found")}
	}
	if driver.IsArangoError(err) {
		return errstack.WrapAsInf(err, "DB: can't get an object")
	}
	return errstack.WrapAsDomain(err, "Object not found")
}

// IsNotFound check if an error is notFound or driver.NotFound error type
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if driver.IsNoMoreDocuments(err) {
		return true
	}
	_, ok := err.(notFound)
	return ok
}

// WrapReadDomainError check the error when reading document and wraps it into
// errstack.E Domain error. `objectType` is the expected object type.
func WrapReadDomainError(err error, objectType string) errstack.E {
	return errstack.WrapAsDomain(err, "Can't decode document into"+objectType+" object")
}
