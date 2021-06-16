package txvalidation

import (
	"encoding/hex"

	"github.com/stellar/go/build"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

type dataMap map[string]map[string]string

// SimplifiedEnvelope is a collection of decoded data from actual tx envelope
// It is only used as a helper representation of TX
type SimplifiedEnvelope struct {
	SourceAccount string // pubkey
	Fee           uint32
	MemoHash      string
	// Can't omit any types and keys of data values
	DataValues          dataMap  // ["type:pubkey" : ["key" : "value"]]
	TotalOperationCount int      // All operations with data and without
	TxHash              [32]byte // Hash of the TX
}

func accountIDToString(acc *xdr.AccountId) (string, error) {
	return strkey.Encode(strkey.VersionByteAccountID, (*acc.Ed25519)[:])
}

// readDataValues aggregates all of the contract's saved data to a single data map
// Since data can be set for multiple accounts in one transaction the upper level of the map is an address
// This means that for each address there is a map of key-value elements
// {
//   addrA: {
//     key1: value1
//     key2: value2
//   }
//   addrB: {
//     key1: value2
//   }
// }
//
// The preservation of this data is needed because we want to know what account's data is updated
// This means that we are able to check if an attacker wants to corrupt his other trade account
func readDataValues(builder build.TransactionEnvelopeBuilder) (dataMap, error) {
	values := make(dataMap)
	for _, o := range builder.E.Tx.Operations {
		if o.Body.Type != xdr.OperationTypeManageData {
			continue
		}
		sourceAddr, err := accountIDToString(o.SourceAccount)
		if err != nil {
			return values, err
		}
		if values[sourceAddr] == nil {
			values[sourceAddr] = make(map[string]string)
		}
		values[sourceAddr][string(o.Body.ManageDataOp.DataName)] = string(*o.Body.ManageDataOp.DataValue)
	}
	return values, nil
}

func getOpCount(builder build.TransactionEnvelopeBuilder) int {
	return len(builder.E.Tx.Operations)
}

func parseMemoHash(memo xdr.Memo) (string, error) {
	if memo.Hash == nil {
		return "", nil
	}
	memoBin, err := memo.Hash.MarshalBinary()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(memoBin), nil
}

// Simplify parses the tx envelope and returns several representations of it
func Simplify(txEnvelope string) (*SimplifiedEnvelope, *build.TransactionEnvelopeBuilder, error) {
	eBuilder, err := ReadEnvelopeBuilder(txEnvelope)
	if err != nil {
		return nil, nil, err
	}
	sourceAddr, err := accountIDToString(&eBuilder.E.Tx.SourceAccount)
	if err != nil {
		return nil, nil, err
	}
	memo, err := parseMemoHash(eBuilder.E.Tx.Memo)
	if err != nil {
		return nil, nil, err
	}
	data, err := readDataValues(*eBuilder)
	if err != nil {
		return nil, nil, err
	}
	txb := build.TransactionBuilder{
		TX:                &eBuilder.E.Tx,
		NetworkPassphrase: "Test SDF Network ; September 2015",
	}
	hash, err := (&txb).Hash()
	if err != nil {
		return nil, nil, err
	}
	e := SimplifiedEnvelope{
		SourceAccount:       sourceAddr,
		Fee:                 uint32(eBuilder.E.Tx.Fee),
		MemoHash:            memo,
		TotalOperationCount: getOpCount(*eBuilder),
		DataValues:          data,
		TxHash:              hash,
	}
	return &e, eBuilder, nil
}
