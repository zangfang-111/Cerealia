package dal

import (
	"context"
	"path/filepath"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

// InsertTradeDoc inserts Trade Doc and associated edge.
// @returns the inserted doc Document Meta object
func InsertTradeDoc(ctx context.Context, db driver.Database, entity *model.Doc, edge model.TradeDocEdge) (driver.DocumentMeta, errstack.E) {
	return InsertIntoGraph(ctx, db, dbconst.ColDocs, dbconst.ColDocEdges, entity, edge)
}

// GetDoc fetches doc from DB
func GetDoc(ctx context.Context, db driver.Database, id string) (*model.Doc, errstack.E) {
	var d model.Doc
	return &d, DBGetOneFromColl(ctx, &d, id, dbconst.ColDocs, db)
}

// InsertOfferDoc inserts new trade Offer document
func InsertOfferDoc(ctx context.Context, db driver.Database, fi model.FileInfo, uid string) (*model.Doc, errstack.E) {
	d := model.Doc{
		Name:      fi.FileName,
		Type:      filepath.Ext(fi.FileName)[1:],
		Hash:      fi.Hash,
		URL:       fi.URL,
		CreatedBy: uid,
		CreatedAt: time.Now().UTC(),
	}
	_, err := insertHasID(ctx, dbconst.ColDocs, &d, db)
	return &d, errstack.WrapAsInfF(err, "Can't insert new doc into coll [%s]", dbconst.ColDocs)
}

// AssertTradeOfferDocExists checks if any tradeoffer-doc edge with the docid exists or not
func AssertTradeOfferDocExists(ctx context.Context, db driver.Database, docID string) errstack.E {
	var exists bool
	query := existsQuery(`FOR d IN doc_tradeoffer_edges FILTER d._from==@docID`)
	bindVars := map[string]interface{}{
		"docID": dbconst.ColDocs.FullID(docID),
	}
	err := DBQueryOne(ctx, &exists, query, bindVars, db)
	if err == nil && !exists {
		err = errstack.NewReq("TradeOffer-Doc edge with docid [" + docID + "] doesn't exist")
	}
	return err
}

// AssertDocExists checks if any doc with docid exists or not
func AssertDocExists(ctx context.Context, db driver.Database, docID string) errstack.E {
	var exists bool
	query := existsQuery(`FOR d IN docs FILTER d._key==@docID`)
	bindVars := map[string]interface{}{
		"docID": docID,
	}
	err := DBQueryOne(ctx, &exists, query, bindVars, db)
	if err == nil && !exists {
		err = errstack.NewReq("Doc with id [" + docID + "] doesn't exists")
	}
	return err
}
