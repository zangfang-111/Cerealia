package model

import (
	"fmt"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/stellar/secretkey"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/keypair"
)

// ParsedTXSourceAcc stores full key pair
type ParsedTXSourceAcc struct {
	TXSourceAcc
	KeyPair keypair.Full
}

// TxSourceAccTypeTrade represents a type for trades
const TxSourceAccTypeTrade TXSourceAccType = "trade"

// TxSourceAccTypePool represents a type for a pool
const TxSourceAccTypePool TXSourceAccType = "pool"

const lockAlreadyOwned = "Other user already locked this trade (%s)."
const lockNotOwned = "Lock is not owned by anyone."
const lockTimeExpired = "Trade lock time expired. Try again."
const lockInvalidated = "Lock invalidated."

// MustBeValidFor checks that user holds the lock
func (c *TXSourceAcc) MustBeValidFor(tradeID, userID string) errstack.E {
	if c.LockUserID == "" || c.LockTradeID == "" {
		return errstack.NewReq(lockNotOwned)
	}
	if c.LockUserID != userID {
		return errstack.NewReq(fmt.Sprintf(lockAlreadyOwned, c.LockUserID))
	}
	if c.LockTradeID != tradeID {
		return errstack.NewReq(fmt.Sprintf(lockAlreadyOwned, c.LockUserID))
	}
	now := time.Now()
	if c.LockExpiresAt.Before(now) {
		return errstack.NewReq(lockTimeExpired)
	}
	if c.LockUnlockedAt != nil && c.LockUnlockedAt.Before(now) {
		return errstack.NewReq(lockInvalidated)
	}
	return nil
}

// Parse returns TXSourceAcc with parsed keypair
func (c *TXSourceAcc) Parse() (*ParsedTXSourceAcc, errstack.E) {
	full, err := secretkey.Parse(string(c.SCSecret))
	if err != nil {
		return nil, errstack.WrapAsInf(err, "Can't parse keypair")
	}
	return &ParsedTXSourceAcc{
		TXSourceAcc: *c,
		KeyPair:     *full,
	}, nil
}
