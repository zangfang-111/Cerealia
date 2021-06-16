package stellar

import (
	"encoding/hex"
	"fmt"

	"bitbucket.org/cerealia/apps/go-lib/model"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

func makeTx(muts ...b.TransactionMutator) (string, error) {
	tx, err := b.Transaction(muts...)
	if err != nil {
		return "", err
	}
	txe := b.TransactionEnvelopeBuilder{
		E: &xdr.TransactionEnvelope{Tx: *tx.TX},
	}
	return txe.Base64()
}

func memoTextToHash(memo string) xdr.Hash {
	bMemo, _ := hex.DecodeString(memo)
	var hash xdr.Hash
	copy(hash[:], bMemo)
	return hash
}

func docPathFormat(id model.TradeStageDocPath) string {
	return fmt.Sprintf("%d:%d", id.StageIdx, id.StageDocIdx)
}
