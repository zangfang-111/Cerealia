package trades

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	"bitbucket.org/cerealia/apps/go-lib/utils"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
	bat "github.com/robert-zaremba/go-bat"
)

// UploadModel field names. Have to be the same as in UploadModel
const (
	tidField       = "tid"
	stageIdxField  = "stageIdx"
	expiresAtField = "expiresAt"
	signedTxField  = "signedTx"
)

// TradeStageDocInput is for document upload data
type (
	TradeStageDocInput struct {
		Tid          string          `form:"tid"`
		StageIdx     uint            `form:"stageIdx"`
		Uploader     string          // Loaded from session
		FileName     string          `form:"fileName"`
		Note         string          `form:"note"`
		ExpiresAt    string          `form:"expiresAt"`
		FileInfo     *model.FileInfo `form:"fileinfos"`
		Hash         string          `form:"hash"`
		SignedTx     string          `form:"signedTx"`
		WithApproval bool            `form:"withApproval"`
	}
	// TradeStageDocInputP contains parsed and validated fields
	TradeStageDocInputP struct {
		TradeStageDocInput
		ExpiresAtTime time.Time
	}
)

func (u TradeStageDocInput) parseFields() (*TradeStageDocInputP, errstack.Builder) {
	errb := u.Validate()
	t, _ := time.Parse(time.RFC3339, u.ExpiresAt)
	return &TradeStageDocInputP{
		TradeStageDocInput: u,
		ExpiresAtTime:      t,
	}, errb
}

// updateDB saves file information in DB
func (u TradeStageDocInputP) updateDB(ctx context.Context, db driver.Database, t *model.Trade, stellarTxHash string) (*model.TradeStageDoc, errstack.E) {
	b := u.ValidateStageIdx(len(t.Stages))
	if b.NotNil() {
		return nil, b.ToReqErr()
	}
	if t.CheckTradeClosed() {
		return nil, errstack.NewReq("You can't modify this trade because this trade has already been closed")
	}
	fi := u.FileInfo
	ext := filepath.Ext(fi.FileName)[1:] // remove '.'
	stage := &t.Stages[u.StageIdx]
	if stage.IsDeletedOrClosed() {
		return nil, errstack.NewReq("“You can’t change deleted or closed stage")
	}
	d := model.Doc{
		Hash:      fi.Hash, // TODO: remove this when using scalars
		Name:      u.FileName,
		Note:      u.Note,
		Type:      ext,
		URL:       fi.URL,
		CreatedBy: u.Uploader,
		CreatedAt: time.Now().UTC(),
	}
	de := model.TradeDocEdge{
		TradeID:     t.ID,
		StageIdx:    u.StageIdx,
		StageDocIdx: uint(len(stage.Docs)),
	}
	meta, errs := dal.InsertTradeDoc(ctx, db, &d, de)
	if errs != nil {
		return nil, errs
	}
	docStatus := model.ApprovalSubmitted
	if u.WithApproval {
		docStatus = model.ApprovalPending
	}
	sd := model.TradeStageDoc{
		DocID:     meta.Key,
		Status:    docStatus,
		ExpiresAt: u.ExpiresAtTime,
		ReqTx:     stellarTxHash,
	}
	stage.Docs = append(stage.Docs, sd)
	_, errs = dal.UpdateTrade(ctx, db, t)
	if errs != nil {
		return nil, errs
	}
	return &sd, errs
}

func tradeStageDocAddNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User,
	stageIdx uint, withApproval bool) (*model.Notification, errstack.E) {
	receiver, _ := t.NotificationReceiver(u)
	idx := strconv.Itoa(len(t.Stages[stageIdx].Docs) - 1)
	action := model.ApprovalSubmitted
	if withApproval {
		action = model.ApprovalPending
	}
	n := model.Notification{
		CreatedAt:   time.Now().UTC(),
		Receiver:    receiver,
		TriggeredBy: u.ID,
		Type:        model.NotifTypeAction,
		Dismissed:   []string{},
		EntityID:    bat.StrJoin("/", string(dbconst.ColTrades)+":"+t.ID, "stages:"+utils.UintToString(stageIdx), "docs:"+idx),
		Msg:         fmt.Sprintf("New stage doc has been created by %s %s", u.FirstName, u.LastName),
		Action:      action,
	}
	return &n, dal.InsertNotification(ctx, db, &n)
}
