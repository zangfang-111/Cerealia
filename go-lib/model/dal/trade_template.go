package dal

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

// GetTradeTemplates return all trade templates
func GetTradeTemplates(ctx context.Context, db driver.Database) ([]model.TradeTemplate, errstack.E) {
	q := "for d in trade_templates sort d._key return d"
	var tts []model.TradeTemplate
	return tts, DBQueryMany(ctx, &tts, q, nil, db)
}

// GetTradeTempate get a trade template by ID
func GetTradeTempate(ctx context.Context, db driver.Database, ttID string) (*model.TradeTemplate, errstack.E) {
	var tt model.TradeTemplate
	return &tt, DBGetOneFromColl(ctx, &tt, ttID, dbconst.ColTradeTemplates, db)
}
