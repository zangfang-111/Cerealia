package txvalidation

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/build"
)

// ValidateTradeCloseReqTX validates tx for trade close transaction
func ValidateTradeCloseReqTX(signedTx string, id string, t *model.Trade, u *model.User, op model.Approval) (*build.TransactionEnvelopeBuilder, *SimplifiedEnvelope, errstack.E) {
	eb, se, vb := prevalidateTradeDataTx(t, u, signedTx)
	if !vb.IsEmpty() {
		return eb, se, vb.ToErrstackBuilder().ToReqErr()
	}
	validateMemo(vb, se.MemoHash, "")
	validateManageData(vb, se.DataValues, t.SCAddr, id, model.TxTradeEntityTradeCloseReqs, op)
	return eb, se, vb.ToErrstackBuilder().ToReqErr()
}
