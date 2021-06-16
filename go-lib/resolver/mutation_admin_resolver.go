package resolver

import (
	"context"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model/dal"

	"github.com/robert-zaremba/errstack"

	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/model"
)

// AdminApproveUser approves user to use the platform.
// If status == rejected then reason is required, otherwise it's ignored.
func (r mutationResolver) AdminApproveUser(ctx context.Context, id string, status model.SimpleApproval, reason *string) (*model.AccessApproval, error) {
	errb := errstack.NewBuilder()
	au, errs := middleware.GetAuthUser(ctx)
	errb.Put("Authentication", errs)
	if !au.IsModerator() {
		errb.Put("Admin", "Admin role required")
	}
	if status == model.SimpleApprovalRejected && *reason == "" {
		errb.Put("ReasonError", "Reason for reject is required")
	}
	u, errs := dal.GetUser(ctx, r.db, id)
	errb.Put("GetUser", errs)
	if errs = errb.ToReqErr(); errs != nil {
		return nil, errs
	}

	approval := model.AccessApproval{
		Status:    status,
		Reason:    reason,
		Approver:  au.ID,
		CreatedAt: time.Now().UTC(),
	}
	u.Approvals = append(u.Approvals, approval)
	return &approval, dal.ReplaceUser(ctx, r.db, u)
}
