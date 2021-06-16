package trades

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *ValidationSuite) TestValidateNegative(c *C) {
	um := TradeStageDocInput{}
	stack := um.Validate()
	c.Check(stack.NotNil(), Equals, true)
	c.Check(stack.Get(tidField), Equals, "validation.required")
	expiresAtExpected := fmt.Sprint([]interface{}{"validation.required", "validation.bad-format"})
	expiresAtResult := fmt.Sprint(stack.Get(expiresAtField))
	c.Check(expiresAtResult, Equals, expiresAtExpected)
	c.Check(stack.Get(signedTxField), Equals, "validation.required")
}

func (s *ValidationSuite) TestValidateStageIdx(c *C) {
	stack := (&TradeStageDocInputP{TradeStageDocInput: TradeStageDocInput{StageIdx: 999}}).ValidateStageIdx(100)
	c.Check(stack.NotNil(), Equals, true)
	c.Check(stack.Get(stageIdxField), Equals, "validation.index-out-of-bounds")
}
