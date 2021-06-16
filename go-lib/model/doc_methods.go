package model

import "bitbucket.org/cerealia/apps/go-lib/model/dbconst"

// SetID implements dal.HasID interface
func (d *Doc) SetID(id string) {
	d.ID = id
}

// ToEdgeDO accepts a full entity ID that involves the collection name
// it returns a DO (database object) that should be inserted directly into DB
func (de TradeDocEdge) ToEdgeDO(absoluteDocID string) interface{} {
	return TradeDocEdgeDO{
		FullTradeID: dbconst.ColTrades.FullID(de.TradeID),
		FullDocID:   absoluteDocID,
		StageIdx:    de.StageIdx,
		StageDocIdx: de.StageDocIdx,
	}
}
