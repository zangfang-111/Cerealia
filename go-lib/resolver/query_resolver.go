package resolver

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"github.com/robert-zaremba/errstack"
)

type queryResolver struct{ *resolver }

// User returns current user if id==nil or user with given id
func (r queryResolver) User(ctx context.Context, id *string) (*model.User, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	if id == nil {
		return u, nil
	}
	return dal.GetUser(ctx, r.db, *id)
}

func (r queryResolver) Users(ctx context.Context) ([]model.User, error) {
	if _, err := middleware.GetAuthUser(ctx); err != nil {
		return nil, err
	}
	return dal.GetApprovedUsers(ctx, r.db)
}

func (r queryResolver) AdminUsers(ctx context.Context) ([]model.AdminUser, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	if !u.IsModerator() {
		return nil, errstack.NewReq("Admin role required")
	}
	return dal.GetAdminUsers(ctx, r.db)
}

func (r queryResolver) TradeTemplates(ctx context.Context) ([]model.TradeTemplate, error) {
	return dal.GetTradeTemplates(ctx, r.db)
}

func (r queryResolver) Trade(ctx context.Context, id string) (*model.Trade, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	t, errs := dal.GetTrade(ctx, r.db, id)
	if errs != nil && t.Buyer.UserID != u.ID && t.Seller.UserID != u.ID {
		return nil, model.ErrUnauthorized
	}
	return t, errs
}

func (r queryResolver) Trades(ctx context.Context) ([]model.Trade, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	return dal.GetTrades(ctx, r.db, u.ID)
}

func (r queryResolver) TradeOffer(ctx context.Context, id string) (*model.TradeOffer, error) {
	if _, err := middleware.GetAuthUser(ctx); err != nil {
		return nil, err
	}
	return dal.GetTradeOffer(ctx, r.db, id)
}

func (r queryResolver) TradeOffers(ctx context.Context) ([]model.TradeOffer, error) {
	if _, err := middleware.GetAuthUser(ctx); err != nil {
		return nil, err
	}
	return dal.GetTradeOffers(ctx, r.db)
}

func (r queryResolver) Notifications(ctx context.Context, from uint) ([]model.Notification, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	return dal.GetUserNotifacations(ctx, r.db, u.ID, from)
}
func (r queryResolver) NotificationsTrade(ctx context.Context, id string) ([]model.Notification, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	return dal.GetNotifacationsByTrade(ctx, r.db, u.ID, id)
}

func (r queryResolver) StellarNet(ctx context.Context) (*model.StellarNet, error) {
	n := r.stellarDriver.Network
	stInfo := model.StellarNet{
		Name:       n.Name,
		URL:        n.URL,
		Passphrase: n.Passphrase.Passphrase}
	return &stInfo, nil
}

func (r queryResolver) AdminTrades(ctx context.Context) ([]model.Trade, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	if !u.IsModerator() {
		return nil, model.ErrUnauthorized
	}
	return dal.GetAllTrades(ctx, r.db)
}

// PubKey retrieves current user's public key for a trade
func (r queryResolver) PubKey(ctx context.Context, tradeID string) (*string, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	trade, err := dal.GetTrade(ctx, r.db, tradeID)
	if err != nil {
		return nil, errstack.WrapAsReqF(err, "Trade not found: %s", tradeID)
	}
	addr, err := trade.FindPubKey(u)
	if err != nil {
		return nil, errstack.WrapAsReqF(err, "Can't get user's '%s' key for trade '%s'", u.ID, trade.ID)
	}
	key := string(addr)
	return &key, nil
}

func (r queryResolver) Organizations(ctx context.Context) ([]model.Organization, error) {
	_, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	return dal.GetAllOrganizations(ctx, r.db)
}
