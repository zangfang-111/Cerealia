package txlogi

import "bitbucket.org/cerealia/apps/go-lib/model"

// NoopTxLogger is a no operation tx logger
type NoopTxLogger struct {
}

// LogTxStatus implements an interface txlogi.Logger
func (n NoopTxLogger) LogTxStatus(transaction string, status model.TxStatusEnum) error {
	return nil
}
