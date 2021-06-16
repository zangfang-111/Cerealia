// Package trades contains REST api services for trades
package trades

import (
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/model/txlog"
	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/errstack"
)

// DocHandler is a route object for docs
type DocHandler struct {
	StellarDriver  *stellar.Driver
	TxSourceDriver txsource.Driver
}

// HandlePostTradeStageDoc is to upload trade document to server
func (h DocHandler) HandlePostTradeStageDoc(c *routing.Context) error {
	ctx, db, u, err := getAndCheckAuthUser(c)
	if err != nil {
		return err
	}
	now := time.Now()
	input, t, stage, err := readAndValidate(ctx, c, db, u, now)
	if err != nil {
		return errstack.WrapAsReq(err, "Request is not valid")
	}
	nextStageDocIdx := uint(len(stage.Docs))
	eBuilder, _, err := validateDocAddTX(input, t, u, nextStageDocIdx, now)
	if err != nil {
		return err
	}
	fi, err := storeDocFile(c.Request, 0, tradeDocDir)
	if err != nil {
		return err
	}
	// check the blake2s hash of the uploading document from frontend
	if input.Hash != fi.Hash {
		return errstack.NewReq("Wrong hash of uploaded document")
	}
	input.FileInfo = &fi
	ld := h.StellarDriver.WithTxLogger(
		txlog.New(ctx, db, model.StellarLedger, input.Tid, &input.StageIdx, &nextStageDocIdx, u.ID),
		h.TxSourceDriver.IsAcquiredFn(ctx, t.ID, u.ID))
	sourceAccs, erre := h.TxSourceDriver.Find(ctx, t.SCAddr, t.ID, u.ID)
	if erre != nil {
		return erre
	}
	txResult, err := ld.SignAndSendEnvelopeSource(eBuilder, sourceAccs)
	if err != nil {
		return err
	}
	modelStageDoc, err := input.updateDB(ctx, db, t, txResult.Hash)
	if err != nil {
		return err
	}
	if _, err = tradeStageDocAddNotif(ctx, db, t, u, input.StageIdx, input.WithApproval); err != nil {
		return err
	}
	return respondWithJSON(c, *modelStageDoc)
}

// HandleGetDocByID retrieves a document
func (h DocHandler) HandleGetDocByID(c *routing.Context) error {
	ctx, db, u, err := getAndCheckAuthUser(c)
	if err != nil {
		return err
	}
	docID := c.Param("docID")
	t, errs := dal.GetTradeOfDocument(ctx, db, docID)
	if errs != nil {
		return errs
	}
	if errs = t.CanBeModifiedBy(u); errs != nil {
		return errs
	}
	return serveDocFile(ctx, c, db, docID)
}

// SetTradeRoutes sets trade routes
func SetTradeRoutes(routerG *routing.RouteGroup, sd *stellar.Driver, txsd txsource.Driver) {
	h := DocHandler{sd, txsd}
	routerG.Post("/stage-docs", h.HandlePostTradeStageDoc)
	routerG.Get("/stage-docs/<docID>", h.HandleGetDocByID)
}
