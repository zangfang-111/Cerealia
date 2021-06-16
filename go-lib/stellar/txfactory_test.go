package stellar

import (
	"encoding/base64"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"

	"bitbucket.org/cerealia/apps/go-lib/stellar/secretkey"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txvalidation"
	. "github.com/robert-zaremba/checkers"
	"github.com/stellar/go/xdr"
	. "gopkg.in/check.v1"

	"github.com/stellar/go/build"
)

type TxFactorySuite struct {
	scAccs1, scAccs2, scAccsWrong txsource.SourceAccs
}

func init() {
	Suite(&TxFactorySuite{})
}

const testNetName = "noop"
const testMasterSecret = "SAHMDCRMCKGGMFENBKSUVYNMTM2GGNPR7DALQZ4SE6IIX2XJWJ6OQ4PG"
const testUserPublicKey = "GDOKCE5VFBB3CCPWG6HLQXW7AL4QELDSHHLABWDBYMRSI4UGYY4BBGS3"

var sampleExpireTime = time.Date(2020, time.February, 1, 0, 0, 0, 0, time.UTC)

// ReadEnvelopeBuilder reads base64 unsigned envelope and sign it
func ReadBuilderAndSign(d *Driver, envelopeBase64 string) (string, error) {
	binary, err := base64.StdEncoding.DecodeString(envelopeBase64)
	if err != nil {
		return "", err
	}
	envelope := xdr.Transaction{}
	err = envelope.UnmarshalBinary(binary)
	if err != nil {
		return "", err
	}
	builder := &build.TransactionBuilder{
		TX:                &envelope,
		NetworkPassphrase: d.Network.Passphrase.Passphrase,
	}
	txe, err := builder.Sign(testUserPublicKey)
	if err != nil {
		return "", err
	}
	return txe.Base64()
}

func (s *TxFactorySuite) SetUpSuite(c *C) {
	key1, err := secretkey.Parse("SDIQCRQVKMNWW6K3LF4MCSSWRI6GRDOII74B6X4IKO76WRFKKOL42XN2")
	c.Assert(err, IsNil)
	key2, err := secretkey.Parse("SAMAL7AWVNZKUFNTLSLAA5ZX6XWPD7A3KG2IOW2IZBH4YWJGSJGPC3S7")
	c.Assert(err, IsNil)
	key3, err := secretkey.Parse("SB5CTVK5SOQEWBWX2H6BDV4CADT5IPI4WKD3I7ERAGOBTROGAUU46Z5P")
	c.Assert(err, IsNil)
	s.scAccs1 = txsource.SourceAccs{
		TradeKeyPair: *key1,
		PoolAcc: model.ParsedTXSourceAcc{
			TXSourceAcc: model.TXSourceAcc{
				PubKey: model.SCAddr("GAQSTS6COMUHLTEJP7GWRAYOW5NPA5XBPWSYLDEHI3CQVZWF442V774V"),
			},
		},
	}
	s.scAccs2 = txsource.SourceAccs{
		TradeKeyPair: *key2,
		PoolAcc: model.ParsedTXSourceAcc{
			TXSourceAcc: model.TXSourceAcc{
				PubKey: model.SCAddr("GDENC4IQ6YSADBWQOVAIRK4PJMYL2HUOQ5SQBY6GKSNJHEYXPJWLCFQZ"),
			},
		},
	}
	s.scAccsWrong = txsource.SourceAccs{ // wrong scAddr
		TradeKeyPair: *key3,
		PoolAcc: model.ParsedTXSourceAcc{
			TXSourceAcc: model.TXSourceAcc{
				PubKey: model.SCAddr("TGAQSTS6COMUHLTEJP7GWRAYOW5NPA5XBPWSYLDEHI3CQVZWF442V774V"),
			},
		},
	}
}

