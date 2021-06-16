package resolver

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
)

type approveReqResolver struct{ *resolver }

func (r approveReqResolver) ReqBy(ctx context.Context, obj *model.ApproveReq) (*model.User, error) {
	u, errs := dal.GetUser(ctx, r.db, obj.ReqBy)
	return u, errs
}

func (r approveReqResolver) ApprovedBy(ctx context.Context, obj *model.ApproveReq) (*model.User, error) {
	if obj.NilOrPending() {
		return nil, nil
	}
	return dal.GetUser(ctx, r.db, obj.ApprovedBy)
}
