package txvalidation

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TxValidationSuite struct {
	sampleExpireTime time.Time
}

const sampleExpireTimeStr = "2020-12-31T00:00:00.833791Z"

var _ = Suite(&TxValidationSuite{})

func (vs *TxValidationSuite) SetUpSuite(c *C) {
	var err error
	vs.sampleExpireTime, err = time.Parse(time.RFC3339, sampleExpireTimeStr)
	c.Assert(err, IsNil)
}
