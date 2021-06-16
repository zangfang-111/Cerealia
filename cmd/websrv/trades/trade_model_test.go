package trades

import (
	. "gopkg.in/check.v1"
)

type TradeModelTest struct{}

var _ = Suite(&TradeModelTest{})

const sampleB64Envelope = "AAAAAMGf7pIPglfeRK1ag6MU1B+eZd8btX1bsunYU/O3EZLNAAABkAAIt3wAAAABAAAAAAAAAAAAAAAEAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAHdHJhZGVJRAAAAAABAAAABzE5OTMxMzUAAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAKc3RhZ2VJbmRleAAAAAAAAQAAAAEzAAAAAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAIZG9jSW5kZXgAAAABAAAAAjIyAAAAAAABAAAAAMGf7pIPglfeRK1ag6MU1B+eZd8btX1bsunYU/O3EZLNAAAACgAAAAdkb2NIYXNoAAAAAAEAAAAJdGVzdF9oYXNoAAAAAAAAAAAAAAG3EZLNAAAAQCi+vtxcqn3WKzX2shMFZE8m5H7bjaY8nwUUehHS0tzgpiAN1WKBtL8mq7HkHlZy+w2i8GJ1b8/ILU1ztZSTxA4="

func (s *TradeModelTest) TestParseUploadModel(c *C) {
	um := TradeStageDocInput{
		Tid:       "string",
		StageIdx:  1,
		Uploader:  "string",
		FileName:  "string",
		Note:      "string",
		ExpiresAt: "2006-01-02T15:04:05+07:00",
		SignedTx:  sampleB64Envelope,
	}
	_, errb := um.parseFields()
	c.Check(errb.NotNil(), Equals, false)
}
