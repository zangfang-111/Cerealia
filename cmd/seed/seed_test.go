package main

import (
	"testing"

	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestGetSeedData(c *C) {
	for _, colname := range seedCollections {
		_, errs := GetSeedData(colname)
		c.Check(errs, IsNil, Comment("Can't seed data of collection "+colname))
	}
}
