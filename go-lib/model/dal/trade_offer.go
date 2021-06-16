package dal

import (
	"context"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"

	"bitbucket.org/cerealia/apps/go-lib/model"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

// GetTradeOffer fetches TradeOffer from DB by ID
func GetTradeOffer(ctx context.Context, db driver.Database, id string) (*model.TradeOffer, errstack.E) {
	var tof model.TradeOffer
	return &tof, DBGetOneFromColl(ctx, &tof, id, dbconst.ColTradeOffers, db)
}

// GetTradeOffers fetches all active offers from db.
func GetTradeOffers(ctx context.Context, db driver.Database) ([]model.TradeOffer, errstack.E) {
	q := "FOR d IN trade_offers FILTER d.closedAt == null RETURN d"
	var tofs []model.TradeOffer
	return tofs, DBQueryMany(ctx, &tofs, q, nil, db)
}

// InsertTradeOffer creates new tradeOffer
func InsertTradeOffer(ctx context.Context, db driver.Database, tof *model.TradeOffer) (*model.TradeOffer, errstack.E) {
	dm, errs := insertHasID(ctx, dbconst.ColTradeOffers, tof, db)
	if errs != nil {
		return nil, errs
	}
	if tof.DocID != nil {
		toe := model.TradeDocOfferEdgeDO{
			FullDocID:        dbconst.ColDocs.FullID(*tof.DocID),
			FullTradeOfferID: dm.ID.String(),
		}
		_, errs = InsertAny(ctx, dbconst.ColDocTradeOfferEdges, &toe, db)
	}
	return tof, errs
}

// CloseTradeOffer update the tradeOffer by setting closeAt with current time
func CloseTradeOffer(ctx context.Context, db driver.Database, offerID string) errstack.E {
	now := time.Now().UTC()
	diff := map[string]*time.Time{"closedAt": &now}
	_, err := UpdateDoc(ctx, db, dbconst.ColTradeOffers, offerID, diff)
	return errstack.WrapAsInf(err, "Failed to update tradeOffer")
}
