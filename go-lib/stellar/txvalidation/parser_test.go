package txvalidation

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	. "github.com/robert-zaremba/checkers"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	. "gopkg.in/check.v1"
)

const tx1 = "AAAAAMGf7pIPglfeRK1ag6MU1B+eZd8btX1bsunYU/O3EZLNAAABkAAIt3wAAAABAAAAAAAAAAAAAAAEAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAHdHJhZGVJRAAAAAABAAAABzE5OTMxMzUAAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAKc3RhZ2VJbmRleAAAAAAAAQAAAAEzAAAAAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAIZG9jSW5kZXgAAAABAAAAAjIyAAAAAAABAAAAAMGf7pIPglfeRK1ag6MU1B+eZd8btX1bsunYU/O3EZLNAAAACgAAAAdkb2NIYXNoAAAAAAEAAAAJdGVzdF9oYXNoAAAAAAAAAAAAAAG3EZLNAAAAQCi+vtxcqn3WKzX2shMFZE8m5H7bjaY8nwUUehHS0tzgpiAN1WKBtL8mq7HkHlZy+w2i8GJ1b8/ILU1ztZSTxA4="
const tx1TargetKey = "GDAZ73USB6BFPXSEVVNIHIYU2QPZ4ZO7DO2X2W5S5HMFH45XCGJM2VCB"
const tx2 = "AAAAAKlrUBzkbvqEZ6PwasQ+vPQBr5duH2PgjMhxxyU8O1GDAAABLAAYVCQAAAAHAAAAAAAAAAPughHebXP+wILuHw2hs//H002cTGYRGx+ZyezyOW1nWgAAAAMAAAABAAAAAKlrUBzkbvqEZ6PwasQ+vPQBr5duH2PgjMhxxyU8O1GDAAAACgAAAAZlbnRpdHkAAAAAAAEAAAAJc3RhZ2VfZG9jAAAAAAAAAQAAAACpa1Ac5G76hGej8GrEPrz0Aa+Xbh9j4IzIccclPDtRgwAAAAoAAAADaWR4AAAAAAEAAAADMDowAAAAAAEAAAAAqWtQHORu+oRno/BqxD689AGvl24fY+CMyHHHJTw7UYMAAAAKAAAACW9wZXJhdGlvbgAAAAAAAAEAAAAHcGVuZGluZwAAAAAAAAAAAYbGOBAAAABA8dw3A3s9IiqzO5mHTE/WLHJWk6VQvU/b9gIIdF4NghYPfHj1/b51LvQh/CIsr3+lHUeNWtK2Aww9xgilnP03Cw=="
const tx2TargetKey = "GCUWWUA44RXPVBDHUPYGVRB6XT2ADL4XNYPWHYEMZBY4OJJ4HNIYGOAE"
const samplePubKey = "GAJE2D632VIQ423C6JSREC6ZUVWJSNHPYQSKFPU7T6SGF7I3K2XBG5IJ"

func (s *TxValidationSuite) TestReadDataValues(c *C) {
	eBuilder, err := ReadEnvelopeBuilder(tx1)
	c.Assert(err, IsNil)
	vals, err := readDataValues(*eBuilder)
	c.Assert(err, IsNil)
	c.Check(
		vals,
		DeepEquals,
		dataMap{
			tx1TargetKey: {
				"tradeID":    "1993135",
				"stageIndex": "3",
				"docIndex":   "22",
				"docHash":    "test_hash"}},
	)
	c.Check(getOpCount(*eBuilder), Equals, 4)

	eBuilder, err = ReadEnvelopeBuilder(tx2)
	c.Assert(err, IsNil)
	vals2, err := readDataValues(*eBuilder)
	c.Assert(err, IsNil)
	c.Check(
		vals2,
		DeepEquals,
		dataMap{
			tx2TargetKey: {
				"entity":    model.TxTradeEntityStageDoc.String(),
				"idx":       "0:0",
				"operation": model.ApprovalPending.String()}})
	c.Check(getOpCount(*eBuilder), Equals, 3)
}

func (s *TxValidationSuite) TestSimplify(c *C) {
	se, eb, err := Simplify(tx1)
	c.Assert(err, IsNil)
	sourceAccID, err := accountIDToString(&eb.E.Tx.SourceAccount)
	c.Assert(err, IsNil)
	c.Check(sourceAccID, Equals, se.SourceAccount)
	c.Check(se, NotNil)
	c.Check(
		se.SourceAccount,
		Equals,
		tx1TargetKey)
	c.Check(se.Fee, Equals, uint32(400))
	c.Check(se.MemoHash, Equals, "")
	c.Check(se.TxHash, Equals, [32]uint8{0xf, 0x9f, 0x40, 0x39, 0x79, 0xce, 0xd1, 0xe8, 0xf7, 0x69, 0xc1, 0x80, 0xfd, 0xc0, 0xcf, 0x18, 0x9b, 0xbf, 0x4f, 0x82, 0x2b, 0xe7, 0x35, 0xa1, 0x9, 0xfa, 0x80, 0xae, 0xf2, 0xb4, 0x84, 0x33})
	c.Check(
		se.DataValues,
		DeepEquals,
		dataMap{
			tx1TargetKey: {
				"tradeID":    "1993135",
				"stageIndex": "3",
				"docIndex":   "22",
				"docHash":    "test_hash"}},
	)

	se, eb, err = Simplify(tx2)
	c.Assert(err, IsNil)
	sourceAccID, err = accountIDToString(&eb.E.Tx.SourceAccount)
	c.Assert(err, IsNil)
	c.Check(sourceAccID, Equals, se.SourceAccount)
	c.Check(se, NotNil)
	c.Check(
		se.SourceAccount,
		Equals,
		tx2TargetKey)
	c.Check(se.Fee, Equals, uint32(300))
	c.Check(se.MemoHash, Equals, "ee8211de6d73fec082ee1f0da1b3ffc7d34d9c4c66111b1f99c9ecf2396d675a")
	c.Check(se.TxHash, Equals, [32]uint8{0x80, 0xd, 0xc1, 0xa9, 0xc7, 0x90, 0xe1, 0xde, 0xee, 0xd7, 0xc1, 0xeb, 0xff, 0x34, 0xff, 0x5d, 0xe2, 0x53, 0x53, 0x12, 0xf1, 0xee, 0x3d, 0xd0, 0x76, 0xf9, 0xd6, 0x45, 0x3a, 0x7c, 0x18, 0x7e})
	c.Check(
		se.DataValues,
		DeepEquals,
		dataMap{
			tx2TargetKey: {
				"entity":    model.TxTradeEntityStageDoc.String(),
				"idx":       "0:0",
				"operation": model.ApprovalPending.String()}})
}

func (s *TxValidationSuite) TestAccountIdToString(c *C) {
	pair, err := keypair.Parse(samplePubKey)
	c.Assert(err, IsNil)
	aid := xdr.AccountId{}
	err = aid.SetAddress(pair.Address())
	c.Assert(err, IsNil)
	decoded, err := accountIDToString(&aid)
	c.Assert(err, IsNil)
	c.Check(samplePubKey, Equals, decoded)
}

func (s *TxValidationSuite) TestSimplifyNegative(c *C) {
	se, eb, err := Simplify(tx1 + "hello world")
	c.Assert(err, ErrorContains, "illegal base64 data at input byte 572")
	c.Check(se, IsNil)
	c.Check(eb, IsNil)
}
