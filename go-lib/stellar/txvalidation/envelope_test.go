package txvalidation

import (
	"github.com/stellar/go/xdr"
	. "gopkg.in/check.v1"
)

const sampleB64Envelope = "AAAAAMGf7pIPglfeRK1ag6MU1B+eZd8btX1bsunYU/O3EZLNAAABkAAIt3wAAAABAAAAAAAAAAAAAAAEAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAHdHJhZGVJRAAAAAABAAAABzE5OTMxMzUAAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAKc3RhZ2VJbmRleAAAAAAAAQAAAAEzAAAAAAAAAQAAAADBn+6SD4JX3kStWoOjFNQfnmXfG7V9W7Lp2FPztxGSzQAAAAoAAAAIZG9jSW5kZXgAAAABAAAAAjIyAAAAAAABAAAAAMGf7pIPglfeRK1ag6MU1B+eZd8btX1bsunYU/O3EZLNAAAACgAAAAdkb2NIYXNoAAAAAAEAAAAJdGVzdF9oYXNoAAAAAAAAAAAAAAG3EZLNAAAAQCi+vtxcqn3WKzX2shMFZE8m5H7bjaY8nwUUehHS0tzgpiAN1WKBtL8mq7HkHlZy+w2i8GJ1b8/ILU1ztZSTxA4="

func (s *TxValidationSuite) TestCanReadEnvelope(c *C) {
	envelope, err := readEnvelope(sampleB64Envelope)
	c.Check(err, IsNil)
	c.Check(envelope, NotNil)
	c.Check(envelope.Signatures, NotNil)
	c.Check(len(envelope.Signatures), Equals, 1)
	c.Check(envelope.Tx, NotNil)
	c.Check(envelope.Tx.SourceAccount, NotNil)
	c.Check(envelope.Tx.SourceAccount.Type.String(), Equals, "PublicKeyTypePublicKeyTypeEd25519")
	c.Check(string(envelope.Tx.SourceAccount.Ed25519[:32]), Equals, "\xc1\x9f\xee\x92\x0f\x82W\xdeD\xadZ\x83\xa3\x14\xd4\x1f\x9ee\xdf\x1b\xb5}[\xb2\xe9\xd8S\xf3\xb7\x11\x92\xcd")
	c.Check(len(envelope.Tx.Operations), Equals, 4)

	c.Check(*envelope.Tx.Operations[0].SourceAccount, NotNil)
	c.Check(envelope.Tx.Operations[0].Body.ManageDataOp.DataName, Equals, xdr.String64("tradeID"))
	tradeIDValLen := len(*envelope.Tx.Operations[0].Body.ManageDataOp.DataValue)
	c.Check(string((*envelope.Tx.Operations[0].Body.ManageDataOp.DataValue)[:tradeIDValLen]), Equals, "1993135")

	c.Check(*envelope.Tx.Operations[3].SourceAccount, NotNil)
	c.Check(envelope.Tx.Operations[3].Body.ManageDataOp.DataName, Equals, xdr.String64("docHash"))
	docHashValLen := len(*envelope.Tx.Operations[3].Body.ManageDataOp.DataValue)
	c.Check(string((*envelope.Tx.Operations[3].Body.ManageDataOp.DataValue)[:docHashValLen]), Equals, "test_hash")
	// TODO: Validate source account and body of each operation
}

func (s *TxValidationSuite) TestCanReadEnvelopeBuilder(c *C) {
	eb, err := ReadEnvelopeBuilder(sampleB64Envelope)
	c.Check(err, IsNil)
	c.Check(eb, NotNil)
	c.Check(eb.E.Signatures, NotNil)
	c.Check(len(eb.E.Signatures), Equals, 1)
}
