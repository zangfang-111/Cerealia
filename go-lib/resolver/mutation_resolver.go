// Package resolver contains all resolver mutations and queries and also its related functions
package resolver

import (
	"context"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/auth"
	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/model/txlog"
	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txvalidation"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/keypair"
)

type mutationResolver struct{ *resolver }

func (r mutationResolver) UserSignup(ctx context.Context, input *model.NewUserInput) (*int, error) {
	errs := input.CleanAndValidate()
	if errs != nil {
		return nil, errs
	}
	_, errs = dal.InsertUser(ctx, r.db, input)
	return nil, errs
}

func (r mutationResolver) UserLogin(ctx context.Context, input model.UserLoginInput) (*model.AuthUser, error) {
	user, errs := dal.UserLogin(ctx, r.db, input)
	if errs != nil {
		return nil, errs
	}
	// Generate JWT token
	token, errs := auth.CreateJWT(user.ID)
	if errs != nil {
		return nil, errs
	}
	return &model.AuthUser{
		ID:    user.ID,
		Token: token,
	}, nil
}

func (r mutationResolver) UserPasswordChange(ctx context.Context, input model.ChangePasswordInput) (*int, error) {
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return nil, errs
	}
	return nil, dal.ChangePassword(ctx, r.db, u, input)
}

func (r mutationResolver) UserEmailChange(ctx context.Context, input []string) (*int, error) {
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return nil, errs
	}
	return nil, dal.ChangeEmail(ctx, r.db, u, input)
}

func (r mutationResolver) UserProfileUpdate(ctx context.Context, input model.UserProfileInput) (*model.User, error) {
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return nil, errs
	}
	return dal.UpdateUserProfile(ctx, r.db, u, input)
}

func (r mutationResolver) OrganizationCreate(ctx context.Context, input model.OrgInput) (*model.Organization, error) {
	newOrg := model.Organization{
		Name:      input.Name,
		Address:   input.Address,
		Email:     input.Email,
		Telephone: input.Telephone,
	}
	if errs := newOrg.ValidateNewOrg(); errs != nil {
		return nil, errs
	}
	return dal.CreateOrganization(ctx, r.db, &newOrg)
}

// TradeCreate creates a new trade.
func (r mutationResolver) TradeCreate(ctx context.Context, input model.NewTradeInput) (*model.Trade, error) {
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return nil, errs
	}
	errb := input.Validate(input.TemplateID, u)
	if errb.NotNil() {
		return nil, errb.ToReqErr()
	}
	pks, errs := getTradeParticipants(ctx, r.db, input.BuyerID, input.SellerID)
	if errs != nil {
		return nil, errstack.WrapAsDomain(errs, "Private key not found")
	}
	keypair, erre := keypair.Random()
	if erre != nil {
		return nil, erre
	}
	t := model.Trade{
		Name:         input.Name,
		TemplateID:   input.TemplateID,
		Buyer:        pks.Buyer,
		Seller:       pks.Seller,
		Description:  input.Description,
		CreatedBy:    u.ID,
		SCAddr:       model.SCAddr(keypair.Address()),
		TradeOfferID: input.TradeOfferID,
		StageAddReqs: []model.TradeStageAddReq{},
		CloseReqs:    []model.ApproveReq{},
		Moderating:   model.DoneStatusNil,
	}
	// get template Data and copy them into trade
	tt, errs := dal.GetTradeTempate(ctx, r.db, t.TemplateID)
	if errs != nil {
		return &t, errs
	}
	t.Stages = tt.BuildStages()
	meta, errs := dal.InsertTrade(ctx, r.db, &t)
	if errs != nil {
		return &t, errs
	}
	t.ID = meta.Key
	sourceAccs, err := r.txSourceDriver.Create(ctx, *keypair, t.ID, u.ID, model.TxSourceAccTypeTrade)
	if err != nil {
		return nil, errstack.WrapAsInf(err)
	}
	defer errstack.CallAndLog(logger, r.txSourceDriver.ReleaseFn(ctx, t.ID, u.ID))
	ld := r.mkStellarLogDriver(ctx, u.ID, &t, nil, nil)
	errs = stellar.CreateTradeAccount(ld, pks, sourceAccs)
	if errs != nil {
		return nil, errs
	}
	if _, errs = tadeCreateNotif(ctx, r.db, &t, u); errs != nil {
		return &t, errs
	}
	return &t, errs
}

