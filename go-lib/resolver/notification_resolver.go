package resolver

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
)

type notificationResolver struct{ *resolver }

func (r notificationResolver) TriggeredBy(ctx context.Context, obj *model.Notification) (*model.User, error) {
	if obj.TriggeredBy == "" {
		return nil, nil
	}
	return dal.GetUser(ctx, r.db, obj.TriggeredBy)
}
