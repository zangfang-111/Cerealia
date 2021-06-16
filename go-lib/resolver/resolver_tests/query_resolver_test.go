package resolvertests

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	. "gopkg.in/check.v1"
)

func (s *TradeIntegrationSuite) TestGetPubKey(c *C) {
	var err error
	tr := s.noopResolver.Trade()
	actorData, err := tr.ActorWallet(s.buyer.Ctx, s.trade)
	c.Assert(err, IsNil)
	c.Assert(actorData, DeepEquals, &model.TradeActorWallet{
		PubKey:   "GD3EPS4EBOK6ZELDEN466I6EU4LW7TK6UL6INZRC6OKLZKYXXESS64VE",
		KeyPath:  "",
		WalletID: "wallet-id",
	})
	actorData, err = tr.ActorWallet(s.seller.Ctx, s.trade)
	c.Assert(err, IsNil)
	c.Assert(actorData, DeepEquals, &model.TradeActorWallet{
		PubKey:   "GCZVKTLPQY54OGK5R3EEAU22Q7COES2XP5544H2C3CRNN7GGRL74IBUM",
		KeyPath:  "",
		WalletID: "default-user2-wallet",
	})
}
