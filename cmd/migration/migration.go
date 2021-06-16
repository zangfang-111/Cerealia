package main

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	"bitbucket.org/cerealia/apps/go-lib/setup"
	"bitbucket.org/cerealia/apps/go-lib/setup/arangodb"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/flag"
	"github.com/robert-zaremba/log15"
)

func init() {
	flag.Parse()
}

var logger = log15.Root()
var db driver.Database

func main() {
	setup.FlagSimpleInit("db-migration", "")

	var err error
	client, dbname, err := arangodb.MkClient()
	if err != nil {
		logger.Fatal("Can't create ArangoDB client", err)
	}

	ctx := context.Background()
	if db, err = client.Database(ctx, dbname); driver.IsNotFound(err) {
		db, err = client.CreateDatabase(ctx, dbname, nil)
	}
	if err != nil {
		logger.Error("Can't connect / create DB", err)
	}
	createCollections(ctx)
	createIndexes(ctx)
	createGraphs(ctx)
	logger.Info("Database migrated successfully")
}

// GraphDef defines a graph
type GraphDef struct {
	From     dbconst.Col
	To       dbconst.Col
	Title    dbconst.Col
	EdgeColl dbconst.Col
}

var graphs = []GraphDef{
	{
		From:     dbconst.ColDocs,
		To:       dbconst.ColTrades,
		Title:    dbconst.GraphDoc,
		EdgeColl: dbconst.ColDocEdges,
	}, {
		From:     dbconst.ColTxEntryLog,
		To:       dbconst.ColTrades,
		Title:    dbconst.GraphTxLog,
		EdgeColl: dbconst.ColTxEntryLogEdges,
	}, {
		From:     dbconst.ColDocs,
		To:       dbconst.ColTradeOffers,
		Title:    dbconst.GraphDocTradeOffer,
		EdgeColl: dbconst.ColDocTradeOfferEdges,
	},
}

func createUsingGraphDef(ctx context.Context, graphDef GraphDef) {
	gr, err := db.CreateGraph(ctx, string(graphDef.Title), nil)
	if err != nil {
		if driver.IsConflict(err) {
			logger.Debug("Graph already exists", "name", graphDef.Title)
			return
		}
		logger.Fatal("Can't create a new graph", "name", graphDef.Title, err)
		return
	}
	logger.Info("Creating graph", "name", graphDef.Title)
	_, err = gr.CreateVertexCollection(ctx, string(graphDef.From))
	if err != nil {
		logger.Fatal("Can't create graph vertex collection", "from_edge", graphDef.From, err)
	}
	_, err = gr.CreateVertexCollection(ctx, string(graphDef.To))
	if err != nil {
		logger.Fatal("Can't create graph vertex collection", "to_edge", graphDef.To, err)
	}
	graphConstraints := driver.VertexConstraints{
		From: []string{string(graphDef.From)},
		To:   []string{string(graphDef.To)},
	}
	_, err = gr.CreateEdgeCollection(ctx, string(graphDef.EdgeColl), graphConstraints)
	if err != nil {
		logger.Fatal("Can't create graph edge collection", "name", graphDef.EdgeColl, err)
	}
}

func createGraphs(ctx context.Context) {
	for _, graphDef := range graphs {
		createUsingGraphDef(ctx, graphDef)
	}
}

func createCollections(ctx context.Context) {
	type collection struct {
		name dbconst.Col
		opts *driver.CreateCollectionOptions
	}
	defaultOpts := driver.CreateCollectionOptions{}
	collections := []collection{
		// Edge collections are omitted, they are created in createUsingGraphDef
		{dbconst.ColUsers, &driver.CreateCollectionOptions{}},
		{dbconst.ColOrganizations, &defaultOpts},
		{dbconst.ColTrades, &defaultOpts},
		{dbconst.ColTradeTemplates, &defaultOpts},
		{dbconst.ColDocs, &defaultOpts},
		{dbconst.ColTxEntryLog, &defaultOpts},
		{dbconst.ColTradeOffers, &defaultOpts},
		{dbconst.ColNotifications, &defaultOpts},
		{dbconst.ColTxSourceAccs, &driver.CreateCollectionOptions{
			WaitForSync: true,
		}},
	}

	for _, c := range collections {
		c.opts.Type = driver.CollectionTypeDocument
		nameStr := string(c.name)
		_, err := db.CreateCollection(ctx, nameStr, c.opts)
		if err != nil {
			if driver.IsConflict(err) {
				logger.Debug("Collection exists", "name", nameStr)
				continue
			}
			logger.Fatal("Can't create a new collection", "name", nameStr)
			return
		}
		logger.Info("Creating new collection", "name", nameStr)
	}
}

func createIndexes(ctx context.Context) {
	// https://docs.arangodb.com/3.4/Manual/Indexing/IndexBasics.html#indexing-array-values
	defaultOptions := driver.EnsureHashIndexOptions{Unique: true, Sparse: true}
	type index struct {
		collection dbconst.Col
		fields     []string
		options    *driver.EnsureHashIndexOptions
	}
	indexes := []index{
		index{dbconst.ColTradeTemplates, []string{"name"}, &defaultOptions},
		// we can't use this index, bug in arangodb: https://github.com/arangodb/arangodb/issues/8359
		// index{dbconst.ColUsers, []string{"emails[*]"}, &defaultOptions},
		index{dbconst.ColOrganizations, []string{"name"}, &defaultOptions},
		index{dbconst.ColOrganizations, []string{"address"}, &defaultOptions},
	}
	for _, idx := range indexes {
		col, err := db.Collection(ctx, string(idx.collection))
		if err != nil {
			logger.Fatal("Can't connect to collection", "name", idx.collection, err)
		}
		_, _, err = col.EnsureHashIndex(ctx, idx.fields, idx.options)
		if err != nil {
			logger.Fatal("Can't create hash index", err)
		}
	}
}