// TradeStageAddReq creates an add request of new trade stage.
func (r mutationResolver) TradeStageAddReq(ctx context.Context, input model.NewStageInput, signedTx string, withApproval bool) (*model.TradeStageAddReq, error) {
	// TODO: move this validation to NewStageInput validation
	if len(input.Reason) < 10 {
		return nil, errstack.NewReq("Reason must be at least 10 characters long")
	}
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return nil, errs
	}
	t, errs := dal.GetTrade(ctx, r.db, input.Tid)
	if errs != nil {
		return nil, errs
	}
	var errb = errstack.NewBuilder()
	owner := model.ParseTradeActor(input.Owner, errb)
	reqActor, errs := t.Requester(u)
	errb.Put("requester", errs)
	errb.Put("permissions", t.CanBeModifiedBy(u))
	if errb.NotNil() {
		return nil, errb.ToReqErr()
	}
	defer errstack.CallAndLog(logger, r.txSourceDriver.ReleaseFn(ctx, t.ID, u.ID))
	var sr = model.TradeStageAddReq{
		Name:        input.Name,
		Description: input.Description,
		Owner:       owner,
		ApproveReq: model.ApproveReq{
			Status:    model.ApprovalPending,
			ReqActor:  reqActor,
			ReqBy:     u.ID,
			ReqAt:     time.Now().UTC(),
			ReqReason: input.Reason,
		},
	}
	operation := dealCreateStageApproval(t, &sr, u, withApproval)
	newStageIdx := uint(len(t.StageAddReqs))
	eBuilder, _, err := txvalidation.ValidateStageAddReqTX(signedTx, newStageIdx, t, u, operation)
	if err != nil {
		return nil, err
	}
	ld := r.mkStellarLogDriver(ctx, u.ID, t, &newStageIdx, nil)
	sourceAccs, erre := r.txSourceDriver.Find(ctx, t.SCAddr, t.ID, u.ID)
	if erre != nil {
		return nil, erre
	}
	txResult, err := ld.SignAndSendEnvelopeSource(eBuilder, sourceAccs)
	if err != nil {
		return nil, err
	}
	sr.ApproveReq.ReqTx = txResult.Hash

	t.StageAddReqs = append(t.StageAddReqs, sr)
	if _, errs = tradeStageAddReqNotif(ctx, r.db, t, u, withApproval); errs != nil {
		return nil, errs
	}
	_, errs = updateTrade(ctx, r.db, t)
	return &sr, errs
}

// TradeStageAddReqApprove approves the add request of new stage.
func (r mutationResolver) TradeStageAddReqApprove(ctx context.Context, id model.TradeStagePath, signedTx string) (*model.TradeStage, error) {
	return tradeStageAddApproval(ctx, r, id, signedTx, "", true)
}

// TradeStageAddReqReject rejects an add request of new stage
func (r mutationResolver) TradeStageAddReqReject(ctx context.Context, id model.TradeStagePath, signedTx string, reason string) (*int, error) {
	_, errs := tradeStageAddApproval(ctx, r, id, signedTx, reason, false)
	return nil, errs
}

