package model

import (
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	"github.com/robert-zaremba/errstack"
)

var (
	// ErrUnauthenticated is thrown when user login is required get the resources
	ErrUnauthenticated = errstack.NewReq("Authentication required")
	// ErrUnauthorized is thrown when user doesn't have permission to get resources
	ErrUnauthorized = errstack.NewReq("Access denied")
	// ErrNoID is the NoneID error
	ErrNoID = errstack.NewDomain("The provided ID is empty")
)

// ErrDbCollection returns fromated error message during connection of db collections
func ErrDbCollection(err error, col dbconst.Col) errstack.E {
	return errstack.WrapAsInf(err, "DB: Can't connect the collection "+string(col))
}

// ResetIfErrNoID returns nil if err == ErrNoID
func ResetIfErrNoID(err error) error {
	if err == ErrNoID || errstack.Cause(err) == ErrNoID {
		return nil
	}
	return err
}
