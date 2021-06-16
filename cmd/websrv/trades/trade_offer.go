package trades

import (
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/log15"
)

// TradeOfferHandler is a route object for trade offers
type TradeOfferHandler struct{}

// HandlePostTradeOffer is to create a new trade offer
func (o TradeOfferHandler) HandlePostTradeOffer(c *routing.Context) error {
	ctx, db, u, errs := getAndCheckAuthUser(c)
	if errs != nil {
		return errs
	}
	err := c.Request.ParseMultipartForm(maxDocSize)
	if len(c.Request.MultipartForm.File["formfile"]) != 1 {
		return errstack.NewReqF("Expecting exactly one files")
	}
	if err != nil {
		return errstack.WrapAsReq(err, "Can't Parse the form data")
	}

	fi, err := storeDocFile(c.Request, 0, tradeDocDir)
	if err != nil {
		return err
	}
	logger.Debug("saving trade offer attached file", log15.Spew(fi))
	d, err := dal.InsertOfferDoc(ctx, db, fi, u.ID)
	if err != nil {
		return err
	}
	return c.Write(d.ID)
}

// HandleGetOfferDocByID retrieves a document for a tradeoffer by id
func (o TradeOfferHandler) HandleGetOfferDocByID(c *routing.Context) error {
	ctx, db, _, err := getAndCheckAuthUser(c)
	if err != nil {
		return err
	}
	docID := c.Param("docID")
	if err = dal.AssertTradeOfferDocExists(ctx, db, docID); err != nil {
		return err
	}
	return serveDocFile(ctx, c, db, docID)
}

// SetTradeOfferRoutes sets trade offer routes
func SetTradeOfferRoutes(routerG *routing.RouteGroup) {
	o := TradeOfferHandler{}
	routerG.Post("", o.HandlePostTradeOffer)
	routerG.Get("/docs/<docID>", o.HandleGetOfferDocByID)
}