// TradeStageDelReq creates a delete request of existing stage.
func (r mutationResolver) TradeStageDelReq(ctx context.Context, id model.TradeStagePath, reason string) (*model.ApproveReq, error) {
	var errb = errstack.NewBuilder()
	if len(reason) < 10 {
		errb.Put("lengthError", errstack.NewReq("Reason must be at least 10 characters long"))
	}
	reqActor, u, t, errs := getTradeRequester(ctx, r.db, id.Tid)
	errb.Put("getRequester", errs)
	if errb.NotNil() {
		return nil, errb.ToReqErr()
	}
	s, errs := t.GetStage(id.StageIdx)
	errb.Put("getStage", errs)
	// check the possibility of making delete request.
	if !s.CanMakeDelReq() {
		errb.Put("checkCanDelete", errstack.NewReq("You can't create a delete request to this stage."))
	}
	if s.Owner != reqActor {
		errb.Put("checkActor", errstack.NewReq("You can't create the request because you are not the owner of this stage."))
	}
	if s.IsDeletedOrClosed() {
		errb.Put("checkStageStatus", errChangeStage)
	}
	if errb.NotNil() {
		return nil, errb.ToReqErr()
	}
	// add delete request
	var ar = model.ApproveReq{
		Status:    model.ApprovalPending,
		ReqActor:  reqActor,
		ReqBy:     u.ID,
		ReqAt:     time.Now().UTC(),
		ReqReason: reason}
	s.DelReqs = append(s.DelReqs, ar)
	if _, errs = tradeStageDelReqNotif(ctx, r.db, t, u, id); errs != nil {
		return nil, errs
	}
	_, errs = updateTrade(ctx, r.db, t)
	return &ar, errs
}

// TradeStageDelReqApprove approves the delete request of existing stage.
func (r mutationResolver) TradeStageDelReqApprove(ctx context.Context, id model.TradeStagePath) (*model.ApproveReq, error) {
	return tradeStageDeleteApproval(ctx, r, id, "", true)
}

// TradeStageDelReqReject rejects the delete request of existing stage.
func (r mutationResolver) TradeStageDelReqReject(ctx context.Context, id model.TradeStagePath, reason string) (*int, error) {
	_, errs := tradeStageDeleteApproval(ctx, r, id, reason, false)
	return nil, errs
}

// TradeStageDocApprove approves the document of a stage.
func (r mutationResolver) TradeStageDocApprove(ctx context.Context, id model.TradeStageDocPath,
	signedTx string) (*model.TradeStageDoc, error) {
	return tradeStageDocApproval(ctx, r, id, signedTx, "", true)
}

// TradeStageDocReject rejects the document of a stage.
func (r mutationResolver) TradeStageDocReject(ctx context.Context, id model.TradeStageDocPath,
	signedTx, reason string) (*int, error) {
	_, errs := tradeStageDocApproval(ctx, r, id, signedTx, reason, false)
	return nil, errs
}

// TradeStageCloseReq creates a trade close request
func (r mutationResolver) TradeStageCloseReq(ctx context.Context, id model.TradeStagePath,
	signedTx, reason string) (*model.ApproveReq, error) {
	var errb = errstack.NewBuilder()
	if len(reason) < 10 {
		errb.Put("lengthError", errstack.NewReq("Reason must be at least 10 characters long"))
	}
	reqActor, u, t, errs := getTradeRequester(ctx, r.db, id.Tid)
	errb.Put("getReqActor", errs)
	if errb.NotNil() {
		return nil, errb.ToReqErr()
	}
	defer errstack.CallAndLog(logger, r.txSourceDriver.ReleaseFn(ctx, t.ID, u.ID))
	s, errs := t.GetStage(id.StageIdx)
	if errs != nil {
		return nil, errs
	}
	if !s.CanMakeCloseReq() {
		errb.Put("can't Close", errstack.NewReq("You can't create a close request to this stage."))
	}
	if s.IsDeletedOrClosed() {
		errb.Put("checkStageStatus", errChangeStage)
	}
	errb.Put("permissions", s.AssertOwnedBy(t, u.ID))
	if errb.NotNil() {
		return nil, errb.ToReqErr()
	}
	var ar = model.ApproveReq{
		Status:    model.ApprovalPending,
		ReqActor:  reqActor,
		ReqBy:     u.ID,
		ReqAt:     time.Now().UTC(),
		ReqReason: reason}
	eBuilder, _, err := txvalidation.ValidateStageCloseReqTX(signedTx, id.StageIdx, t, u, model.ApprovalPending)
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
	ar.ReqTx = txResult.Hash
	s.CloseReqs = append(s.CloseReqs, ar)
	if _, errs = tradeStageCloseReqNotif(ctx, r.db, t, u, id); errs != nil {
		return nil, errs
	}
	_, errs = updateTrade(ctx, r.db, t)
	return &ar, errs
}

