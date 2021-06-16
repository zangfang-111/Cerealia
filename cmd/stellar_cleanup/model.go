package main

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	txsourcedal "bitbucket.org/cerealia/apps/go-lib/model/dal/txsource"
	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/stellar/secretkey"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
)

func tradeToAddresses(ctx context.Context, locker txsourcedal.SourceLocker, t model.Trade) (*keypair.Full, error) {
	lock, err := locker.FindByKey(ctx, t.SCAddr)
	if err != nil {
		return nil, err
	}
	return secretkey.Parse(string(lock.SCSecret))
}

func newMergeTX(ld *stellar.WrappedDriver, dismantleAcc string, pickUpAcc string) (*b.TransactionBuilder, errstack.E) {
	t, err := b.Transaction(
		b.SourceAccount{dismantleAcc},
		b.AutoSequence{SequenceProvider: ld.Client},
		ld.Network.Passphrase,
		b.ClearData("entity"),
		b.ClearData("idx"),
		b.ClearData("operation"),
		b.AccountMerge(
			b.Destination{pickUpAcc},
		),
	)
	return t, errstack.WrapAsDomain(err, "Can't construct 'merge' tx")
}

// GetTrades gets all trade data
func GetTrades(ctx context.Context, db driver.Database) ([]model.Trade, errstack.E) {
	query := "for d in trades return d"
	cursor, err := db.Query(ctx, query, make(map[string]interface{}))
	if err != nil {
		return nil, errstack.WrapAsInf(err, "Can't get trade data")
	}
	var ts []model.Trade
	defer errstack.CallAndLog(logger, cursor.Close) // idempotent, okay to call twice
	for {
		var t model.Trade
		_, err := cursor.ReadDocument(ctx, &t)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, errstack.WrapAsInf(err, "Can't read trade data")
		}
		ts = append(ts, t)
	}
	return ts, nil
}
