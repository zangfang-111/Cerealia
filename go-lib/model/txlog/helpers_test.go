package txlog

import (
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

func decodeAccountID(pubKey string) (*xdr.AccountId, error) {
	decoded, err := strkey.Decode(strkey.VersionByteAccountID, pubKey)
	if err != nil {
		return nil, errstack.WrapAsInfF(err, "Malformed transaction source account")
	}
	var raw xdr.Uint256
	copy(raw[:], decoded[:32])
	accID, err := xdr.NewAccountId(xdr.PublicKeyTypePublicKeyTypeEd25519, raw)
	return &accID, err
}

func makeTxEnvelope(nonce int, sourceAcc string) (*xdr.TransactionEnvelope, error) {
	accID, err := decodeAccountID(sourceAcc)
	if err != nil {
		return nil, errstack.WrapAsInfF(err, "Can't decode accountID")
	}
	return &xdr.TransactionEnvelope{
		Tx: xdr.Transaction{
			SeqNum:        xdr.SequenceNumber(nonce),
			SourceAccount: *accID,
		},
	}, nil
}
