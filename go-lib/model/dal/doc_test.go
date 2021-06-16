package dal

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	. "gopkg.in/check.v1"
)

type DocSuite struct{}

var _ = Suite(&DocSuite{})

func (s *DocSuite) TestMakeDocEdge(c *C) {
	stageIdx := uint(123132)
	docIdx := uint(993)
	dto, ok := model.TradeDocEdge{
		TradeID:     "sample-tradeID",
		StageIdx:    stageIdx,
		StageDocIdx: docIdx,
	}.ToEdgeDO("my_doc_coll/asd-123").(model.TradeDocEdgeDO)
	c.Assert(ok, Equals, true)
	c.Check(dto.FullTradeID, Equals, "trades/sample-tradeID")
	c.Check(dto.FullDocID, Equals, "my_doc_coll/asd-123")
	c.Check(dto.StageIdx, Equals, stageIdx)
	c.Check(dto.StageDocIdx, Equals, docIdx)
}
