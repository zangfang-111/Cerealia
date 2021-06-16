package model

import (
	"fmt"
	"testing"

	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ValidationSuite struct {
	sampleTrade Trade
}

type TradeMethodsSuite struct {
	sampleTrade Trade
}

var _ = Suite(&ValidationSuite{
	Trade{
		Buyer: TradeParticipant{
			UserID:            "buyer-id",
			KeyDerivationPath: "buyer-wallet-path",
		},
		Seller: TradeParticipant{
			UserID:            "seller-id",
			KeyDerivationPath: "seller-wallet-path",
		},
	}})

var _ = Suite(&TradeMethodsSuite{
	Trade{
		ID: "1234",
		Buyer: TradeParticipant{
			UserID:            "buyer-id",
			KeyDerivationPath: "buyer-wallet-path",
		},
		Seller: TradeParticipant{
			UserID:            "seller-id",
			KeyDerivationPath: "seller-wallet-path",
		},
	}})

func (s *ValidationSuite) TestSellerIsNotBuyerPositive(c *C) {
	vb := tradeValidatorBuilder{}
	vb.sellerShouldNotBeBuyer("123", "124")
	stack := vb.ToErrstackBuilder()
	c.Check(stack.NotNil(), Equals, false)
	c.Check(stack.Get(sellerIDField), IsNil, Comment("Should appear only on failure"))
	c.Check(stack.Get(buyerIDField), IsNil, Comment("Should appear only on failure"))
}

func (s *ValidationSuite) TestSellerIsNotBuyerNegative(c *C) {
	vb := tradeValidatorBuilder{}
	vb.sellerShouldNotBeBuyer("123", "123")
	stack := vb.ToErrstackBuilder()
	c.Check(stack.NotNil(), Equals, true)
	c.Check(stack.Get(sellerIDField), Equals, buyerIsSeller)
	c.Check(stack.Get(buyerIDField), Equals, buyerIsSeller)
}

func (s *ValidationSuite) TestCreatorIsSellerOrBuyerPositive(c *C) {
	vb := tradeValidatorBuilder{}
	vb.creatorIsSellerOrBuyer("123", "124", "123")
	stack := vb.ToErrstackBuilder()
	c.Check(stack.NotNil(), Equals, false)
	c.Check(stack.Get(sellerIDField), IsNil, Comment("Should appear only on failure"))
	c.Check(stack.Get(buyerIDField), IsNil, Comment("Should appear only on failure"))
	vb.creatorIsSellerOrBuyer("123", "124", "124")
	stack = vb.ToErrstackBuilder()
	c.Check(stack.NotNil(), Equals, false)
	c.Check(stack.Get(sellerIDField), IsNil, Comment("Should appear only on failure"))
	c.Check(stack.Get(buyerIDField), IsNil, Comment("Should appear only on failure"))
}

func (s *ValidationSuite) TestCreatorIsSellerOrBuyerNegative(c *C) {
	vb := tradeValidatorBuilder{}
	vb.creatorIsSellerOrBuyer("123", "124", "125")
	stack := vb.ToErrstackBuilder()
	c.Check(stack.NotNil(), Equals, true)
	c.Check(stack.Get(sellerIDField), Equals, creatorIsNotSellerOrBuyer)
	c.Check(stack.Get(buyerIDField), Equals, creatorIsNotSellerOrBuyer)
}

func (s *ValidationSuite) TestValidateNegative(c *C) {
	stack := NewTradeInput{}.Validate("template-id", &User{})
	c.Check(stack.NotNil(), Equals, true)
	c.Check(stack.Get(nameField), Equals, "validation.required")
	sellerStrOutput := fmt.Sprint(stack.Get(sellerIDField))
	sellerStrExpected := fmt.Sprint([]interface{}{"validation.required", "validation.buyer-is-seller"})
	c.Check(sellerStrOutput, DeepEquals, sellerStrExpected)
	buyerStrOutput := fmt.Sprint(stack.Get(buyerIDField))
	buyerStrExpected := fmt.Sprint([]interface{}{"validation.required", "validation.buyer-is-seller"})
	c.Check(buyerStrOutput, DeepEquals, buyerStrExpected)
}

func (s *ValidationSuite) TestUserCanModifyPositive(c *C) {
	erre := s.sampleTrade.CanBeModifiedBy(&User{ID: "buyer-id"})
	c.Assert(erre, IsNil, Comment("Buyer should be able to modify the trade"))
	erre = s.sampleTrade.CanBeModifiedBy(&User{ID: "seller-id"})
	c.Assert(erre, IsNil, Comment("Seller should be able to modify the trade"))
}

func (s *ValidationSuite) TestUserCanModifyNegative(c *C) {
	err := s.sampleTrade.CanBeModifiedBy(&User{ID: "stranger-id"})
	c.Assert(err, Equals, ErrUnauthorized, Comment("Stranger shouldn't be able to modify the trade"))
	c.Assert(err.IsReq(), IsTrue, Comment("Is request error"))
	err = s.sampleTrade.CanBeModifiedBy(&User{ID: "stranger-id-2"})
	c.Assert(err, Equals, ErrUnauthorized, Comment("Stranger shouldn't be able to modify the trade"))
	c.Assert(err.IsReq(), IsTrue, Comment("Is request error"))
}

func (s *ValidationSuite) TestIsDeletedOrClosed(c *C) {
	stage := TradeStage{
		Name:      "testStage",
		CloseReqs: []ApproveReq{},
		DelReqs:   []ApproveReq{},
	}
	stage.CloseReqs = append(stage.CloseReqs, ApproveReq{
		Status: ApprovalRejected,
	})
	c.Check(stage.IsDeletedOrClosed(), Not(IsTrue))

	stage.DelReqs = append(stage.DelReqs, ApproveReq{
		Status: ApprovalRejected,
	})
	c.Check(stage.IsDeletedOrClosed(), Not(IsTrue))

	stage.CloseReqs = append(stage.CloseReqs, ApproveReq{
		Status: ApprovalPending,
	})
	c.Check(stage.IsDeletedOrClosed(), Not(IsTrue))

	stage.DelReqs = append(stage.DelReqs, ApproveReq{
		Status: ApprovalPending,
	})
	c.Check(stage.IsDeletedOrClosed(), Not(IsTrue))

	stage.CloseReqs = append(stage.CloseReqs, ApproveReq{
		Status: ApprovalExpired,
	})
	c.Check(stage.IsDeletedOrClosed(), Not(IsTrue))

	stage.DelReqs = append(stage.DelReqs, ApproveReq{
		Status: ApprovalExpired,
	})
	c.Check(stage.IsDeletedOrClosed(), Not(IsTrue))

	stage.CloseReqs = append(stage.CloseReqs, ApproveReq{
		Status: ApprovalApproved,
	})
	c.Check(stage.IsDeletedOrClosed(), IsTrue)

	stage.DelReqs = append(stage.DelReqs, ApproveReq{
		Status: ApprovalApproved,
	})
	c.Check(stage.IsDeletedOrClosed(), IsTrue)
}

func (s *TradeMethodsSuite) TestFindTradeParticipant(c *C) {
	// Find buyer
	participant, err := s.sampleTrade.FindParticipant(&User{ID: "buyer-id"})
	c.Assert(err, IsNil)
	c.Assert(participant, NotNil)
	c.Check(participant.UserID, Equals, "buyer-id")
	// Find seller
	participant, err = s.sampleTrade.FindParticipant(&User{ID: "seller-id"})
	c.Assert(err, IsNil)
	c.Assert(participant, NotNil)
	c.Check(participant.UserID, Equals, "seller-id")
	// Find seller
	participant, err = s.sampleTrade.FindParticipant(&User{
		ID:              "moderator-id",
		Roles:           []UserRole{"moderator"},
		DefaultWalletID: "mod's default wallet id",
		StaticWallets: map[string]StaticWallet{
			"mod's default wallet id": StaticWallet{
				PubKey: "mod's pubkey",
			},
		},
	})
	c.Assert(err, IsNil)
	c.Assert(participant, NotNil)
	c.Check(participant.UserID, Equals, "moderator-id")
	// Negative
	participant, err = s.sampleTrade.FindParticipant(&User{ID: "hello world"})
	c.Assert(err, ErrorContains, "User 'hello world' is not a participant of trade '1234'")
	c.Assert(participant, IsNil)
}

func (s *S) TestKeyForUser(c *C) {
	key, err := tradeWithParticipants.FindPubKey(&User{ID: "1234"})
	c.Assert(err, IsNil)
	c.Check(key, Equals, SCAddr("pub-key-1234"))

	key, err = tradeWithParticipants.FindPubKey(&User{ID: "nonexistent"})
	c.Assert(err, ErrorContains, "User 'nonexistent' is not a participant of trade 'Trade-id'")
	c.Check(key, Equals, SCAddr(""))

	key, err = tradeWithParticipantsHD.FindPubKey(&User{ID: "abcde"})
	c.Assert(err, IsNil)
	c.Check(key, Equals, SCAddr("pubkey-abcde"))
	// Repeating should result in a same wallet output
	key, err = tradeWithParticipantsHD.FindPubKey(&User{ID: "abcde"})
	c.Assert(err, IsNil)
	c.Check(key, Equals, SCAddr("pubkey-abcde"))

	key, err = tradeWithParticipantsHD.FindPubKey(&User{ID: "1234"})
	c.Assert(err, IsNil)
	c.Check(key, Equals, SCAddr("pubkey-1234"))
}
