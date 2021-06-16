package resolver

import (
	"context"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

func getTradeRequester(ctx context.Context, db driver.Database, tradeID string) (model.TradeActor, *model.User, *model.Trade, errstack.E) {
	reqActor := model.TradeActorB
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return reqActor, nil, nil, errs
	}
	t, errs := dal.GetTrade(ctx, db, tradeID)
	if errs != nil {
		return reqActor, nil, nil, errs
	}
	errs = t.CanBeModifiedBy(u)
	if errs != nil {
		return reqActor, nil, nil, errs
	}
	reqActor, errs = t.Requester(u)
	return reqActor, u, t, errs
}

func getTradeParticipants(ctx context.Context, db driver.Database, bid, sid string) (model.TradeParticipants, errstack.E) {
	ub, us, errs := dal.Get2Users(ctx, db, bid, sid)
	if errs != nil {
		return model.TradeParticipants{}, errs
	}
	buyer, errs := newTradeParticipant(ctx, db, ub)
	if errs != nil {
		return model.TradeParticipants{}, errs
	}
	seller, errs := newTradeParticipant(ctx, db, us)
	return model.TradeParticipants{
		Buyer:  *buyer,
		Seller: *seller,
	}, errs
}

func prepareTradeStageAddReqApproval(ctx context.Context, db driver.Database, id model.TradeStagePath,
	appendStage bool) (*model.Trade, *model.TradeStageAddReq, *model.User, error) {
	_, u, t, errs := getTradeRequester(ctx, db, id.Tid)
	if errs != nil {
		return nil, nil, nil, errs
	}
	sr, errs := t.GetStageAddReq(id.StageIdx)
	if errs != nil {
		return nil, nil, nil, errs
	}
	now := time.Now().UTC()
	sr.ApprovedBy = u.ID
	sr.ApprovedAt = &now
	if appendStage {
		t.Stages = append(t.Stages,
			model.NewTradeStage(sr.Name, sr.Description, int(id.StageIdx), sr.Owner))
	}
	return t, sr, u, sr.CanBeApproved()
}

func prepareTradeStageReqApproval(ctx context.Context, db driver.Database, id model.TradeStagePath,
	isDel bool) (*model.Trade, *model.ApproveReq, *model.User, error) {
	_, u, t, errs := getTradeRequester(ctx, db, id.Tid)
	if errs != nil {
		return nil, nil, nil, errs
	}
	s, errs := t.GetStage(id.StageIdx)
	if errs != nil {
		return nil, nil, nil, errs
	}
	var req *model.ApproveReq
	if isDel {
		req, errs = s.GetLastDeletionRequest()
	} else {
		req, errs = s.GetLastClosingRequest()
	}
	if errs != nil {
		return nil, nil, nil, errs
	}
	if req.Status != model.ApprovalPending {
		return nil, nil, nil, errstack.NewReq("This stage has no pending approval requests")
	}
	if req.ReqBy == u.ID {
		return nil, nil, nil, errSelfApprove
	}
	if s.IsDeletedOrClosed() {
		return nil, nil, nil, errChangeStage
	}
	req.SetApprovedBy(u.ID)
	return t, req, u, nil
}

func prepareTradeCloseReqApproval(ctx context.Context, db driver.Database, id string) (*model.Trade, *model.ApproveReq, *model.User, error) {
	_, u, t, errs := getTradeRequester(ctx, db, id)
	if errs != nil {
		return nil, nil, nil, errs
	}

	lastReq := &t.CloseReqs[len(t.CloseReqs)-1]
	if lastReq.ReqBy == u.ID {
		return nil, nil, nil, errSelfApprove
	}
	lastReq.SetApprovedBy(u.ID)
	return t, lastReq, u, nil
}

func prepareTradeStageDocApproval(ctx context.Context, db driver.Database,
	id model.TradeStageDocPath) (*model.Trade, *model.TradeStageDoc, *model.User, errstack.E) {
	_, u, t, errs := getTradeRequester(ctx, db, id.Tid)
	if errs != nil {
		return nil, nil, nil, errs
	}
	s, d, errs := t.GetStageDoc(id.StageIdx, id.StageDocIdx)
	if errs != nil {
		return nil, nil, nil, errs
	}
	if errs := d.Status.ShouldBePending(); errs != nil {
		return nil, nil, nil, errs
	}
	if d.ExpiresAt.Before(time.Now()) {
		return nil, nil, nil, errstack.NewReq("This document has expired, you can't approve or reject it")
	}
	if s.IsDeletedOrClosed() {
		return nil, nil, nil, errChangeStage
	}
	doc, errs := dal.GetDoc(ctx, db, d.DocID)
	if errs != nil {
		return nil, nil, nil, errs
	}
	if doc.CreatedBy == u.ID {
		return nil, nil, nil, errSelfApprove
	}
	now := time.Now().UTC()
	d.ApprovedBy = u.ID
	d.ApprovedAt = &now
	return t, d, u, nil
}

func dealCreateStageApproval(t *model.Trade, sr *model.TradeStageAddReq, u *model.User, withApproval bool) model.Approval {
	if withApproval {
		return model.ApprovalPending
	}
	sr.ApproveReq.Status = model.ApprovalNil
	s := model.NewTradeStage(sr.Name, sr.Description, len(t.StageAddReqs), sr.Owner)
	now := time.Now().UTC()
	if t.Buyer.UserID != u.ID && t.Seller.UserID != u.ID && u.IsModerator() {
		s.Moderator = model.StageModerator{
			UserID:    u.ID,
			CreatedAt: &now,
		}
	}
	t.Stages = append(t.Stages, s)
	return model.ApprovalApproved
}

func validateOpTypeAndGetTrade(ctx context.Context, db driver.Database, ld txsource.Driver, tid string, operationType model.Approval) (*model.Trade, *txsource.SourceAccs, errstack.E) {
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return nil, nil, errs
	}
	if !operationType.IsValid() {
		return nil, nil, errstack.NewReq("OperationType is not correct")
	}
	t, err := dal.GetTrade(ctx, db, tid)
	if err != nil {
		return nil, nil, err
	}
	if err = t.CanBeModifiedBy(u); err != nil {
		return nil, nil, err
	}
	sourceAcc, errAcq := ld.Acquire(ctx, t.SCAddr, t.ID, u.ID)
	return t, sourceAcc, errstack.WrapAsReq(errAcq, "Couldn't acquire trade's sourceAcc")
}

func parseExpireTime(expiresAt string) (*time.Time, errstack.E) {
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return nil, errstack.WrapAsReq(err, "Can't parse the expire time")
	}
	if t.Before(time.Now()) {
		return nil, errstack.NewReq("Expire time can't be in the past")
	}
	return &t, nil
}

func updateTrade(ctx context.Context, db driver.Database, t *model.Trade) (driver.DocumentMeta, errstack.E) {
	if t.CheckTradeClosed() {
		return driver.DocumentMeta{}, errstack.NewReq("You can't modify closed trade")
	}
	return dal.UpdateTrade(ctx, db, t)
}
