package trades

import (
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txvalidation"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/build"
)

// validateDocAddTX validates document tx according to the validation criteria for approvals
func validateDocAddTX(input *TradeStageDocInputP, t *model.Trade, u *model.User, nextStageDocIdx uint, now time.Time) (*build.TransactionEnvelopeBuilder, *txvalidation.SimplifiedEnvelope, errstack.E) {
	docPath := model.TradeStageDocPath{
		StageIdx:    input.StageIdx,
		StageDocIdx: nextStageDocIdx,
	}
	if input.WithApproval {
		return txvalidation.ValidateDocAddExpireTX(
			input.SignedTx,
			docPath,
			t,
			u,
			input.Hash,
			input.ExpiresAtTime,
			now,
		)
	}
	return txvalidation.ValidateDocAddTX(
		input.SignedTx,
		docPath,
		t,
		u,
		input.Hash,
	)
}
