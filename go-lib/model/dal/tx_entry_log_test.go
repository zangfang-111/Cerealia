package dal

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	. "gopkg.in/check.v1"
)

type TxEntryLogTest struct{}

var _ = Suite(&TxEntryLogTest{})

func (s *TxEntryLogTest) TestMakeTxEntryEdge(c *C) {
	stageIdx := uint(15129)
	docIdx := uint(5959)
	dto, ok := model.TxLogEdge{
		TradeID:     "sample-tradeID",
		StageIdx:    &stageIdx,
		StageDocIdx: &docIdx,
	}.ToEdgeDO("my_edge_coll/asd-123").(model.TxLogEdgeDTO)
	c.Assert(ok, Equals, true)
	c.Check(dto.FullTradeID, Equals, "trades/sample-tradeID")
	c.Check(dto.FullTxLogID, Equals, "my_edge_coll/asd-123")
	c.Check(*dto.StageIdx, Equals, stageIdx)
	c.Check(*dto.StageDocIdx, Equals, docIdx)
}
