package model

import (
	"bitbucket.org/cerealia/apps/go-lib/validation"
	"github.com/robert-zaremba/errstack"
)

// ValidateInput validates the new tradeOfferInput data
func (ti TradeOfferInput) ValidateInput() errstack.E {
	errb := errstack.NewBuilder()
	validation.Required(ti.ComType, errb.Putter("ComType"))
	validation.Required(ti.Commodity, errb.Putter("Commodity"))
	validation.Required(ti.Currency, errb.Putter("Currency"))
	validation.Required(ti.Incoterm, errb.Putter("Incoterm"))
	validation.Required(ti.MarketLoc, errb.Putter("MarketLoc"))
	validation.Required(ti.Origin, errb.Putter("Origin"))
	validation.Required(ti.Price, errb.Putter("Price"))
	validation.Required(ti.PriceType, errb.Putter("PriceType"))
	validation.Required(ti.Quality, errb.Putter("Quality"))
	validation.Required(ti.Shipment, errb.Putter("Shipment"))
	validation.Required(ti.Vol, errb.Putter("Vol"))
	validation.Required(ti.OrgID, errb.Putter("OrgID"))
	if len(ti.Shipment) != 2 || !ti.Shipment[1].After(ti.Shipment[0]) {
		errb.Put("Shipment", "date is not correct")
	}
	return errb.ToReqErr()
}

// SetID implements dal.HasID interface
func (o *TradeOffer) SetID(id string) {
	o.ID = id
}
