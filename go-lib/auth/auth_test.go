package auth

import (
	"testing"

	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestCreateJWT(c *C) {
	samToken, err := CreateJWT("1")
	c.Check(err, IsNil, Comment("Failed to generate Sam's JWT token, please ensure the correct info"))
	userID, errs := Authorize(samToken)
	c.Check(userID, Equals, "1", Comment("Parsed userID is wrong"))
	c.Check(errs, IsNil, Comment("Failed to parse the user token"))

	benToken, err := CreateJWT("2")
	c.Check(err, IsNil, Comment("Failed to generate Ben's JWT token, please ensure the correct info"))
	userID, errs = Authorize(benToken)
	c.Check(userID, Equals, "2", Comment("Parsed userID is wrong"))
	c.Check(errs, IsNil, Comment("Failed to parse the user token"))

	// negative test
	antonToken, err := CreateJWT("10")
	c.Check(err, IsNil, Comment("Failed to generate Ben's JWT token, please ensure the correct info"))
	userID, errs = Authorize(antonToken)
	c.Check(userID, Not(Equals), "100", Comment("Parsed userID should be different with expected value but now same"))
	c.Check(errs, IsNil, Comment("Failed to parse the user token"))
}
