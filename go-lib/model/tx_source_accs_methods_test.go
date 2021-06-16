package model

import (
	"time"

	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

type TxSourceAccsSuite struct{}

var _ = Suite(&TxSourceAccsSuite{})

func (s *TxSourceAccsSuite) TestValidateOwnershipUser(c *C) {
	// Check for user ID
	lock := TXSourceAcc{
		PubKey:        "pubKey",
		Type:          "type",
		LockUserID:    "userID",
		LockTradeID:   "tradeID",
		LockExpiresAt: time.Now().Add(time.Minute),
	}
	// Valid user
	err := lock.MustBeValidFor("tradeID", "userID")
	c.Assert(err, IsNil)
	// Bad user, Good trade
	err = lock.MustBeValidFor("tradeID", "my user ID")
	c.Assert(err, ErrorContains, "Other user already locked this trade")
	// Good user, Bad trade
	err = lock.MustBeValidFor("bad trade ID", "userID")
	c.Assert(err, ErrorContains, "Other user already locked this trade")
	// Everything's bad
	err = lock.MustBeValidFor("bad trade ID", "bad user ID")
	c.Assert(err, ErrorContains, "Other user already locked this trade")
}

func (s *TxSourceAccsSuite) TestValidateOwnershipTime(c *C) {
	// Valid time
	lock := TXSourceAcc{
		PubKey:        "pubKey",
		Type:          "type",
		LockUserID:    "userID",
		LockTradeID:   "tradeID",
		LockExpiresAt: time.Now().Add(time.Minute),
	}
	err := lock.MustBeValidFor("tradeID", "userID")
	c.Assert(err, IsNil)
	// Invalid time
	lock = TXSourceAcc{
		PubKey:        "pubKey",
		Type:          "type",
		LockUserID:    "userID",
		LockTradeID:   "tradeID",
		LockExpiresAt: time.Now().Add(-5 * time.Second),
	}
	err = lock.MustBeValidFor("tradeID", "userID")
	c.Assert(err, ErrorContains, "Trade lock time expired. Try again.")
}

func (s *TxSourceAccsSuite) TestValidateOwnershipInvalid(c *C) {
	// Valid time
	invalidatedAt := time.Now()
	lock := TXSourceAcc{
		PubKey:         "pubKey",
		Type:           "type",
		LockUserID:     "userID",
		LockTradeID:    "tradeID",
		LockExpiresAt:  time.Now().Add(time.Minute),
		LockUnlockedAt: &invalidatedAt,
	}
	err := lock.MustBeValidFor("tradeID", "userID")
	c.Assert(err, ErrorContains, "Lock invalidated")
}
