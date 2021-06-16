package arangodb

import (
	"context"
	"testing"

	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestCreateDbConnection(c *C) {
	_ = dbURL.Set("http://root:birthday@localhost:8529/myNewDatabase")
	_, err := OpenDB(context.Background())
	c.Check(err, IsNil, Comment("Failed to create DB connection, plaese ensure the correct connection info"))
}
