package middleware

import (
	"context"
	"encoding/json"

	"github.com/99designs/gqlgen/graphql"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/errstack"
	bat "github.com/robert-zaremba/go-bat"
	"github.com/vektah/gqlparser/gqlerror"
)

func restError(c *routing.Context) error {
	err := c.Next()
	logError(err)
	if e, ok := err.(errorJSONMarshaler); err != nil && ok {
		return errJSONDisplay{e}
	}
	return err
}

// GraphQLError handles errors for GraphQL
func GraphQLError(ctx context.Context, e error) *gqlerror.Error {
	logError(e)
	if gqlerr, ok := e.(*gqlerror.Error); ok {
		gqlerr.Path = graphql.GetResolverContext(ctx).Path()
		return gqlerr
	}
	status := statusFromError(e)
	if status >= 500 {
		return &gqlerror.Error{
			Message:    "Internal Server Error",
			Path:       graphql.GetResolverContext(ctx).Path(),
			Extensions: map[string]interface{}{"status": status},
		}
	}
	return graphql.DefaultErrorPresenter(ctx, e)
}

type hasLog interface {
	Log()
}

func logError(err error) {
	if err == nil {
		return
	}
	if e, ok := err.(hasLog); ok {
		e.Log()
		return
	}
	status := statusFromError(err)
	if status < 400 {
		return
	}
	if status < 500 {
		logger.Debug("Request error", "status", status, err)
	} else {
		logger.Error("Server error", "status", status, err)
	}
}

// it expects err != nil
func statusFromError(err error) int {
	switch e := err.(type) {
	case errstack.HasStatusCode:
		return e.StatusCode()
	case errstack.HasUnderlying:
		return statusFromError(e.Cause())
	}
	return 500
}

type errorJSONMarshaler interface {
	json.Marshaler
	error
}

type errJSONDisplay struct {
	err errorJSONMarshaler
}

func (e errJSONDisplay) Error() string {
	j, err := e.err.MarshalJSON()
	if err != nil {
		logger.Error("can't marshal horizon error", err)
		return e.err.Error()
	}
	return bat.UnsafeByteArrayToStr(j)
}
