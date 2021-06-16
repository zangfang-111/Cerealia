package validation

import (
	"testing"

	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ValidationSuite struct{}

var _ = Suite(&ValidationSuite{})

func (s *ValidationSuite) TestBuilderAdd(c *C) {
	vb := Builder{}
	vb.Append("key", "12345")
	c.Check(vb.Accumulated, HasLen, 1)
	c.Check(vb.Accumulated[0].Key, Equals, "key")
	c.Check(vb.Accumulated[0].Value, Equals, "12345")
	vb.Append("key 2", "123456")
	c.Check(vb.Accumulated, HasLen, 2)
	c.Check(vb.Accumulated[1].Key, Equals, "key 2")
	c.Check(vb.Accumulated[1].Value, Equals, "123456")
}

func (s *ValidationSuite) TestRequiredNegative(c *C) {
	vb := Builder{}
	vb.Required("field-name", "")
	c.Check(vb.Accumulated, HasLen, 1)
	c.Check(vb.Accumulated[0].Key, Equals, "field-name")
	c.Check(vb.Accumulated[0].Value, Equals, required)
}

func (s *ValidationSuite) TestRequiredPositive(c *C) {
	vb := Builder{}
	vb.Required("field-name", "I am here")
	c.Check(vb.Accumulated, HasLen, 0)
	vb.Required("field-name 1", "I am here as well")
	c.Check(vb.Accumulated, HasLen, 0)
	vb.Required("field-name 2", "0")
	c.Check(vb.Accumulated, HasLen, 0)
}

func (s *ValidationSuite) TestMinLengthNegative(c *C) {
	vb := Builder{}
	vb.MinLength("field-name", "value", 10)
	c.Check(vb.Accumulated, HasLen, 1)
	c.Check(vb.Accumulated[0].Key, Equals, "field-name")
	c.Check(vb.Accumulated[0].Value, Equals, insufficientLength)
	vb.MinLength("field-name", "", 1)
	c.Check(vb.Accumulated, HasLen, 2)
	c.Check(vb.Accumulated[1].Key, Equals, "field-name")
	c.Check(vb.Accumulated[1].Value, Equals, required)
}

func (s *ValidationSuite) TestMinLengthPositive(c *C) {
	vb := Builder{}
	vb.MinLength("field-name", "I am here", 5)
	c.Check(vb.Accumulated, HasLen, 0)
	vb.MinLength("field-name 1", "I am here as well", 10)
	c.Check(vb.Accumulated, HasLen, 0)
	vb.MinLength("field-name 2", "0", 1)
	c.Check(vb.Accumulated, HasLen, 0)
	vb.MinLength("field-name", "", 0)
	c.Check(vb.Accumulated, HasLen, 0)
}

func (s *ValidationSuite) TestToErrstackBuilder(c *C) {
	vb := Builder{}
	vb.Required("good-field-name", "I am here")
	vb.Required("bad-field-name", "")
	vb.Required("field-name 1", "I am here as well")
	errb := vb.ToErrstackBuilder()
	c.Check(errb.Get("good-field-name"), IsNil, Comment("Validation should succeed"))
	c.Check(errb.Get("bad-field-name"), Equals, required)
	c.Check(errb.Get("field-name 1"), IsNil, Comment("Validation should succeed"))
}

func (s *ValidationSuite) TestUniquePositive(c *C) {
	vb := Builder{}
	vb.Unique("name-a", "name-b", "val-a", "val-b", "custom.message")
	c.Check(vb.Accumulated, HasLen, 0)
}

func (s *ValidationSuite) TestUniqueNegative(c *C) {
	vb := Builder{}
	vb.Unique("name-a", "name-b", "qqqqq", "qqqqq", "custom.message")
	c.Check(vb.Accumulated, HasLen, 2)
	c.Check(vb.Accumulated[0].Key, Equals, "name-a")
	c.Check(vb.Accumulated[0].Value, Equals, "custom.message")
	c.Check(vb.Accumulated[1].Key, Equals, "name-b")
	c.Check(vb.Accumulated[1].Value, Equals, "custom.message")
	vb.Unique("name-w", "name-q", "", "", "custom.message.1")
	c.Check(vb.Accumulated, HasLen, 4)
	c.Check(vb.Accumulated[2].Key, Equals, "name-w")
	c.Check(vb.Accumulated[2].Value, Equals, "custom.message.1")
	c.Check(vb.Accumulated[3].Key, Equals, "name-q")
	c.Check(vb.Accumulated[3].Value, Equals, "custom.message.1")
}

func (s *ValidationSuite) TestIndexPositive(c *C) {
	vb := Builder{}
	vb.IndexLessThan("index-field", 12, 15)
	c.Check(vb.Accumulated, HasLen, 0)
	vb.IndexLessThan("index-field", 2, 3)
	c.Check(vb.Accumulated, HasLen, 0)
}

func (s *ValidationSuite) TestIndexNegative(c *C) {
	vb := Builder{}
	vb.IndexLessThan("index-field", 15, 14)
	c.Check(vb.Accumulated, HasLen, 1)
	c.Check(vb.Accumulated[0].Key, Equals, "index-field")
	c.Check(vb.Accumulated[0].Value, Equals, indexOutOfBounds)
	vb.IndexLessThan("index-field-1", -13, 14)
	c.Check(vb.Accumulated, HasLen, 2)
	c.Check(vb.Accumulated[1].Key, Equals, "index-field-1")
	c.Check(vb.Accumulated[1].Value, Equals, indexOutOfBounds)
	vb.IndexLessThan("index-field-2", 0, 0)
	c.Check(vb.Accumulated, HasLen, 3)
	c.Check(vb.Accumulated[2].Key, Equals, "index-field-2")
	c.Check(vb.Accumulated[2].Value, Equals, indexOutOfBounds)
}

func (s *ValidationSuite) TestTimePositive(c *C) {
	vb := Builder{}
	vb.Time("time-field", "2012-11-01T22:08:41+50:00")
	c.Check(vb.Accumulated, HasLen, 0)
}

func (s *ValidationSuite) TestTimeNegative(c *C) {
	vb := Builder{}
	vb.Time("time-field", "3")
	c.Check(vb.Accumulated, HasLen, 1)
	c.Check(vb.Accumulated[0].Key, Equals, "time-field")
	c.Check(vb.Accumulated[0].Value, Equals, badFormat)
	vb.Time("time-field-1", "3h")
	c.Check(vb.Accumulated, HasLen, 2)
	c.Check(vb.Accumulated[1].Key, Equals, "time-field-1")
	c.Check(vb.Accumulated[1].Value, Equals, badFormat)
}

func (s *ValidationSuite) TestMatchPositive(c *C) {
	vb := Builder{}
	vb.Match("match-field", "2012-11-01T22:08:41+50:00", "2012-11-01T22:08:41+50:00")
	c.Check(vb.Accumulated, HasLen, 0)
}

func (s *ValidationSuite) TestMatchNegative(c *C) {
	vb := Builder{}
	vb.Match("match-field", "3", "4")
	c.Check(vb.Accumulated, HasLen, 1)
	c.Check(vb.Accumulated[0].Key, Equals, "match-field")
	c.Check(vb.Accumulated[0].Value, Equals, doesNotMatch)
	vb.Match("match-field-1", "3h", "3")
	c.Check(vb.Accumulated, HasLen, 2)
	c.Check(vb.Accumulated[1].Key, Equals, "match-field-1")
	c.Check(vb.Accumulated[1].Value, Equals, doesNotMatch)
}

func (s *ValidationSuite) TestIsErrorNoError(c *C) {
	vb := Builder{}
	empty := vb.IsEmpty()
	c.Check(empty, Equals, true)
}

func (s *ValidationSuite) TestIsErrorWithError(c *C) {
	vb := Builder{}
	vb.Append("error-id", "something's broken")
	empty := vb.IsEmpty()
	c.Check(empty, Equals, false)
}
