package stellar

import (
	"math/rand"
	"strconv"

	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/xdr"
)

const fakeNetName = "noop"

// NoopClient is fake stellar client to enable transaction without blockchain
type NoopClient struct{}

// SubmitTransaction is impl of fake stellar client and returns empty TransactionSuccess.
func (nc *NoopClient) SubmitTransaction(tx string) (horizon.TransactionSuccess, error) {
	return horizon.TransactionSuccess{
		Hash: "noop-driver-" + strconv.Itoa(rand.Int()),
	}, nil
}

// SequenceForAccount is impl of fake stellar client and returns empty SequenceNumber.
func (nc *NoopClient) SequenceForAccount(accountID string) (xdr.SequenceNumber, error) {
	return 0, nil
}
