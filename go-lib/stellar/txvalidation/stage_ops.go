package txvalidation

import (
	"fmt"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/build"
)

// ValidateStageCloseReqTX validates tx for stage close transaction
func ValidateStageCloseReqTX(signedTx string, stageID uint, t *model.Trade, u *model.User, op model.Approval) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, errstack.E) {
	eb, se, vb := prevalidateTradeDataTx(t, u, signedTx)
	if !vb.IsEmpty() {
		return eb, se, vb.ToErrstackBuilder().ToReqErr()
	}
	validateMemo(vb, se.MemoHash, "")
	validateManageData(vb, se.DataValues, t.SCAddr, fmt.Sprint(stageID), model.TxTradeEntityStageCloseReqs, op)
	return eb, se, vb.ToErrstackBuilder().ToReqErr()
}

// ValidateStageAddReqTX validates tx for stage add transaction
func ValidateStageAddReqTX(signedTx string, stageAddReqID uint, t *model.Trade, u *model.User, op model.Approval) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, errstack.E) {
	eb, se, vb := prevalidateTradeDataTx(t, u, signedTx)
	if !vb.IsEmpty() {
		return eb, se, vb.ToErrstackBuilder().ToReqErr()
	}
	validateMemo(vb, se.MemoHash, "")
	validateManageData(vb, se.DataValues, t.SCAddr, fmt.Sprint(stageAddReqID), model.TxTradeEntityStageAdd, op)
	return eb, se, vb.ToErrstackBuilder().ToReqErr()
}
