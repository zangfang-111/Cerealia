package dal

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	driver "github.com/arangodb/go-driver"
)

// GetCollByName abstracts db graph and returns a collection instead
// If we instantiate the edge collection without the graph object then the constraints won’t be loaded.
func GetCollByName(ctx context.Context, db driver.Database, colName dbconst.Col) (driver.Collection, error) {
	c := string(colName)
	switch colName {
	case dbconst.ColTrades:
	case dbconst.ColDocs:
		gr, err := db.Graph(ctx, string(dbconst.GraphDoc))
		if err != nil {
			return nil, err
		}
		return gr.VertexCollection(ctx, c)
	case dbconst.ColDocEdges:
		gr, err := db.Graph(ctx, string(dbconst.GraphDoc))
		if err != nil {
			return nil, err
		}
		coll, _, err := gr.EdgeCollection(ctx, c)
		return coll, err
	case dbconst.ColTxEntryLogEdges:
		gr, err := db.Graph(ctx, string(dbconst.GraphTxLog))
		if err != nil {
			return nil, err
		}
		coll, _, err := gr.EdgeCollection(ctx, c)
		return coll, err
	}
	return db.Collection(ctx, c)
}
