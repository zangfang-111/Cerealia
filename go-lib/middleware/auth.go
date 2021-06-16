package middleware

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/auth"
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	driver "github.com/arangodb/go-driver"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/errstack"
)

type contextKey struct {
	name string
}

// ctxUserKey is a key used to store UserID in a context
var ctxUserKey = &contextKey{"user-id"}

// WithAuth attaches userID in the context
func WithAuth(db driver.Database) routing.Handler {
	return func(c *routing.Context) error {
		tokenStr, _ := auth.TokenFromAuthHeader(c.Request)
		if tokenStr == "" {
			return nil
		}
		userid, _ := auth.Authorize(tokenStr)
		if userid == "" {
			return nil
		}
		ctx := c.Request.Context()
		u, errs := dal.GetUser(ctx, db, userid)
		if errs != nil {
			logger.Error("Get User Failed", errs)
		}
		ctx = context.WithValue(ctx, ctxUserKey, u)
		c.Request = c.Request.WithContext(ctx)
		return nil
	}
}

// GetAuthUser returns User from the current request context.
// If there is no authenticated user returns nil and unauthenticated error.
func GetAuthUser(ctx context.Context) (*model.User, errstack.E) {
	raw := ctx.Value(ctxUserKey)
	if raw == nil {
		return nil, model.ErrUnauthenticated
	}
	u, _ := raw.(*model.User)
	return u, nil
}
