package resolver

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
)

type docResolver struct{ *resolver }

func (r docResolver) Hash(ctx context.Context, obj *model.Doc) ([]byte, error) {
	return []byte(obj.Hash), nil
}

func (r docResolver) CreatedBy(ctx context.Context, obj *model.Doc) (*model.User, error) {
	return dal.GetUser(ctx, r.db, obj.CreatedBy)
}
