package resolvertests

import (
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"

	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

var sampleExpireTime = time.Date(2020, 12, 31, 0, 0, 0, 323443, time.UTC)

func (s *TradeIntegrationSuite) TestMkTradeStageDocTx(c *C) {
	mr := s.noopResolver.Mutation()
	var err error
	id := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     0,
		StageDocIdx:  0,
		StageDocHash: "e0b976c1249699e30daa59007fb681939008a077a64856244233a3280acf8cee",
	}
	c.Assert(err, IsNil)
	time, err := time.Parse(time.RFC3339, "2096-01-02T15:04:05+00:00")
	c.Assert(err, IsNil)
	_, err = mr.MkTradeStageDocTx(s.buyer.Ctx, id, model.ApprovalPending, &time)
	c.Assert(err, IsNil)
	err = s.txSourceDriver.ReleaseFn(s.buyer.Ctx, s.trade.ID, s.buyer.ID)()
	c.Assert(err, IsNil)
	_, err = mr.MkTradeStageDocTx(s.buyer.Ctx, id, model.ApprovalApproved, &time)
	c.Assert(err, IsNil)
	err = s.txSourceDriver.ReleaseFn(s.buyer.Ctx, s.trade.ID, s.buyer.ID)()
	c.Assert(err, IsNil)
	_, err = mr.MkTradeStageDocTx(s.buyer.Ctx, id, model.ApprovalRejected, &time)
	c.Assert(err, IsNil)
	err = s.txSourceDriver.ReleaseFn(s.buyer.Ctx, s.trade.ID, s.buyer.ID)()
	c.Assert(err, IsNil)

	// negative test for wrong operationType
	_, err = mr.MkTradeStageDocTx(s.buyer.Ctx, id, "Bad operation", &time)
	c.Assert(err, NotNil, Comment("Expected an error for wrong operationType"))
}

func (s *TradeIntegrationSuite) TestMkTradeStageCloseTx(c *C) {
	var err error
	mr := s.noopResolver.Mutation()

	id := model.TradeStagePath{
		Tid:      s.trade.ID,
		StageIdx: 0,
	}
	c.Assert(err, IsNil)
	_, err = mr.MkTradeStageCloseTx(s.buyer.Ctx, id, "pending")
	c.Check(err, IsNil)
	err = s.txSourceDriver.ReleaseFn(s.buyer.Ctx, s.trade.ID, s.buyer.ID)()
	c.Assert(err, IsNil)
	_, err = mr.MkTradeStageCloseTx(s.buyer.Ctx, id, "approved")
	c.Check(err, IsNil)
	err = s.txSourceDriver.ReleaseFn(s.buyer.Ctx, s.trade.ID, s.buyer.ID)()
	c.Assert(err, IsNil)
	_, err = mr.MkTradeStageCloseTx(s.buyer.Ctx, id, "rejected")
	c.Check(err, IsNil)
	err = s.txSourceDriver.ReleaseFn(s.buyer.Ctx, s.trade.ID, s.buyer.ID)()
	c.Assert(err, IsNil)

	// negative test for wrong operationType
	_, err = mr.MkTradeStageCloseTx(s.buyer.Ctx, id, "Bad operation")
	c.Check(err, NotNil, Comment("Expected an error for wrong operationType"))
}

func (s *TradeIntegrationSuite) TestMkTradeCloseTx(c *C) {
	var err error
	mr := s.noopResolver.Mutation()
	_, err = mr.MkTradeCloseTx(s.buyer.Ctx, s.trade.ID, "pending")
	c.Assert(err, IsNil)
	err = s.txSourceDriver.ReleaseFn(s.buyer.Ctx, s.trade.ID, s.buyer.ID)()
	c.Assert(err, IsNil)
	_, err = mr.MkTradeCloseTx(s.buyer.Ctx, s.trade.ID, "approved")
	c.Assert(err, IsNil)
	err = s.txSourceDriver.ReleaseFn(s.buyer.Ctx, s.trade.ID, s.buyer.ID)()
	c.Assert(err, IsNil)
	_, err = mr.MkTradeCloseTx(s.buyer.Ctx, s.trade.ID, "rejected")
	c.Assert(err, IsNil)
	err = s.txSourceDriver.ReleaseFn(s.buyer.Ctx, s.trade.ID, s.buyer.ID)()
	c.Assert(err, IsNil)

	// negative test for wrong operationType
	_, err = mr.MkTradeCloseTx(s.buyer.Ctx, s.trade.ID, "Bad operation")
	c.Check(err, NotNil, Comment("Expected an error for wrong operationType"))
}
