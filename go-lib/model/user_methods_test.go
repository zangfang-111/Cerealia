package model

import (
	"reflect"

	. "github.com/robert-zaremba/checkers"
	rzc "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

type S struct{}

var _ = Suite(&S{})

func (s *S) TestCheckPassword(c *C) {
	c.Check(checkPassword("abcDEF@123"), IsNil, Comment("Failed to check password"))
	c.Check(checkPassword("15!@ERdvhe"), IsNil, Comment("Failed to check password"))
	c.Check(checkPassword("tue$dl@dfiDd"), IsNil, Comment("Failed to check password"))

	// negative test
	// negative test for short length password
	errs := checkPassword("As@3gs")
	c.Assert(errs, ErrorContains, shortLengthMessage, Comment("Expected error to check for short length"))

	// negative test for wrong password without lowercase letter
	errs = checkPassword("AB#1235VET")
	c.Assert(errs, ErrorContains, characterPwdMessage)

	// negative test for wrong password without uppercase letter
	errs = checkPassword("abc@#154ft")
	c.Assert(errs, ErrorContains, characterPwdMessage)

	// negative test for wrong password without special character
	errs = checkPassword("abctSD321")
	c.Assert(errs, ErrorContains, characterPwdMessage)

}

// User with no keys (something's not right)
var emptyUser = User{
	ID:                "1234",
	DefaultWalletID:   "qwe",
	StaticWallets:     map[string]StaticWallet{},
	HDCerealiaWallets: map[string]HDCerealiaWallet{},
}

// User with static default wallet
var staticUser = User{
	ID:              "1234",
	DefaultWalletID: "w1234",
	StaticWallets: map[string]StaticWallet{
		"w1234": StaticWallet{
			PubKey: "aabbccddee",
			Wallet: Wallet{Note: "my note"},
		},
	},
	HDCerealiaWallets: map[string]HDCerealiaWallet{},
}

// User with a static default wallet
var staticUserBadWalletID = User{
	ID:              "John",
	DefaultWalletID: "john-wallet-id",
	StaticWallets: map[string]StaticWallet{
		"my-nice-wallet": StaticWallet{
			PubKey: "aabbccddee",
			Wallet: Wallet{Note: "my note"},
		},
	},
	HDCerealiaWallets: map[string]HDCerealiaWallet{},
}

// User with hd default wallet
var hdUser = User{
	ID:              "abcde",
	DefaultWalletID: "one",
	StaticWallets:   map[string]StaticWallet{},
	HDCerealiaWallets: map[string]HDCerealiaWallet{
		"one": HDCerealiaWallet{
			DerivationIndex: 1336,
		},
		"two": HDCerealiaWallet{
			DerivationIndex: 840,
		},
		"three": HDCerealiaWallet{
			DerivationIndex: 532,
		},
	},
}

// User with hd default wallet
var hdAndStaticSameIDUser = User{
	ID:              "abcde",
	DefaultWalletID: "shadow",
	StaticWallets: map[string]StaticWallet{
		"shadow": StaticWallet{},
	},
	HDCerealiaWallets: map[string]HDCerealiaWallet{
		"shadow": HDCerealiaWallet{},
	},
}

func (s *S) TestGetDefaultWallet(c *C) {
	// Non existent wallet
	wallet, err := emptyUser.DefaultWallet()
	c.Assert(err, ErrorContains, "User '1234' does not have a wallet with ID: 'qwe'")
	c.Check(wallet, IsNil)
	// Static
	staticWallet, err := staticUser.DefaultWallet()
	c.Assert(err, IsNil)
	c.Check(staticWallet, NotNil)
	// Static with wallet wrong ID
	staticWallet, err = staticUserBadWalletID.DefaultWallet()
	c.Assert(err, ErrorContains, "User 'John' does not have a wallet with ID: 'john-wallet-id'")
	c.Check(staticWallet, IsNil)
	// HD wallet
	hdWallet, err := hdUser.DefaultWallet()
	c.Assert(err, IsNil)
	c.Check(hdWallet, NotNil)
	// HD wallets should shadow the static ones
	hdWallet, err = hdAndStaticSameIDUser.DefaultWallet()
	c.Assert(err, IsNil)
	c.Check(hdWallet, NotNil)
	c.Check(reflect.TypeOf(hdWallet), Equals, reflect.TypeOf(&HDCerealiaWallet{}))
}

func (s *S) TestWalletByID(c *C) {
	// Non existent wallet
	wallet, err := emptyUser.FindWallet(emptyUser.DefaultWalletID)
	c.Assert(err, rzc.ErrorContains, "User '1234' does not have a wallet with ID: 'qwe'")
	c.Check(wallet, IsNil)
	// Static
	staticWallet, err := staticUser.FindWallet(staticUser.DefaultWalletID)
	c.Assert(err, IsNil)
	c.Check(staticWallet, NotNil)
	// Static with wallet wrong ID
	staticWallet, err = staticUser.FindWallet("some unknown wallet")
	c.Assert(err, ErrorContains, "User '1234' does not have a wallet with ID: 'some unknown wallet'")
	c.Check(staticWallet, IsNil)
	// HD wallet
	hdWallet, err := hdUser.FindWallet("three")
	c.Assert(err, IsNil)
	c.Check(hdWallet, NotNil)
	// HD wallets should shadow the static ones
	hdWallet, err = hdAndStaticSameIDUser.FindWallet(hdAndStaticSameIDUser.DefaultWalletID)
	c.Assert(err, IsNil)
	c.Check(hdWallet, NotNil)
	c.Check(reflect.TypeOf(hdWallet), Equals, reflect.TypeOf(&HDCerealiaWallet{}))
}

var tradeWithParticipants = Trade{
	ID: "Trade-id",
	Buyer: TradeParticipant{
		UserID:   "1234",
		WalletID: "w1234",
		PubKey:   "pub-key-1234",
	},
	Seller: TradeParticipant{
		UserID:   "John",
		WalletID: "wJOHN",
		PubKey:   "pub-key-John",
	},
}

var tradeWithParticipantsHD = Trade{
	Buyer: TradeParticipant{
		UserID:   "1234",
		WalletID: "my-precious-wallet",
		PubKey:   "pubkey-1234",
	},
	Seller: TradeParticipant{
		UserID:   "abcde",
		WalletID: "one",
		PubKey:   "pubkey-abcde",
	},
}
