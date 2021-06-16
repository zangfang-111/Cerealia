package resolver

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txvalidation"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/build"
)

func tradeStageAddApproval(ctx context.Context, r mutationResolver, id model.TradeStagePath, signedTx, reason string, isApprove bool) (*model.TradeStage, error) {
	t, sr, u, errs := prepareTradeStageAddReqApproval(ctx, r.db, id, isApprove)
	if errs != nil {
		return nil, errs
	}
	defer errstack.CallAndLog(logger, r.txSourceDriver.ReleaseFn(ctx, t.ID, u.ID))
	op := model.ApprovalRejected
	if isApprove {
		op = model.ApprovalApproved
	}
	eBuilder, _, err := txvalidation.ValidateStageAddReqTX(signedTx, id.StageIdx, t, u, op)
	if err != nil {
		return nil, err
	}

	ld := r.mkStellarLogDriver(ctx, u.ID, t, &id.StageIdx, nil)
	sourceAccs, erre := r.txSourceDriver.Find(ctx, t.SCAddr, t.ID, u.ID)
	if erre != nil {
		return nil, erre
	}
	txResult, err := ld.SignAndSendEnvelopeSource(eBuilder, sourceAccs)
	if err != nil {
		return nil, err
	}
	sr.ApproveReq.ApprovedTx = txResult.Hash
	sr.Status = op
	sr.RejectReason = reason
	var stage *model.TradeStage

	if isApprove {
		if l := len(t.Stages); l > 0 {
			stage = &t.Stages[l-1]
		}
	}
	if _, errs = tradeStageAddApprovalNotif(ctx, r.db, t, u, id, isApprove); errs != nil {
		return nil, errs
	}
	_, errs = updateTrade(ctx, r.db, t)
	return stage, errs
}

func tradeStageDeleteApproval(ctx context.Context, r mutationResolver, id model.TradeStagePath, reason string, isApprove bool) (*model.ApproveReq, error) {
	t, delReq, u, errs := prepareTradeStageReqApproval(ctx, r.db, id, true)
	if errs != nil {
		return nil, errs
	}
	if isApprove {
		delReq.Status = model.ApprovalApproved
	} else {
		delReq.Status = model.ApprovalRejected
	}
	delReq.RejectReason = reason
	if _, errs = tradeStageDeleteApprovalNotif(ctx, r.db, t, u, id, isApprove); errs != nil {
		return nil, errs
	}
	_, errs = updateTrade(ctx, r.db, t)
	return delReq, errs
}

func tradeStageDocApproval(ctx context.Context, r mutationResolver, id model.TradeStageDocPath, signedTx string, reason string, isApprove bool) (*model.TradeStageDoc, error) {
	t, d, u, errs := prepareTradeStageDocApproval(ctx, r.db, id)
	if errs != nil {
		return nil, errs
	}
	defer errstack.CallAndLog(logger, r.txSourceDriver.ReleaseFn(ctx, t.ID, u.ID))
	var eBuilder *build.TransactionEnvelopeBuilder
	var err error
	if isApprove {
		eBuilder, _, err = txvalidation.ValidateDocApproveTX(signedTx, id, t, u)
	} else {
		eBuilder, _, err = txvalidation.ValidateDocRejectTX(signedTx, id, t, u)
	}
	if err != nil {
		return nil, err
	}
	ld := r.mkStellarLogDriver(ctx, u.ID, t, &id.StageIdx, &id.StageDocIdx)
	sourceAccs, erre := r.txSourceDriver.Find(ctx, t.SCAddr, t.ID, u.ID)
	if erre != nil {
		return nil, erre
	}
	txResult, err := ld.SignAndSendEnvelopeSource(eBuilder, sourceAccs)
	if err != nil {
		return nil, err
	}
	d.ApprovedTx = txResult.Hash
	if isApprove {
		d.Status = model.ApprovalApproved
	} else {
		d.Status = model.ApprovalRejected
	}
	d.RejectReason = reason
	if _, errs = tradeStageDocApprovalNotif(ctx, r.db, t, u, id, isApprove); errs != nil {
		return nil, errs
	}
	_, errs = updateTrade(ctx, r.db, t)
	return d, errs
}

func tradeStageCloseApproval(ctx context.Context, r mutationResolver, id model.TradeStagePath, signedTx, reason string, isApprove bool) (*model.ApproveReq, error) {
	t, closeReq, u, errs := prepareTradeStageReqApproval(ctx, r.db, id, false)
	if errs != nil {
		return nil, errs
	}
	defer errstack.CallAndLog(logger, r.txSourceDriver.ReleaseFn(ctx, t.ID, u.ID))
	op := model.ApprovalRejected
	if isApprove {
		op = model.ApprovalApproved
	}
	eBuilder, _, err := txvalidation.ValidateStageCloseReqTX(signedTx, id.StageIdx, t, u, op)
	if err != nil {
		return nil, err
	}
	ld := r.mkStellarLogDriver(ctx, u.ID, t, &id.StageIdx, nil)
	sourceAccs, erre := r.txSourceDriver.Find(ctx, t.SCAddr, t.ID, u.ID)
	if erre != nil {
		return nil, erre
	}
	txResult, err := ld.SignAndSendEnvelopeSource(eBuilder, sourceAccs)
	if err != nil {
		return nil, err
	}
	closeReq.ApprovedTx = txResult.Hash
	closeReq.Status = op
	closeReq.RejectReason = reason
	if _, errs = tradeStageCloseApprovalNotif(ctx, r.db, t, u, id, isApprove); errs != nil {
		return nil, errs
	}
	// close tradeOffer if it exists in the trade
	if id.StageIdx == 0 && t.TradeOfferID != nil && isApprove {
		if errs = dal.CloseTradeOffer(ctx, r.db, *t.TradeOfferID); errs != nil {
			return nil, errs
		}
	}
	_, errs = updateTrade(ctx, r.db, t)
	return closeReq, errs
}

func tradeCloseApproval(ctx context.Context, r mutationResolver, tid, signedTx, reason string, isApprove bool) (*model.ApproveReq, error) {
	t, closeReq, u, errs := prepareTradeCloseReqApproval(ctx, r.db, tid)
	if errs != nil {
		return nil, errs
	}
	defer errstack.CallAndLog(logger, r.txSourceDriver.ReleaseFn(ctx, t.ID, u.ID))
	op := model.ApprovalRejected
	if isApprove {
		op = model.ApprovalApproved
	}
	eBuilder, _, err := txvalidation.ValidateTradeCloseReqTX(signedTx, tid, t, u, op)
	if err != nil {
		return nil, err
	}
	ld := r.mkStellarLogDriver(ctx, u.ID, t, nil, nil)
	sourceAccs, erre := r.txSourceDriver.Find(ctx, t.SCAddr, t.ID, u.ID)
	if erre != nil {
		return nil, erre
	}
	txResult, err := ld.SignAndSendEnvelopeSource(eBuilder, sourceAccs)
	if err != nil {
		return nil, err
	}
	closeReq.ApprovedTx = txResult.Hash
	closeReq.Status = op
	closeReq.RejectReason = reason
	if _, errs = tradeCloseApprovalNotif(ctx, r.db, t, u, isApprove); errs != nil {
		return nil, errs
	}
	_, errs = updateTrade(ctx, r.db, t)
	return closeReq, errs
}
