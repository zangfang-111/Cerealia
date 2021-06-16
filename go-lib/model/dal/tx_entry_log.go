package dal

import (
	"context"
	"fmt"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

// InsertTxLogEntry inserts tx with it's context
// @returns the inserted txEntry Document Meta object
func InsertTxLogEntry(ctx context.Context, db driver.Database, entry *model.TxLog, edge *model.TxLogEdge) (driver.DocumentMeta, errstack.E) {
	return InsertIntoGraph(ctx, db, dbconst.ColTxEntryLog, dbconst.ColTxEntryLogEdges, entry, edge)
}

// FindTxLogEntry finds a tx log entry using it's tradeID, stageIDx and docIDx
func FindTxLogEntry(ctx context.Context, db driver.Database, tradeID string, stageIdx, docIdx *uint) (*model.TxLog, errstack.E) {
	entryDest := model.TxLog{}
	err := DBQueryOne(
		ctx,
		&entryDest,
		fmt.Sprintf(`
for parent in %s
    for trade, edge, p in 1..1 outbound parent tx_entry_log_edges
        filter edge._to == @_to && edge.stageIdx == @stageIdx && edge.stageDocIdx == @stageDocIdx
        return parent`,
			dbconst.ColTxEntryLog),
		map[string]interface{}{
			"_to":         dbconst.ColTrades.FullID(tradeID),
			"stageIdx":    stageIdx,
			"stageDocIdx": docIdx,
		},
		db)
	return &entryDest, err
}

// UpdateTxLogEntry updates existing TxLogEntry
func UpdateTxLogEntry(ctx context.Context, db driver.Database, t model.TxLog) (driver.DocumentMeta, errstack.E) {
	return UpdateDoc(ctx, db, dbconst.ColTxEntryLog, t.ID, t)
}

// UpsertTxLogEntry creates or updates a TxLogEntry
func UpsertTxLogEntry(ctx context.Context, db driver.Database, entry *model.TxLog, edge *model.TxLogEdge) (driver.DocumentMeta, errstack.E) {
	found, ferr := FindTxLogEntry(ctx, db, edge.TradeID, edge.StageIdx, edge.StageDocIdx)
	if ferr != nil {
		return InsertTxLogEntry(ctx, db, entry, edge)
	}
	entry.SetID(found.ID)
	return UpdateTxLogEntry(ctx, db, *entry)
}
