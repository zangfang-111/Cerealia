package main

import (
	"context"
	"os"
	"path/filepath"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	"bitbucket.org/cerealia/apps/go-lib/setup"
	dbs "bitbucket.org/cerealia/apps/go-lib/setup/arangodb"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/flag"
	bat "github.com/robert-zaremba/go-bat"
	"github.com/robert-zaremba/log15"
	"github.com/stellar/go/keypair"
)

var fixturesPath setup.PathFlag
var logger = log15.Root()

const tradeContractName = "trade contract"

var seedCollections = []dbconst.Col{
	dbconst.ColUsers,
	dbconst.ColOrganizations,
	dbconst.ColTradeTemplates,
	dbconst.ColTrades,
	dbconst.ColDocs,
	dbconst.ColDocEdges,
	dbconst.ColTxEntryLog,
	dbconst.ColTxEntryLogEdges,
	dbconst.ColDocTradeOfferEdges,
	dbconst.ColTradeOffers,
	dbconst.ColNotifications,
	dbconst.ColTxSourceAccs,
}

func main() {
	var defaultPath = filepath.Join(setup.RootDir, "fixtures")
	flag.Var(&fixturesPath, "fixtures-path",
		"path to fixtures. ["+defaultPath+"]")
	flag.Parse()
	if fixturesPath.Path == "" {
		fixturesPath.Path = defaultPath
	}

	ctx := context.Background()
	db, errs := dbs.GetDb(ctx)
	if errs != nil {
		logger.Error("Can not connect to database", errs)
	}

	for _, colName := range seedCollections {
		colNameS := string(colName)
		col, err := dal.GetCollByName(ctx, db, colName)
		if err != nil {
			logger.Error("Can not connect to "+colNameS+" collection", err)
		}
		documents, errs := GetSeedData(colName)
		if errs != nil {
			logger.Error("Can't load seed data", "collection_name", colNameS, errs)
			return
		}
		if documents == nil {
			continue
		}
		_, errSlice, err := col.CreateDocuments(ctx, documents)
		if err1 := errSlice.FirstNonNil(); err != nil || err1 != nil {
			logger.Error("Failed to seed data", "collection_name", colNameS, log15.Spew(errSlice), err)
			return
		}
	}
}

// GetSeedData gets userdata from users.json file and unmarshal it.
func GetSeedData(colName dbconst.Col) (interface{}, errstack.E) {
	seedFile := filepath.Join(fixturesPath.Path, string(colName)+".json")
	if _, err := os.Stat(seedFile); os.IsNotExist(err) {
		logger.Info("file doesn't exist. Skipping", "collection", colName)
		return nil, nil
	}
	switch colName {
	case dbconst.ColUsers:
		var items = new([]model.User)
		errs := bat.DecodeJSONFile(seedFile, items, logger)
		if errs != nil {
			return *items, errs
		}
		return *items, checkUserPubKey(items)
	case dbconst.ColOrganizations:
		var items = new([]model.Organization)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	case dbconst.ColTrades:
		var items = new([]model.Trade)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	case dbconst.ColTradeTemplates:
		var items = new([]model.TradeTemplate)
		errs := bat.DecodeJSONFile(seedFile, items, logger)
		if errs != nil {
			return *items, errs
		}
		return *items, checkTradeContract(items)
	case dbconst.ColDocs:
		var items = new([]model.Doc)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	case dbconst.ColDocEdges:
		var items = new([]model.TradeDocEdgeDO)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	case dbconst.ColTxEntryLog:
		var items = new([]model.TxLog)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	case dbconst.ColTxEntryLogEdges:
		var items = new([]model.TxLogEdgeDTO)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	case dbconst.ColTradeOffers:
		var items = new([]model.TradeOffer)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	case dbconst.ColDocTradeOfferEdges:
		var items = new([]model.TradeDocOfferEdgeDO)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	case dbconst.ColNotifications:
		var items = new([]model.Notification)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	case dbconst.ColTxSourceAccs:
		var items = new([]model.TXSourceAcc)
		return *items, bat.DecodeJSONFile(seedFile, items, logger)
	default:
		return nil, errstack.NewDomain("Wrong collection name")
	}
}

func checkTradeContract(t *[]model.TradeTemplate) errstack.E {
	res := []model.TradeTemplate{}
	for _, tt := range *t {
		if len(tt.Stages) == 0 {
			return errstack.NewDomain("the tradeTemplate " + tt.Name + " doesn't include any stages!")
		} else if tt.Stages[0].Name != tradeContractName {
			return errstack.NewDomain("the tradeTemplate " + tt.Name + " doesn't include trade contract as first stage!")
		} else {
			res = append(res, tt)
		}
	}
	*t = res
	return nil
}

func checkUserPubKey(t *[]model.User) errstack.E {
	res := []model.User{}
	for _, user := range *t {
		for _, wallet := range user.StaticWallets {
			_, err := keypair.Parse(wallet.PubKey)
			if err != nil {
				return errstack.NewInfF("the user %s %s wallet %s has invalid public key: '%s'", user.FirstName, user.LastName, wallet, wallet.PubKey)
			}
		}
		res = append(res, user)
	}
	*t = res
	return nil
}