func (r mutationResolver) TradeStageCloseReqApprove(ctx context.Context, id model.TradeStagePath,
	signedTx string) (*model.ApproveReq, error) {
	return tradeStageCloseApproval(ctx, r, id, signedTx, "", true)
}

func (r mutationResolver) TradeStageCloseReqReject(ctx context.Context, id model.TradeStagePath,
	signedTx, reason string) (*int, error) {
	_, errs := tradeStageCloseApproval(ctx, r, id, signedTx, reason, false)
	return nil, errs
}

func (r mutationResolver) TradeCloseReq(ctx context.Context, id string, reason string, signedTx string) (*model.ApproveReq, error) {
	var errb = errstack.NewBuilder()
	if len(reason) < 10 {
		errb.Put("lengthError", errstack.NewReq("Reason must be at least 10 characters long"))
	}
	reqActor, u, t, errs := getTradeRequester(ctx, r.db, id)
	errb.Put("getReqActor", errs)
	if errb.NotNil() {
		return nil, errb.ToReqErr()
	}
	defer errstack.CallAndLog(logger, r.txSourceDriver.ReleaseFn(ctx, t.ID, u.ID))
	// check the possibility of making close request.
	if !t.CanMakeCloseReq() {
		return nil, errstack.NewReq("You can't create a close request to this trade.")
	}
	var ar = model.ApproveReq{
		Status:    model.ApprovalPending,
		ReqActor:  reqActor,
		ReqBy:     u.ID,
		ReqAt:     time.Now().UTC(),
		ReqReason: reason}
	eBuilder, _, err := txvalidation.ValidateTradeCloseReqTX(signedTx, id, t, u, model.ApprovalPending)
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
	ar.ReqTx = txResult.Hash
	t.CloseReqs = append(t.CloseReqs, ar)
	if _, errs = tradeCloseReqNotif(ctx, r.db, t, u); errs != nil {
		return nil, errs
	}
	_, errs = updateTrade(ctx, r.db, t)
	return &ar, errs
}

func (r mutationResolver) TradeCloseReqApprove(ctx context.Context, id string, signedTx string) (*model.ApproveReq, error) {
	return tradeCloseApproval(ctx, r, id, signedTx, "", true)
}

func (r mutationResolver) TradeCloseReqReject(ctx context.Context, id string, reason string, signedTx string) (*int, error) {
	_, errs := tradeCloseApproval(ctx, r, id, signedTx, reason, false)
	return nil, errs
}

func (r mutationResolver) TradeStageSetExpireTime(ctx context.Context, id model.TradeStagePath, expiresAt string) (*int, error) {
	reqActor, u, t, errs := getTradeRequester(ctx, r.db, id.Tid)
	if errs != nil {
		return nil, errs
	}
	s, errs := t.GetStage(id.StageIdx)
	if errs != nil {
		return nil, errs
	}
	var errb = errstack.NewBuilder()
	if s.Owner == reqActor {
		errb.Put("getReqActor", errstack.NewReq("You can't set expire time for this stage since you are the owner."))
	}
	if s.IsDeletedOrClosed() {
		errb.Put("checkStageStatus", errChangeStage)
	}
	expTime, errs := parseExpireTime(expiresAt)
	errb.Put("timeParseErr", errs)
	if errb.NotNil() {
		return nil, errb.ToReqErr()
	}
	s.ExpiresAt = expTime
	if _, errs = tradeStageDelReqNotif(ctx, r.db, t, u, id); errs != nil {
		return nil, errs
	}
	_, errs = updateTrade(ctx, r.db, t)
	return nil, errs
}

