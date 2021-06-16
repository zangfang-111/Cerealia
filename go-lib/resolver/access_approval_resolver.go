package resolver

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
)

type accessApprovalResolver struct{ *resolver }

func (r accessApprovalResolver) Approver(ctx context.Context, obj *model.AccessApproval) (*model.User, error) {
	return dal.GetUser(ctx, r.db, obj.Approver)
}
