package resolver

import (
	"time"

	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

func (s *UserHelpersSuite) TestParseExpireTime(c *C) {
	now := time.Now().UTC()
	// positive test for right time
	t := now.Add(time.Hour * 10)
	expTime1, errs := parseExpireTime(t.Format(time.RFC3339))
	c.Check(errs, IsNil)
	c.Check(expTime1, NotNil)

	t = now.Add(time.Second * 10)
	expTime2, errs := parseExpireTime(t.Format(time.RFC3339))
	c.Check(errs, IsNil)
	c.Check(expTime2, NotNil)

	// negative time for wrong time
	t = now.Add(time.Hour * 10)
	expTime3, errs := parseExpireTime(t.Format(time.RFC822))
	c.Check(errs, NotNil)
	c.Check(expTime3, IsNil)

	t = now.Add(time.Hour * -10)
	expTime4, errs := parseExpireTime(t.Format(time.RFC3339))
	c.Assert(errs, NotNil)
	c.Check(errs, ErrorContains, "Expire time can't be in the past")
	c.Check(expTime4, IsNil)
}
