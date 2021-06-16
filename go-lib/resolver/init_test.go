package resolver

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type UserHelpersSuite struct {
}

type TradeHelpersSuite struct {
}

var _ = Suite(&UserHelpersSuite{})

var _ = Suite(&TradeHelpersSuite{})
