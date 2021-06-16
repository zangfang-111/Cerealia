package txvalidation

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

const stageActionTX = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAABAAAAAAAAAAPFf9gr8+SN1JV1sO6FuvKHiMI8oi54YTIMobZf5i5ovQAAAAMAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAJc3RhZ2UuZG9jAAAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAADaWR4AAAAAAEAAAADMDowAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAHcGVuZGluZwAAAAAAAAAAARe5JS8AAABAHQJp0zsYbAjKcvWsT5R+00r3rRQGllqJn68UiPX7jfp6M7MtYljipr8ZccwJEr3K6nkviMh65y5fGNlXyaKFAw=="

const stageActionTXTradeAccountKey = "GC7U6TFCSS75IHEOJRLSYRCIUGEKBZ7GLWJHP6VEPX7H4HL6LGXOIN3J"
const stageActionTXSigner = "SCVZ2D7TXSH7RNIYJOQSJOBTJSWE7AMDOHSVBQJA3OZABTJQ5MC47XID"

const invalidTX = "AAAAAMGf7pIPglfeRK1ag6MU1B+eZd8btX1bsunYU/O3EZLNAAABkAAIt3wAAAABAAAAAAAAAAAAAAAEAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAHdHJhZGVJRAAAAAABAAAABzE5OTMxMzUAAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAKc3RhZ2VJbmRleAAAAAAAAQAAAAEzAAAAAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAIZG9jSW5kZXgAAAABAAAAAjIyAAAAAAABAAAAAMGf7pIPglfeRK1ag6MU1B+eZd8btX1bsunYU/O3EZLNAAAACgAAAAdkb2NIYXNoAAAAAAEAAAAJdGVzdF9oYXNoAAAAAAAAAAAAAAG3EZLNAAAAQCi+vtxcqn3WKzX2shMFZE8m5H7bjaY8nwUUehHS0tzgpiAN1WKBtL8mq7HkHlZy+w2i8GJ1b8/ILU1ztZSTxA4="
const invalidTXExcessData = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAEsAAQCyAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABmVudGl0eQAAAAAAAQAAAAlzdGFnZS5kb2MAAAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAANpZHgAAAAAAQAAAAMwOjAAAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAJb3BlcmF0aW9uAAAAAAAAAQAAAAdwZW5kaW5nAAAAAAEAAAAA9kfLhAuV7JFjI3nvI8SnF2/NXqL8huYi85S8qxe5JS8AAAAKAAAACW1hbGljaW91cwAAAAAAAAEAAAAXc3RlYWwgYXNzZXRzLCBvYnZpb3VzbHkAAAAAAAAAAAEXuSUvAAAAQK3CcY8/GzV9zsT4OxA9C8AIXjq80JW9+lBEHRuI2CmwaOTvLKHKiuNChdxd0HmeHPpRAFrGm+sbit2pMKKkdwA="
const invalidTXNonMatchingDataKey = "AAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAABLAAQCyAAAAACAAAAAAAAAAAAAAADAAAAAQAAAAC/T0yilL/UHI5MVyxESKGIoOfmXZJ3+qR9/n4dflmu5AAAAAoAAAAGZW50aXR5AAAAAAABAAAACXN0YWdlLmRvYwAAAAAAAAEAAAAAv09MopS/1ByOTFcsREihiKDn5l2Sd/qkff5+HX5ZruQAAAAKAAAABGlkeDEAAAABAAAAAzA6MAAAAAABAAAAAL9PTKKUv9QcjkxXLERIoYig5+Zdknf6pH3+fh1+Wa7kAAAACgAAAAlvcGVyYXRpb24AAAAAAAABAAAAB3BlbmRpbmcAAAAAAAAAAAEXuSUvAAAAQH1hu2DBLIWCx5cdr/3fJY27czkWLA66g54VAs6QbIy9AU+EMYFjBcwWkCVHRWXReiQAEoDFgMwlwBzq47enHAM="

var stageActionValidTrade = &model.Trade{
	SCAddr: stageActionTXTradeAccountKey,
	Buyer: model.TradeParticipant{
		UserID:   "test-id",
		WalletID: "my-wallet",
		PubKey:   stageActionTXSigner,
	},
	Seller: model.TradeParticipant{
		UserID:   "test-id-2",
		WalletID: "my-wallet-2",
		PubKey:   stageActionTXSigner,
	},
}

var stageActionValidUser = &model.User{
	ID: "test-id",
	StaticWallets: map[string]model.StaticWallet{
		"my-wallet": model.StaticWallet{
			PubKey: stageActionTXSigner,
		},
	},
}

func (s *TxValidationSuite) TestStageActionParseErr(c *C) {
	builder, sb, err := prevalidateTradeDataTx(stageActionValidTrade, stageActionValidUser, "hello world")
	c.Assert(err.ToErrstackBuilder().ToReqErr(), ErrorContains, txUnparsable)
	c.Check(builder, IsNil)
	c.Check(sb, IsNil)
}

func (s *TxValidationSuite) TestStageActionGreenPath(c *C) {
	builder, sb, errb := prevalidateTradeDataTx(stageActionValidTrade, stageActionValidUser, stageActionTX)
	c.Assert(errb.IsEmpty(), Equals, true)
	c.Check(builder, NotNil)
	c.Check(sb, NotNil)
}

func (s *TxValidationSuite) TestStageActionNoKey(c *C) {
	var StageActionUserNoKey = &model.User{}
	builder, sb, err := prevalidateTradeDataTx(stageActionValidTrade, StageActionUserNoKey, stageActionTX)
	c.Assert(err.ToErrstackBuilder().ToReqErr(), ErrorContains, txSignature)
	c.Check(builder, IsNil)
	c.Check(sb, IsNil)
}

func (s *TxValidationSuite) TestStageActionBadSignature(c *C) {
	var StageActionUserBadSigner = &model.User{
		DefaultWalletID: "my-wallet",
		StaticWallets: map[string]model.StaticWallet{
			"my-wallet": model.StaticWallet{
				PubKey: "GDAZ73USB6BFPXSEVVNIHIYU2QPZ4ZO7DO2X2W5S5HMFH45XCGJM2VCB",
			},
		},
	}
	builder, sb, err := prevalidateTradeDataTx(stageActionValidTrade, StageActionUserBadSigner, stageActionTX)
	c.Assert(err.ToErrstackBuilder().ToReqErr(), ErrorContains, txSignature)
	c.Check(builder, IsNil)
	c.Check(sb, IsNil)
}
