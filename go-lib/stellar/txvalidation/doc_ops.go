package txvalidation

import (
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/validation"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/build"
)

// validateDataDocOp does the actual validation of the data values
func validateDataDocOp(vb *validation.Builder, obtained dataMap, tradeAcc model.SCAddr, id model.TradeStageDocPath, op model.Approval) bool {
	return validateManageData(vb, obtained, tradeAcc,
		idxJoin(id.StageIdx, id.StageDocIdx),
		model.TxTradeEntityStageDoc,
		op)
}

func prevalidateDocAddTx(signedTx string, id model.TradeStageDocPath, t *model.Trade, u *model.User,
	docHash string) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, *validation.Builder) {
	eb, se, vb := prevalidateTradeDataTx(t, u, signedTx)
	if !vb.IsEmpty() {
		return eb, se, vb
	}
	validateMemo(vb, se.MemoHash, docHash)
	return eb, se, vb
}

// ValidateDocAddExpireTX validates tx for new doc and checks tx expiration time
func ValidateDocAddExpireTX(signedTx string, id model.TradeStageDocPath, t *model.Trade, u *model.User,
	docHash string, expiresAt, now time.Time) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, errstack.E) {
	eb, se, vb := prevalidateDocAddTx(signedTx, id, t, u, docHash)
	validateManageDataWithExpireTime(vb, se.DataValues, t.SCAddr, idxJoin(id.StageIdx, id.StageDocIdx),
		model.TxTradeEntityStageDoc, model.ApprovalPending, expiresAt)
	if now.After(expiresAt) { // expiresAt equality to tx data is checked in validateManageDataWithExpireTime
		vb.Append(validationFieldTX, badData)
	}
	return eb, se, vb.ToErrstackBuilder().ToReqErr()
}

// ValidateDocAddTX validates tx for new doc txs that don't have expiration (without confirmation)
func ValidateDocAddTX(signedTx string, id model.TradeStageDocPath, t *model.Trade, u *model.User,
	docHash string) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, errstack.E) {
	eb, se, vb := prevalidateDocAddTx(signedTx, id, t, u, docHash)
	validateDataDocOp(vb, se.DataValues, t.SCAddr, id, model.ApprovalSubmitted)
	return eb, se, vb.ToErrstackBuilder().ToReqErr()
}

func validateDocOpTX(signedTx string, id model.TradeStageDocPath, t *model.Trade, u *model.User, op model.Approval) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, errstack.E) {
	eb, se, vb := prevalidateTradeDataTx(t, u, signedTx)
	if !vb.IsEmpty() {
		return eb, se, vb.ToErrstackBuilder().ToReqErr()
	}
	validateMemo(vb, se.MemoHash, id.StageDocHash)
	validateDataDocOp(vb, se.DataValues, t.SCAddr, id, op)
	return eb, se, vb.ToErrstackBuilder().ToReqErr()
}

// ValidateDocApproveTX validates tx for doc approval
func ValidateDocApproveTX(signedTx string, id model.TradeStageDocPath, t *model.Trade, u *model.User) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, errstack.E) {
	return validateDocOpTX(signedTx, id, t, u, model.ApprovalApproved)
}

// ValidateDocRejectTX validates tx for doc rejection
func ValidateDocRejectTX(signedTx string, id model.TradeStageDocPath, t *model.Trade, u *model.User) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, errstack.E) {
	return validateDocOpTX(signedTx, id, t, u, model.ApprovalRejected)
}
