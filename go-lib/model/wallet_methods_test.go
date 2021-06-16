package model

import (
	//. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

type WalletSuite struct{}

var _ = Suite(&WalletSuite{})

func (s *WalletSuite) TestDeriveStatic(c *C) {
	w := StaticWallet{
		PubKey: "Hello, I am a key",
	}
	k, path, err := w.DeriveNewKey()
	c.Assert(err, IsNil)
	c.Assert(k, Equals, "Hello, I am a key")
	c.Assert(path, Equals, "")
	// Repeated attempt should produce the same key
	k, path, err = w.DeriveNewKey()
	c.Assert(err, IsNil)
	c.Assert(k, Equals, "Hello, I am a key")
	c.Assert(path, Equals, "")
}

func (s *WalletSuite) TestDeriveNewKeyHD(c *C) {
	w := HDCerealiaWallet{
		// TODO: add parameters
		DerivationIndex: 42858,
	}
	k, path, err := w.DeriveNewKey()
	c.Assert(err, IsNil)
	c.Check(k, Equals, "GD3EPS4EBOK6ZELDEN466I6EU4LW7TK6UL6INZRC6OKLZKYXXESS64VE")
	c.Check(w.DerivationIndex, Equals, 42859)
	c.Check(path, Equals, "m/44'/148'/42858'")
	// Repeated attempt should produce a new key
	//	k2, path2, err := w.DeriveNewKey()
	//	c.Assert(err, IsNil)
	//	c.Check(k2, Not(Equals), k) //TODO: make it work with proper derivation
	//	c.Check(w.DerivationIndex, Equals, 42860)
	//	c.Check(path2, Equals, "m/44'/148'/42859'")
	//	c.Check(path2, Not(Equals), path)
}
