package txlog

import (
	"testing"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/test/daltest"
	"github.com/stellar/go/xdr"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ValidationSuite struct{}

var _ = Suite(&ValidationSuite{})

var db = daltest.DBStub{}

func (s *ValidationSuite) TestNewForDoc(c *C) {
	stageID := uint(12385858)
	docID := uint(585858)
	logger := New(nil, db, "stellar", "trade-id", &stageID, &docID, "initiator-id")
	underlying, ok := logger.(txLogger)
	c.Check(ok, Equals, true)
	c.Check(underlying.db, Equals, db)
	c.Check(underlying.ledger, Equals, model.StellarLedger)
	c.Check(underlying.tradeID, Equals, "trade-id")
	c.Check(*underlying.stageID, Equals, stageID)
	c.Check(*underlying.docID, Equals, docID)
	c.Check(underlying.userID, Equals, "initiator-id")
}

func (s *ValidationSuite) TestMakeTxEntry(c *C) {
	stageID := uint(12385858)
	docID := uint(585858)
	sampleLogger := txLogger{
		db:      nil,
		tradeID: "sample tradeID",
		stageID: &stageID,
		docID:   &docID,
		userID:  "sample initiatorID",
		ledger:  "sample ledger",
	}
	e, err := makeTxEnvelope(534, "GCXWGTAHGQEAM46DXRABW5SEKX22DA3XVUKG3P7YJ2WELA7B2CLVPB6I")
	c.Assert(err, IsNil)
	entry := sampleLogger.makeTxEntry("raw tx string", "create 123", e)
	c.Check(entry.ID, Equals, "")
	c.Check(string(entry.TxStatus), Equals, "create 123")
	c.Check(string(entry.Ledger), Equals, "sample ledger")
	c.Check(entry.RawTx, Equals, "raw tx string")
	c.Check(entry.CreatedBy, Equals, "sample initiatorID")
	c.Check(entry.UpdatedAt, NotNil)
	c.Check(entry.Nonce, Equals, xdr.SequenceNumber(534))
	c.Check(entry.Notes, Equals, "")
	c.Check(entry.SourceAcc, Equals, "GCXWGTAHGQEAM46DXRABW5SEKX22DA3XVUKG3P7YJ2WELA7B2CLVPB6I")
}

func (s *ValidationSuite) TestMakeTxEntryEdge(c *C) {
	stageID := uint(12385858)
	docID := uint(585858)
	sampleLogger := txLogger{
		db:      nil,
		tradeID: "sample tradeID",
		stageID: &stageID,
		docID:   &docID,
		userID:  "sample initiatorID",
		ledger:  "sample ledger",
	}
	entry := sampleLogger.makeTxEntryEdge()
	c.Check(entry.TradeID, Equals, "sample tradeID")
	c.Check(*entry.StageIdx, Equals, stageID)
	c.Check(*entry.StageDocIdx, Equals, docID)
}
