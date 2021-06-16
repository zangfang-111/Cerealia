// Package txlogi is created because implementation of TxLogger causes a dependency import cycle
package txlogi

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	"github.com/stellar/go/xdr"
)

// Logger is an interface for transaction logging
type Logger interface {
	// LogTxStatus logs the transaction
	LogTxStatus(transaction string, e *xdr.TransactionEnvelope, status model.TxStatusEnum) error
}
