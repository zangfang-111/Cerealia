// Cleanup script that removes an account associated with trade
// Sends all funds to fund destination account (-fund-destination-addr config or command line key)
//
// Usage examples:
//
// Mode 1: clean only single trade
// ./bin/stellar_cleanup -delete-trade-account 2119023 --signer <secret1> --signer <secret2> -fund-destination-addr <receiver>
//
// Mode 2: clean all trades from DB
// ./bin/stellar_cleanup --clean-all-trades --signer <key1> --signer <key2>
package main

import (
	"context"
	"fmt"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	txsourcedal "bitbucket.org/cerealia/apps/go-lib/model/dal/txsource"
	"bitbucket.org/cerealia/apps/go-lib/model/txlog"
	"bitbucket.org/cerealia/apps/go-lib/setup"
	dbs "bitbucket.org/cerealia/apps/go-lib/setup/arangodb"
	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/stellar/secretkey"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/flag"
	"github.com/robert-zaremba/log15"
	"github.com/robert-zaremba/log15/log15setup"
	"github.com/stellar/go/keypair"
)

var logger = log15.Root()

var noopLogger = txlog.NoopTxLogger{}
var noopSourceDriver = txsource.NoopDriver{}

func decodeAddress(seed string) (string, error) {
	pair, err := keypair.Parse(seed)
	if err != nil {
		return "", err
	}
	return pair.Address(), nil
}

func dismantleSingleTradeByID(ctx context.Context, db driver.Database, ld *stellar.WrappedDriver, tradeID, destAcc string, signSeeds []keypair.Full) error {
	t, err := dal.GetTrade(ctx, db, tradeID)
	if err != nil {
		logger.Error("Can't read trade "+tradeID, err)
		return err
	}
	return dismantleSingleTrade(ctx, txsourcedal.NewSourceLocker(db), ld, *t, destAcc, signSeeds)
}

func dismantleSingleTrade(ctx context.Context, locker txsourcedal.SourceLocker, ld *stellar.WrappedDriver, t model.Trade, fundDestAcc string, signSeeds []keypair.Full) error {
	dismantleKey, err := tradeToAddresses(ctx, locker, t)
	if err != nil {
		return err
	}
	mergeTx, err := newMergeTX(ld, dismantleKey.Address(), fundDestAcc)
	if err != nil {
		return err
	}
	allSigners := append(signSeeds, *dismantleKey)
	_, err = ld.SignAndSend(*mergeTx, allSigners...)
	if err != nil {
		return err
	}
	return nil
}

func convertKeypairsToStrings(signSecrets []string) (*[]keypair.Full, errstack.E) {
	var a []keypair.Full
	for _, secret := range signSecrets {
		parsed, err := secretkey.Parse(secret)
		if err != nil {
			return nil, errstack.WrapAsInfF(err, "Can't parse key")
		}
		a = append(a, *parsed)
	}
	return &a, nil
}

func dismantleAllTrades(ctx context.Context, db driver.Database, ld *stellar.WrappedDriver, fundDestAddr string, signSeeds []keypair.Full) errstack.E {
	ts, erre := GetTrades(ctx, db)
	if erre != nil {
		logger.Error("Can't read trade list", erre)
		return erre
	}
	successful := 0
	for i, t := range ts {
		logger.Info(fmt.Sprintf("Cleaning trade %s: [%d/%d]", t.ID, i, len(ts)))
		err := dismantleSingleTrade(ctx, txsourcedal.NewSourceLocker(db), ld, t, fundDestAddr, signSeeds)
		if err != nil {
			logger.Warn("Trade account can't be dismantled", err)
			continue
		}
		successful++
	}
	logger.Info(fmt.Sprintf("Success: [%d/%d]", successful, len(ts)))
	return nil
}

func main() {
	log15setup.MustLogger("cleaner_env", "cleaner", setup.GitVersion, "", "sec", "INFO", true)
	flag.Parse() // Loads flags defined in flags.go
	if *F.DeleteAccountOfTrade == "" && !(*F.DismantleAllTrades) {
		logger.Fatal(delTradeAccConfigKey + " was not provided")
	}
	fundDestAddr, err := decodeAddress(*F.FundDestinationAddr)
	if err != nil {
		logger.Fatal("Can't decode fund destination addr", err)
	}
	ctx := context.Background()
	db, erre := dbs.GetDb(ctx)
	if erre != nil {
		logger.Fatal("Can't get db", erre)
	}
	stellarDriver, err := stellar.NewDriver(*F.StellarNetwork)
	logDriver := stellarDriver.WithTxLogger(noopLogger, noopSourceDriver.IsAcquiredFn(ctx, "", ""))
	if err != nil {
		logger.Fatal("Can't build stellar.Driver", err)
	}
	additionalSigners, err := convertKeypairsToStrings(F.AdditionalSigners)
	if err != nil {
		logger.Fatal("Can't decode additional signers", err)
	}
	if *F.DismantleAllTrades {
		err = dismantleAllTrades(ctx, db, logDriver, fundDestAddr, *additionalSigners)
		if err != nil {
			logger.Fatal("Dismantle all returned an error", err)
		}
	} else {
		err = dismantleSingleTradeByID(ctx, db, logDriver, *F.DeleteAccountOfTrade, fundDestAddr, *additionalSigners)
		if err != nil {
			logger.Fatal("Dismantle single trade returned an error", err)
		}
	}
}
