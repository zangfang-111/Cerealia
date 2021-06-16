// Package txlog contains transaction record and log functions
package txlog

import (
	"context"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/model/txlog/txlogi"
	driver "github.com/arangodb/go-driver"
	"github.com/stellar/go/xdr"
)

// New creates a new instance of TxLogger
func New(ctx context.Context, db driver.Database, ledger model.LedgerEnum, tradeID string, stageID, docID *uint, userID string) txlogi.Logger {
	return txLogger{
		ctx:     ctx,
		db:      db,
		userID:  userID,
		tradeID: tradeID,
		stageID: stageID,
		docID:   docID,
		ledger:  ledger,
	}
}

type txLogger struct {
	ctx     context.Context
	db      driver.Database
	tradeID string
	stageID *uint            // optional
	docID   *uint            // optional
	userID  string           // user that initiates the tx
	ledger  model.LedgerEnum // Stellar or other
}

func (l txLogger) makeTxEntry(rawTx string, status model.TxStatusEnum, e *xdr.TransactionEnvelope) model.TxLog {
	return model.TxLog{
		TxStatus:  status,
		Ledger:    l.ledger,
		RawTx:     rawTx,
		CreatedBy: l.userID,
		UpdatedAt: time.Now().UTC(),
		Nonce:     e.Tx.SeqNum,
		SourceAcc: e.Tx.SourceAccount.Address(),
	}
}

func (l txLogger) makeTxEntryEdge() model.TxLogEdge {
	return model.TxLogEdge{
		TradeID:     l.tradeID,
		StageIdx:    l.stageID,
		StageDocIdx: l.docID,
	}
}

// LogTxStatus implements interface txlogi.Logger
func (l txLogger) LogTxStatus(rawTx string, e *xdr.TransactionEnvelope, status model.TxStatusEnum) error {
	entry := l.makeTxEntry(rawTx, status, e)
	edge := l.makeTxEntryEdge()
	_, err := dal.UpsertTxLogEntry(
		l.ctx,
		l.db,
		&entry,
		&edge,
	)
	return err
}
