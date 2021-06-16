// Package dal is data access layer for the project
package dal

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

// InsertTrade creates new trade between two people selected
func InsertTrade(ctx context.Context, db driver.Database, t *model.Trade) (driver.DocumentMeta, errstack.E) {
	if t == nil {
		return driver.DocumentMeta{}, errstack.NewDomain("Can't insert nil object")
	}
	t.CreatedAt = time.Now().UTC()
	return insertHasID(ctx, dbconst.ColTrades, t, db)
}

// GetTrade get a trade by ID
func GetTrade(ctx context.Context, db driver.Database, tradeID string) (*model.Trade, errstack.E) {
	var t model.Trade
	return &t, DBGetOneFromColl(ctx, &t, tradeID, dbconst.ColTrades, db)
}

// GetTradeOfDocument finds trade by document's ID
func GetTradeOfDocument(ctx context.Context, db driver.Database, docID string) (*model.Trade, errstack.E) {
	tradeDest := model.Trade{}
	err := DBQueryOne(
		ctx,
		&tradeDest,
		fmt.Sprintf(`
			for doc in %s filter doc._id == @full_doc_id
    		for trade in 1..1 outbound doc doc_edges return trade`,
			dbconst.ColDocs),
		map[string]interface{}{
			"full_doc_id": dbconst.ColDocs.FullID(docID)},
		db)
	return &tradeDest, err
}

// GetTrades gets all trade by user id
func GetTrades(ctx context.Context, db driver.Database, uid string) ([]model.Trade, errstack.E) {
	q := "for d in trades filter d.buyer.userID == @uid || d.seller.userID == @uid sort d.createdAt return d"
	vars := map[string]interface{}{
		"uid": uid}
	var ts []model.Trade
	err := DBQueryMany(ctx, &ts, q, vars, db)
	return ts, err
}

// GetAllTrades gets all trade data
func GetAllTrades(ctx context.Context, db driver.Database) ([]model.Trade, errstack.E) {
	q := "for d in trades sort d.createdAt return d"
	var ts []model.Trade
	return ts, DBQueryMany(ctx, &ts, q, nil, db)
}

// DeleteTradeData delete the selected trade data
func DeleteTradeData(ctx context.Context, db driver.Database, tradeID string) errstack.E {
	return deleteDoc(ctx, db, dbconst.ColTrades, tradeID)
}

// UpdateTrade updates trade
func UpdateTrade(ctx context.Context, db driver.Database, t *model.Trade) (driver.DocumentMeta, errstack.E) {
	return UpdateDoc(ctx, db, dbconst.ColTrades, t.ID, t)
}

// GetTradeAndStage retrieves Trade And Stage from DB
func GetTradeAndStage(ctx context.Context, db driver.Database, id model.TradeStagePath) (*model.Trade, *model.TradeStage, errstack.E) {
	t, errs := GetTrade(ctx, db, id.Tid)
	if errs != nil {
		return nil, nil, errs
	}
	stage, errs := t.GetStage(id.StageIdx)
	return t, stage, errs
}

// GetTradeStageDoc retrieves doc data from DB
func GetTradeStageDoc(ctx context.Context, db driver.Database, loc model.TradeStageDocPath) (*model.TradeStageDoc, errstack.E) {
	_, stage, errs := GetTradeAndStage(ctx, db, model.TradeStagePath{
		Tid: loc.Tid, StageIdx: loc.StageIdx})
	if errs != nil {
		return nil, errs
	}
	if uint(len(stage.Docs)) <= loc.StageDocIdx {
		return nil, errstack.NewReq("Stage Document Index out of range")
	}
	return &stage.Docs[loc.StageDocIdx], nil
}