func (r mutationResolver) TradeOfferCreate(ctx context.Context, input model.TradeOfferInput) (*model.TradeOffer, error) {
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return nil, errs
	}
	errs = input.ValidateInput()
	if errs != nil {
		return nil, errs
	}
	tof := model.TradeOffer{
		Price:       input.Price,
		PriceType:   input.PriceType,
		IsSell:      input.IsSell,
		Currency:    input.Currency,
		CreatedBy:   u.ID,
		CreatedAt:   time.Now().UTC(),
		ExpiresAt:   input.ExpiresAt,
		Origin:      input.Origin,
		IsAnonymous: input.IsAnonymous,
		Commodity:   input.Commodity,
		ComType:     input.ComType,
		Quality:     input.Quality,
		OrgID:       input.OrgID,
		Incoterm:    input.Incoterm,
		MarketLoc:   input.MarketLoc,
		Vol:         input.Vol,
		Shipment:    input.Shipment,
		Note:        input.Note,
		DocID:       input.DocID,
	}
	if tof.DocID != nil {
		if errs = dal.AssertDocExists(ctx, r.db, *tof.DocID); errs != nil {
			return nil, errs
		}
	}
	return dal.InsertTradeOffer(ctx, r.db, &tof)
}

func (r mutationResolver) TradeOfferClose(ctx context.Context, id string) (*int, error) {
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return nil, errs
	}
	tof, errs := dal.GetTradeOffer(ctx, r.db, id)
	if errs != nil {
		return nil, errs
	}
	if tof.CreatedBy != u.ID {
		return nil, model.ErrUnauthorized
	}
	return nil, dal.CloseTradeOffer(ctx, r.db, id)
}

func (r mutationResolver) NotificationDismiss(ctx context.Context, id string) (*int, error) {
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	return nil, dal.NotificationDismiss(ctx, r.db, u.ID, id)
}

func (r mutationResolver) MkTradeStageDocTx(ctx context.Context, id model.TradeStageDocPath, operationType model.Approval, expiresAt *time.Time) (string, error) {
	_, sources, errs := validateOpTypeAndGetTrade(ctx, r.db, r.txSourceDriver, id.Tid, operationType)
	if errs != nil {
		return "", errs
	}
	if operationType == model.ApprovalPending {
		return stellar.MkTradeDocApprovalExpireTx(r.stellarDriver, sources, id, model.TxTradeEntityStageDoc, operationType, *expiresAt)
	}
	return stellar.MkTradeDocApprovalTx(r.stellarDriver, sources, id, model.TxTradeEntityStageDoc, operationType)
}

func (r mutationResolver) MkTradeStageCloseTx(ctx context.Context, id model.TradeStagePath, operationType model.Approval) (string, error) {
	_, sources, errs := validateOpTypeAndGetTrade(ctx, r.db, r.txSourceDriver, id.Tid, operationType)
	if errs != nil {
		return "", errs
	}
	return stellar.MkTradeStageOperationTx(r.stellarDriver, sources, id.StageIdx, model.TxTradeEntityStageCloseReqs, operationType)
}

func (r mutationResolver) MkTradeStageAddTx(ctx context.Context, id model.TradeStagePath, operationType model.Approval) (string, error) {
	_, sources, errs := validateOpTypeAndGetTrade(ctx, r.db, r.txSourceDriver, id.Tid, operationType)
	if errs != nil {
		return "", errs
	}
	return stellar.MkTradeStageOperationTx(r.stellarDriver, sources, id.StageIdx, model.TxTradeEntityStageAdd, operationType)
}

func (r mutationResolver) mkStellarLogDriver(ctx context.Context, userID string, t *model.Trade, stageID, docID *uint) *stellar.WrappedDriver {
	l := txlog.New(ctx, r.db, model.StellarLedger, t.ID, stageID, docID, userID)
	return r.stellarDriver.WithTxLogger(l, r.txSourceDriver.IsAcquiredFn(ctx, t.ID, userID))
}

func (r mutationResolver) MkTradeCloseTx(ctx context.Context, id string, operationType model.Approval) (string, error) {
	_, sources, errs := validateOpTypeAndGetTrade(ctx, r.db, r.txSourceDriver, id, operationType)
	if errs != nil {
		return "", errs
	}
	return stellar.MkTradeCloseTx(r.stellarDriver, sources, id, model.TxTradeEntityTradeCloseReqs, operationType)
}
