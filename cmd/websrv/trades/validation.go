package trades

import (
	"bitbucket.org/cerealia/apps/go-lib/validation"
	"github.com/robert-zaremba/errstack"
)

// Validate validates errors in upload input
func (u *TradeStageDocInput) Validate() errstack.Builder {
	vb := validation.Builder{}
	vb.Required(tidField, u.Tid)
	vb.Required(expiresAtField, u.ExpiresAt)
	vb.Time(expiresAtField, u.ExpiresAt)
	// TODO security: check that tx has one signature, that data change operations exist
	vb.Required(signedTxField, u.SignedTx)
	//vb.Required(fileInfos, upload.FileInfos) // Not sure why this field exists
	return vb.ToErrstackBuilder()
}

// ValidateStageIdx validates stage index bounds
func (u *TradeStageDocInputP) ValidateStageIdx(stageLen int) errstack.Builder {
	vb := validation.Builder{}
	vb.IndexLessThan(stageIdxField, int(u.StageIdx), stageLen)
	return vb.ToErrstackBuilder()
}
