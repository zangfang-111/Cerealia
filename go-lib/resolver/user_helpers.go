package resolver

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

func newTradeParticipant(ctx context.Context, db driver.Database, u *model.User) (*model.TradeParticipant, errstack.E) {
	wallet, err := u.DefaultWallet()
	if err != nil {
		return nil, errstack.WrapAsDomain(err, "Default wallet not found")
	}
	derivedKey, derivationPath, err := wallet.DeriveNewKey()
	if err != nil {
		return nil, errstack.WrapAsDomain(err, "Can't create a new derivation param")
	}
	if err = dal.ReplaceUser(ctx, db, u); err != nil {
		return nil, errstack.WrapAsInf(err, "Can't update user's data")
	}
	return &model.TradeParticipant{
		UserID:            u.ID,
		KeyDerivationPath: derivationPath,
		WalletID:          u.DefaultWalletID,
		PubKey:            derivedKey,
	}, nil
}
