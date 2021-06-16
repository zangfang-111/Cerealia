package resolver

import (
	context "context"

	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/utils"
	"github.com/robert-zaremba/errstack"
)

type tradeResolver struct{ *resolver }

func (r tradeResolver) Template(ctx context.Context, obj *model.Trade) (*model.TradeTemplate, error) {
	return dal.GetTradeTempate(ctx, r.db, obj.TemplateID)
}

func (r tradeResolver) Buyer(ctx context.Context, obj *model.Trade) (*model.User, error) {
	return dal.GetUser(ctx, r.db, obj.Buyer.UserID)
}

func (r tradeResolver) Seller(ctx context.Context, obj *model.Trade) (*model.User, error) {
	return dal.GetUser(ctx, r.db, obj.Seller.UserID)
}

func (r tradeResolver) CreatedBy(ctx context.Context, obj *model.Trade) (*model.User, error) {
	return dal.GetUser(ctx, r.db, obj.CreatedBy)
}

func (r tradeResolver) ScAddr(ctx context.Context, obj *model.Trade) (string, error) {
	return string(obj.SCAddr), nil
}

func (r tradeResolver) TradeOffer(ctx context.Context, obj *model.Trade) (*model.TradeOffer, error) {
	if utils.PointerToString(obj.TradeOfferID) == "" {
		return nil, nil
	}
	return dal.GetTradeOffer(ctx, r.db, *obj.TradeOfferID)
}

type tradeStageAddReqResolver struct{ *resolver }

func (r tradeStageAddReqResolver) ReqBy(ctx context.Context, obj *model.TradeStageAddReq) (*model.User, error) {
	return dal.GetUser(ctx, r.db, obj.ReqBy)
}

func (r tradeStageAddReqResolver) ApprovedBy(ctx context.Context, obj *model.TradeStageAddReq) (*model.User, error) {
	if obj.NilOrPending() {
		return nil, nil
	}
	u, errs := dal.GetUser(ctx, r.db, obj.ApprovedBy)
	return u, model.ResetIfErrNoID(errs)
}

type tradeStageDocResolver struct{ *resolver }

func (r tradeStageDocResolver) ApprovedBy(ctx context.Context, obj *model.TradeStageDoc) (*model.User, error) {
	if obj.NilOrPending() {
		return nil, nil
	}
	u, errs := dal.GetUser(ctx, r.db, obj.ApprovedBy)
	return u, model.ResetIfErrNoID(errs)
}

func (r tradeStageDocResolver) Doc(ctx context.Context, obj *model.TradeStageDoc) (*model.Doc, error) {
	d, errs := dal.GetDoc(ctx, r.db, obj.DocID)
	return d, model.ResetIfErrNoID(errs)
}

type stageModeratorResolver struct{ *resolver }

func (r stageModeratorResolver) User(ctx context.Context, obj *model.StageModerator) (*model.User, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	if !u.IsModerator() || obj == nil {
		return nil, nil
	}
	foundUser, errs := dal.GetUser(ctx, r.db, obj.UserID)
	return foundUser, model.ResetIfErrNoID(errs)
}

// ActorData retrieves current user's public key and other info of a trade
func (r tradeResolver) ActorWallet(ctx context.Context, t *model.Trade) (*model.TradeActorWallet, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	party, err := t.FindParticipant(u)
	if err != nil {
		return nil, errstack.WrapAsReqF(err, "User '%s' is not a participant of trade '%s'", u.ID, t.ID)
	}
	return &model.TradeActorWallet{
		PubKey:   party.PubKey,
		KeyPath:  party.KeyDerivationPath,
		WalletID: party.WalletID,
	}, nil
}
