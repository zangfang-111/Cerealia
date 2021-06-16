package txlog

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	"github.com/stellar/go/xdr"
)

// NoopTxLogger does no logging
type NoopTxLogger struct {
}

// LogTxStatus immplements interface for tx logger
func (n NoopTxLogger) LogTxStatus(string, *xdr.TransactionEnvelope, model.TxStatusEnum) error {
	return nil
}
