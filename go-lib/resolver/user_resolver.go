package resolver

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
)

type userResolver struct{ *resolver }

// OrgMap resolves orgMap field
func (r userResolver) OrgMap(ctx context.Context, u *model.User) ([]model.UserOrgMap, error) {
	return dal.GetOrgMap(ctx, r.db, u.ID)
}

// PubKey returns user default public key
func (r userResolver) PubKey(ctx context.Context, u *model.User) (*string, error) {
	// TODO: this endpoint does not include tradeID. It is deprecated.
	// Front-end should use the `Query.PubKey` endpoint for trade key retrieval.
	authUser, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return nil, errs
	}
	if u.ID != authUser.ID {
		// Only owner should see his own key.
		// Returning an error will result in many meaningless errors.
		return nil, nil
	}
	// HACK: public key should be loaded for each trade specifically.
	staticWallet, ok := u.StaticWallets[u.DefaultWalletID]
	if !ok {
		// User doesn't have a static wallet
		// TODO: Derive the key
		return nil, nil
	}
	k := staticWallet.PubKey
	return &k, nil
}
