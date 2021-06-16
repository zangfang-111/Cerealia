package model

import "bitbucket.org/cerealia/apps/go-lib/model/dbconst"

// LedgerEnum enum
type LedgerEnum string

// StellarLedger from stellar.org
const StellarLedger LedgerEnum = "stellar"

// TxStatusEnum defines logging statusses
type TxStatusEnum string

const (
	// TxStatusPending means waiting for tx to fail or complete
	TxStatusPending TxStatusEnum = "pending"
	// TxStatusOk means that tx had no execution errors
	TxStatusOk TxStatusEnum = "ok"
	// TxStatusFailed means that tx had an execution error
	TxStatusFailed TxStatusEnum = "failed"
)

// SetID implements dal.HasID interface
func (d *TxLog) SetID(id string) {
	d.ID = id
}

// ToEdgeDO accepts a full entity ID that involves the collection name
// it returns a DO (database object) that should be inserted directly into DB
func (e TxLogEdge) ToEdgeDO(fullTxLogEntryID string) interface{} {
	return TxLogEdgeDTO{
		FullTradeID: dbconst.ColTrades.FullID(e.TradeID),
		FullTxLogID: fullTxLogEntryID,
		StageIdx:    e.StageIdx,
		StageDocIdx: e.StageDocIdx,
	}
}
