package resolver

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
)

type tradeOfferRes struct{ *resolver }

func (r tradeOfferRes) CreatedBy(ctx context.Context, obj *model.TradeOffer) (*model.User, error) {
	return dal.GetUser(ctx, r.db, obj.CreatedBy)
}

func (r tradeOfferRes) Org(ctx context.Context, obj *model.TradeOffer) (*model.Organization, error) {
	if obj.IsAnonymous {
		return nil, nil
	}
	return dal.GetOrganization(ctx, r.db, obj.OrgID)
}

func (r tradeOfferRes) Terms(ctx context.Context, obj *model.TradeOffer) (*model.Doc, error) {
	if obj.DocID == nil {
		return nil, nil
	}
	return dal.GetDoc(ctx, r.db, *obj.DocID)
}
