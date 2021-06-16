package setup

import (
	"testing"

	rzcheck "github.com/robert-zaremba/checkers"
	gocheck "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gocheck.TestingT(t) }

type FlagsSuite struct{}

func init() {
	gocheck.Suite(&FlagsSuite{})
}

func (s *FlagsSuite) TestValidateDirExists(c *gocheck.C) {
	ff := PathFlag{}
	err := ff.Set("/tmp")
	c.Assert(err, gocheck.IsNil)
	c.Assert(ff.String(), gocheck.Equals, "/tmp")
}

func (s *FlagsSuite) TestValidatePermissionDenied(c *gocheck.C) {
	ff := PathFlag{}
	_ = ff.Set("/root/secret")
	err := ff.Check()
	c.Assert(err, gocheck.ErrorMatches, "stat /root/secret: permission denied")
}

func (s *FlagsSuite) TestValidateNoFile(c *gocheck.C) {
	ff := PathFlag{}
	_ = ff.Set("")
	err := ff.Check()
	c.Assert(err, rzcheck.ErrorContains, "File path can't be empty")
}

func (s *FlagsSuite) TestHandleBadDefaultWithPanic(c *gocheck.C) {
	ff := PathFlag{"hello-world"}
	err := ff.Check()
	c.Assert(err, gocheck.ErrorMatches, "stat hello-world: no such file or directory")
}

func (s *FlagsSuite) TestHandleDefault(c *gocheck.C) {
	ff := PathFlag{"/tmp"}
	c.Assert(ff.String(), gocheck.Equals, "/tmp")
	ff = PathFlag{"/"}
	c.Assert(ff.String(), gocheck.Equals, "/")
	err := ff.Set("/tmp")
	c.Assert(err, gocheck.IsNil)
	c.Assert(ff.String(), gocheck.Equals, "/tmp")
}