func (s *TxFactorySuite) TestMkTradeDocApprovalTx(c *C) {
	testDriver, err := NewDriver(testNetName)
	c.Assert(err, IsNil, Comment("Failed to create new test stellar driver"))

	// stageDoc request and approve/reject tx builder test
	entityID := model.TradeStageDocPath{
		StageIdx:     1,
		StageDocIdx:  1,
		StageDocHash: "4db6c6c35d0f161ad199be03f03106674483ffb95c5aa342c7ac4e6f731b3758",
	}
	entity := model.TxTradeEntityStageDoc
	operation := model.ApprovalPending
	txStr, err := MkTradeDocApprovalExpireTx(testDriver, &s.scAccs2, entityID, entity, operation, sampleExpireTime)
	c.Assert(err, IsNil, Comment("Failed to make tx base64 string"))
	c.Check(txStr, Not(Equals), "", Comment("Generated tx is empty"))
	teb, err := ReadBuilderAndSign(testDriver, txStr)
	c.Assert(err, NotNil)
	_, _, err = txvalidation.Simplify(teb)
	c.Assert(err, NotNil)
	c.Check(txStr, Not(Equals), "", Comment("Generated tx seems not to be correct"))

	// Negative test for stageDoc request and approve/reject tx builder
	txStr, err = MkTradeDocApprovalExpireTx(testDriver, &s.scAccsWrong, entityID, entity, operation, sampleExpireTime)
	c.Assert(err, NotNil, Comment("Expected error doesn't happen"))
	c.Check(txStr, Equals, "", Comment("Generated tx should be empty"))
}

func (s *TxFactorySuite) TestMkTradeStageOperationTx(c *C) {
	testDriver, err := NewDriver(testNetName)
	c.Assert(err, IsNil, Comment("Failed to create new test stellar driver"))

	// stage close request and approve/reject tx builder test
	entityID := uint(2)
	entity := model.TxTradeEntityStageCloseReqs
	operation := model.ApprovalApproved
	txStr, err := MkTradeStageOperationTx(testDriver, &s.scAccs1, entityID, entity, operation)
	c.Assert(err, IsNil, Comment("Failed to make tx base64 string"))
	c.Check(txStr, Not(Equals), "", Comment("Generated tx is empty"))
	teb, err := ReadBuilderAndSign(testDriver, txStr)
	c.Assert(err, NotNil)
	_, _, err = txvalidation.Simplify(teb)
	c.Assert(err, NotNil)
	c.Check(txStr, Not(Equals), "", Comment("Generated tx seems not to be correct"))

	// stage add request and approve/reject tx builder test
	entity = model.TxTradeEntityStageAdd
	operation = model.ApprovalPending
	txStr, err = MkTradeStageOperationTx(testDriver, &s.scAccs1, entityID, entity, operation)
	c.Assert(err, IsNil, Comment("Failed to make tx base64 string"))
	c.Check(txStr, Not(Equals), "", Comment("Generated tx is empty"))
	teb, err = ReadBuilderAndSign(testDriver, txStr)
	c.Assert(err, NotNil)
	_, _, err = txvalidation.Simplify(teb)
	c.Assert(err, NotNil)
	c.Check(txStr, Not(Equals), "", Comment("Generated tx seems not to be correct"))

	// Negative test for stage close request and approve/reject tx builder
	txStr, err = MkTradeStageOperationTx(testDriver, &s.scAccsWrong, entityID, entity, operation)
	c.Assert(err, NotNil, Comment("Expected error doesn't happen"))
	c.Check(txStr, Equals, "", Comment("Generated tx should be empty"))
}

func (s *TxFactorySuite) TestMkTradeCloseTx(c *C) {
	testDriver, err := NewDriver(testNetName)
	c.Assert(err, IsNil, Comment("Failed to create new test stellar driver"))

	// trade close request and approve/reject tx builder test
	entityID := "2"
	entity := model.TxTradeEntityTradeCloseReqs
	operation := model.ApprovalApproved
	txStr, err := MkTradeCloseTx(testDriver, &s.scAccs1, entityID, entity, operation)
	c.Assert(err, IsNil, Comment("Failed to make tx base64 string"))
	c.Check(txStr, Not(Equals), "", Comment("Generated tx is empty"))
	teb, err := ReadBuilderAndSign(testDriver, txStr)
	c.Assert(err, NotNil)
	_, _, err = txvalidation.Simplify(teb)
	c.Assert(err, NotNil)
	c.Check(txStr, Not(Equals), "", Comment("Generated tx seems not to be correct"))

	// Negative test for trade close request and approve/reject tx builder
	txStr, err = MkTradeCloseTx(testDriver, &s.scAccsWrong, entityID, entity, operation)
	c.Assert(err, NotNil, Comment("Expected error doesn't happen"))
	c.Check(txStr, Equals, "", Comment("Generated tx should be empty"))
}
