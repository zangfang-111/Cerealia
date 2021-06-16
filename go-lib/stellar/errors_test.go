package stellar

import (
	"encoding/json"
	"testing"

	"github.com/facebookgo/stack"
	. "github.com/robert-zaremba/checkers"
	"github.com/stellar/go/clients/horizon"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type StellarErrors struct{}

func init() {
	Suite(&StellarErrors{})
}

func (s *StellarErrors) TestFormatStellarErrorNoErr(c *C) {
	hp := &horizon.Error{}
	hp.Problem.Title = "sample-problem-title"
	hp.Problem.Type = "sample-problem-type"
	hp.Problem.Extras = map[string]json.RawMessage{"lalala": []byte("['123', '345']")}
	e, _ := wrapErr(hp, "message").(herror)
	jsonBytes, outputError := e.MarshalJSON()
	c.Assert(outputError, IsNil)
	c.Check(
		string(jsonBytes),
		Equals,
		"{\"message\":\"message\",\"result_codes\":null,\"status\":0,\"title\":\"Horizon error: sample-problem-title\",\"type\":\"sample-problem-type\"}")
}

func (s *StellarErrors) TestFormatStellarErrorWithErr(c *C) {
	hp := &horizon.Error{}
	hp.Problem.Title = "sample-problem-title-1"
	hp.Problem.Type = "sample-problem-type-1"
	hp.Problem.Extras = map[string]json.RawMessage{"result_codes": []byte("{\"transaction\":\"tx_bad_auth\"}")}
	e, _ := wrapErr(hp, "message").(herror)
	jsonBytes, outputError := e.MarshalJSON()
	c.Assert(e, ErrorContains, "sample-problem-title-1")
	c.Assert(outputError, IsNil)
	c.Check(
		string(jsonBytes),
		Equals,
		"{\"message\":\"message\",\"result_codes\":{\"transaction\":\"tx_bad_auth\"},\"status\":0,\"title\":\"Horizon error: sample-problem-title-1\",\"type\":\"sample-problem-type-1\"}")
}

func (s *StellarErrors) TestSkipFrames(c *C) {
	var st stack.Stack
	// won't panic on empty stack
	skipInternalStack(st)

	// won't panic on not relevant stack
	st = stack.Stack{
		stack.Frame{File: "abc1"},
		stack.Frame{File: "abc2"},
		stack.Frame{File: pkgVendor + "/abc"}}
	stOut := skipInternalStack(st)
	c.Check(stOut, DeepEquals, st[:1], Commentf("When whole stack is not relevant it should keep the first frame.	"))

	st = stack.Stack{stack.Frame{File: pkgImport}}
	stOut = skipInternalStack(st)
	c.Check(stOut, DeepEquals, st,
		Commentf("It should work with a stack composed only with cerealia packages"))

	st = stack.Stack{
		stack.Frame{File: "abc"},
		stack.Frame{File: pkgVendor + "/abc"},
		stack.Frame{File: pkgImport + "/abc"},
		stack.Frame{File: "abc2"},
		stack.Frame{File: pkgVendor + "/abc"}}

	stOut = skipInternalStack(st)
	c.Check(stOut, DeepEquals, st[:3])
}
