package resolvertests

import (
	"net/http"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/resolver/testutil"
	. "github.com/robert-zaremba/checkers"
	"github.com/robert-zaremba/errstack"
	bat "github.com/robert-zaremba/go-bat"
	. "gopkg.in/check.v1"
)

var sampleDesc = "123456789012345678901234567890"

const validReason = "this is such a valid reason"
const testDocHash = "0de5620066bd089d06fbe45dc3bd80959502a75c865f915321bbbfb78f9d8f08"

func (s *TradeIntegrationSuite) TestMakeNewTrade(c *C) {
	var err error
	trade, err := s.noopResolver.Mutation().TradeCreate(s.buyer.Ctx, testutil.MakeTradeInput("test-trade", s.buyer.ID, s.seller.ID, &sampleDesc))
	c.Assert(err, IsNil)
	c.Check(trade, NotNil)
	notifications, err := s.noopResolver.Query().NotificationsTrade(s.seller.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.seller.ID)
	c.Check(notifications[0].EntityID, Contains, s.trade.FullID2())
	c.Check(notifications[0].Action, Equals, model.ApprovalApproved)
	foundTrade, err := testutil.GetTrade(s.buyer.Ctx, s.noopResolver, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(foundTrade, NotNil)
	c.Check(foundTrade.ID, Equals, s.trade.ID)
	c.Check(foundTrade.Buyer.UserID, Equals, s.buyer.ID)
	c.Check(foundTrade.Seller.UserID, Equals, s.seller.ID)
	status := fetchHTTPStatus("https://horizon-testnet.stellar.org/accounts/" + string(s.trade.SCAddr))
	c.Check("404 Not Found", Equals, status)

	// TxLogEntry should exist
	foundTxLog, err := dal.FindTxLogEntry(s.buyer.Ctx, s.db, foundTrade.ID, nil, nil)
	c.Check(foundTxLog, NotNil)
	c.Check(err, IsNil)
	c.Check(foundTxLog.TxStatus, Equals, model.TxStatusOk)
	c.Check(foundTxLog.Ledger, Equals, model.StellarLedger)
	c.Check(foundTxLog.CreatedBy, Equals, s.buyer.ID)
}

func fetchHTTPStatus(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return "fail: " + err.Error()
	}
	defer errstack.CallAndLog(logger, resp.Body.Close)
	return resp.Status
}

func (s *TradeIntegrationSuite) TestMakeNewTradeWithTestnetBlockchain(c *C) {
	var err error
	trade, err := s.testnetResolver.Mutation().TradeCreate(s.buyer.Ctx, testutil.MakeTradeInput("test-trade", s.buyer.ID, s.seller.ID, &sampleDesc))
	c.Assert(err, IsNil)
	c.Assert(trade, NotNil)
	foundTrade, err := testutil.GetTrade(s.buyer.Ctx, s.testnetResolver, trade.ID)
	c.Check(err, IsNil)
	c.Check(foundTrade, NotNil)
	c.Check(foundTrade.ID, Equals, trade.ID)
	c.Check(foundTrade.Buyer.UserID, Equals, s.buyer.ID)
	c.Check(foundTrade.Seller.UserID, Equals, s.seller.ID)
	status := fetchHTTPStatus("https://horizon-testnet.stellar.org/accounts/" + string(trade.SCAddr))
	c.Check("200 OK", Equals, status)
	// Let's add the stage
	stageInput := model.NewStageInput{
		Tid:         trade.ID,
		Owner:       "b",
		Name:        "My test stage",
		Description: "Description of my test stage",
		Reason:      "Integration tests",
	}
	stagePath := model.TradeStagePath{
		Tid:      trade.ID,
		StageIdx: 0,
	}
	rawStageTx, err := s.testnetResolver.Mutation().MkTradeStageAddTx(s.buyer.Ctx, stagePath, model.ApprovalApproved)
	c.Assert(err, IsNil)
	rawTxSigned, err := testutil.SignTx(
		*s.testnetDriver,
		rawStageTx,
		testutil.SampleUser1Seed,
	)
	stage, err := s.testnetResolver.Mutation().TradeStageAddReq(s.buyer.Ctx, stageInput, rawTxSigned, false)
	c.Assert(err, IsNil)
	c.Assert(stage.Name, Equals, "My test stage")

	// Two users concurrently acting on a same trade
	// first user
	stagePath1 := model.TradeStagePath{
		Tid:      trade.ID,
		StageIdx: 1,
	}
	rawStageTx1, err := s.testnetResolver.Mutation().MkTradeStageAddTx(s.buyer.Ctx, stagePath1, model.ApprovalApproved)
	c.Assert(err, IsNil)
	rawTxSigned1, err := testutil.SignTx(*s.testnetDriver, rawStageTx1, testutil.SampleUser1Seed)
	// second user
	stagePath2 := model.TradeStagePath{
		Tid:      trade.ID,
		StageIdx: 2,
	}
	rawStageTx2, err := s.testnetResolver.Mutation().MkTradeStageAddTx(s.seller.Ctx, stagePath2, model.ApprovalApproved)
	c.Assert(err, IsNil)
	rawTxSigned2, err := testutil.SignTx(*s.testnetDriver, rawStageTx2, testutil.SampleUser2Seed)
	// first user executes
	_, err = s.testnetResolver.Mutation().TradeStageAddReq(s.buyer.Ctx, stageInput, rawTxSigned1, false)
	c.Assert(err, IsNil)
	// second user executes
	_, err = s.testnetResolver.Mutation().TradeStageAddReq(s.seller.Ctx, stageInput, rawTxSigned2, false)
	c.Assert(err, IsNil)
}

func (s *TradeIntegrationSuite) TestMakeNewTradeNilDesc(c *C) {
	_, err := s.noopResolver.Mutation().TradeCreate(s.buyer.Ctx, testutil.MakeTradeInput("test-trade", s.buyer.ID, s.seller.ID, nil))
	c.Assert(err, IsNil)
}

func (s *TradeIntegrationSuite) TestMakeNewTradeWrongName(c *C) {
	trade, err := s.noopResolver.Mutation().TradeCreate(s.buyer.Ctx, model.NewTradeInput{
		TemplateID:  "1471516",
		Name:        "123",
		BuyerID:     s.buyer.ID,
		SellerID:    s.seller.ID,
		Description: &sampleDesc,
	})
	c.Assert(err, NotNil)
	c.Assert(trade, IsNil)
	c.Check(err, ErrorContains, "validation.insufficient-length")
}

func (s *TradeIntegrationSuite) TestMakeNewTradeSellerDoesNotExist(c *C) {
	trade, err := s.noopResolver.Mutation().TradeCreate(s.buyer.Ctx, testutil.MakeTradeInput("test-trade", s.buyer.ID, "qqqqq", &sampleDesc))
	c.Assert(err, NotNil)
	c.Assert(trade, IsNil)
	c.Check(err, ErrorContains, "DB: object not found [document not found]")
}

func (s *TradeIntegrationSuite) TestNewStageWithConfirmation(c *C) {
	c.Check(len(s.trade.Stages), Equals, 13)

	// Let's add the stage
	input := model.NewStageInput{
		Tid:         s.trade.ID,
		Owner:       "b",
		Name:        "My test stage",
		Description: "Description of my test stage",
		Reason:      "Integration tests",
	}
	stagePath := model.TradeStagePath{
		Tid:      s.trade.ID,
		StageIdx: 0,
	}
	mr := s.noopResolver.Mutation()
	rawTx, err := mr.MkTradeStageAddTx(s.buyer.Ctx, stagePath, model.ApprovalPending)
	c.Assert(err, IsNil)
	rawTxSigned, err := testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser1Seed,
	)
	c.Assert(err, IsNil)
	_, err = mr.TradeStageAddReq(s.buyer.Ctx, input, rawTxSigned, true)
	c.Assert(err, IsNil)
	notifications, err := s.noopResolver.Query().NotificationsTrade(s.seller.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.seller.ID)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stageAddReqs:0"))
	c.Check(notifications[0].Action, Equals, model.ApprovalPending)
	updatedTrade, err := testutil.GetTrade(s.buyer.Ctx, s.noopResolver, s.trade.ID)
	c.Assert(err, IsNil)
	c.Assert(updatedTrade, NotNil)
	c.Assert(len(updatedTrade.Stages), Equals, 13)
	c.Assert(len(updatedTrade.StageAddReqs), Equals, 1)
	newestStageRequest := updatedTrade.StageAddReqs[0]
	c.Check(newestStageRequest.Name, Equals, "My test stage")
	c.Check(newestStageRequest.Description, Equals, "Description of my test stage")
	c.Check(newestStageRequest.Owner, Equals, model.TradeActorB)
	c.Check(newestStageRequest.Status, Equals, model.ApprovalPending)
	c.Check(newestStageRequest.ReqBy, Equals, "3")
	c.Check(len(newestStageRequest.ReqAt.String()) > 30, Equals, true)
	c.Check(newestStageRequest.ReqTx, Matches, "^noop-driver-[0-9]+")
	c.Check(newestStageRequest.ApprovedBy, Equals, "")
	c.Check(newestStageRequest.ApprovedAt, IsNil)
	c.Check(newestStageRequest.ApprovedTx, Equals, "")
	c.Check(newestStageRequest.ReqReason, Equals, "Integration tests")
	c.Check(newestStageRequest.RejectReason, Equals, "")

	// Let other user approve
	rawTx, err = mr.MkTradeStageAddTx(s.seller.Ctx, stagePath, model.ApprovalApproved)
	c.Assert(err, IsNil)
	rawTxSigned, err = testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser2Seed,
	)
	c.Assert(err, IsNil)
	_, err = mr.TradeStageAddReqApprove(s.seller.Ctx, stagePath, rawTxSigned)
	c.Assert(err, IsNil)
	notifications, err = s.noopResolver.Query().NotificationsTrade(s.buyer.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.buyer.ID)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stageAddReqs:0"))
	c.Check(notifications[0].Action, Equals, model.ApprovalApproved)
	updatedTrade, err = testutil.GetTrade(s.buyer.Ctx, s.noopResolver, s.trade.ID)
	c.Assert(err, IsNil)
	c.Assert(updatedTrade, NotNil)
	c.Assert(len(updatedTrade.Stages), Equals, 14)
	c.Assert(len(updatedTrade.StageAddReqs), Equals, 1)
	newestStageRequest = updatedTrade.StageAddReqs[0]
	c.Check(newestStageRequest.Name, Equals, "My test stage")
	c.Check(newestStageRequest.Description, Equals, "Description of my test stage")
	c.Check(newestStageRequest.Owner, Equals, model.TradeActorB)
	c.Check(newestStageRequest.Status, Equals, model.ApprovalApproved)
	c.Check(newestStageRequest.ReqBy, Equals, "3")
	c.Check(len(newestStageRequest.ReqAt.String()) > 30, Equals, true)
	c.Check(newestStageRequest.ReqTx, Matches, "^noop-driver-[0-9]+")
	c.Check(newestStageRequest.ApprovedBy, Equals, "2")
	c.Check(len(newestStageRequest.ApprovedAt.String()) > 30, Equals, true)
	c.Check(newestStageRequest.ApprovedTx, Matches, "^noop-driver-[0-9]+")
	c.Check(newestStageRequest.ReqReason, Equals, "Integration tests")
	c.Check(newestStageRequest.RejectReason, Equals, "")
}

func (s *TradeIntegrationSuite) TestNewStageForModerator(c *C) {
	mr := s.noopResolver.Mutation()
	// Let's add the stage
	input := model.NewStageInput{
		Tid:         s.trade.ID,
		Owner:       "b",
		Name:        "My test stage",
		Description: "Description of my test stage",
		Reason:      "Integration tests",
	}
	stagePath := model.TradeStagePath{
		Tid:      s.trade.ID,
		StageIdx: 0,
	}
	rawTx, err := mr.MkTradeStageAddTx(s.moderator.Ctx, stagePath, model.ApprovalApproved)
	c.Assert(err, IsNil)
	rawTxSigned, err := testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUserModeratorSeed,
	)
	c.Assert(err, IsNil)
	_, err = mr.TradeStageAddReq(s.moderator.Ctx, input, rawTxSigned, false)
	c.Assert(err, IsNil)
	notifications, err := s.noopResolver.Query().NotificationsTrade(s.buyer.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.buyer.ID)
	c.Check(len(notifications[0].Receiver), Equals, 2)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stageAddReqs:0"))
	c.Check(notifications[0].Action, Equals, model.ApprovalSubmitted)
	updatedTrade, err := testutil.GetTrade(s.buyer.Ctx, s.noopResolver, s.trade.ID)
	c.Assert(err, IsNil)
	c.Assert(updatedTrade, NotNil)
	c.Assert(len(updatedTrade.StageAddReqs), Equals, 1)
	newestStageRequest := updatedTrade.StageAddReqs[0]
	c.Check(newestStageRequest.Name, Equals, "My test stage")
	c.Check(newestStageRequest.Status, Equals, model.ApprovalNil)
	c.Check(newestStageRequest.ReqBy, Equals, "4") // Moderator's userID
	c.Check(newestStageRequest.ReqReason, Equals, "Integration tests")

	newestStage := updatedTrade.Stages[len(updatedTrade.Stages)-1]
	c.Check(newestStage.Name, Equals, "My test stage")
	c.Check(newestStage.Description, Equals, "Description of my test stage")
	c.Check(newestStage.Owner, Equals, model.TradeActorB)
	c.Check(newestStage.Moderator.UserID, Equals, s.moderator.ID)
}

func (s *TradeIntegrationSuite) TestDeleteStageWithConfirmation(c *C) {
	var err error
	mr := s.noopResolver.Mutation()
	s.trade.Stages = append([]model.TradeStage{}, model.TradeStage{
		Name:        "testStage",
		Description: sampleDesc,
		Owner:       model.TradeActorB,
		Docs:        []model.TradeStageDoc{},
	})
	tradeStagePath := model.TradeStagePath{Tid: s.trade.ID, StageIdx: 0}
	_, err = dal.UpdateTrade(s.buyer.Ctx, s.db, s.trade)
	c.Check(err, IsNil)
	// test for trade stage delete request
	approveReq, err := mr.TradeStageDelReq(s.buyer.Ctx, tradeStagePath, validReason)
	c.Check(err, IsNil)
	c.Check(approveReq.Status, Equals, model.ApprovalPending)
	c.Check(approveReq.ReqBy, Equals, s.buyer.ID)
	c.Check(approveReq.ReqReason, Equals, validReason)

	notifications, err := s.noopResolver.Query().NotificationsTrade(s.seller.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.seller.ID)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stages:0", "delReqs:0"))
	c.Check(notifications[0].Action, Equals, model.ApprovalPending)

	// test for trade stage delete request rejecting
	_, err = mr.TradeStageDelReqReject(s.seller.Ctx, tradeStagePath, validReason)
	c.Check(err, IsNil)

	notifications, err = s.noopResolver.Query().NotificationsTrade(s.buyer.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.buyer.ID)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stages:0", "delReqs:0"))
	c.Check(notifications[0].Action, Equals, model.ApprovalRejected)
	_, err = mr.TradeStageDelReq(s.buyer.Ctx, tradeStagePath, validReason)
	c.Check(err, IsNil)

	// test for trade stage delete request approving
	approveReq, err = mr.TradeStageDelReqApprove(s.seller.Ctx, tradeStagePath)
	c.Check(err, IsNil)
	c.Check(approveReq.Status, Equals, model.ApprovalApproved)
	c.Check(approveReq.ApprovedBy, Equals, s.seller.ID)

	notifications, err = s.noopResolver.Query().NotificationsTrade(s.buyer.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.buyer.ID)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stages:0", "delReqs:1"))
	c.Check(notifications[0].Action, Equals, model.ApprovalApproved)
}

func (s *TradeIntegrationSuite) TestCloseStageReqWithPendingDocs(c *C) {
	var err error
	mr := s.noopResolver.Mutation()
	s.trade.Stages = append([]model.TradeStage{}, model.TradeStage{
		Name:        "testStage",
		Description: sampleDesc,
		Owner:       model.TradeActorB,
		Docs: []model.TradeStageDoc{model.TradeStageDoc{
			DocID:      "aaa1",
			Status:     model.ApprovalApproved,
			ExpiresAt:  time.Now().UTC(),
			ApprovedBy: s.seller.ID,
		}, model.TradeStageDoc{
			DocID:      "aaa2",
			Status:     model.ApprovalPending,
			ExpiresAt:  time.Now().UTC(),
			ApprovedBy: s.seller.ID,
		}, model.TradeStageDoc{
			DocID:      "aaa3",
			Status:     model.ApprovalRejected,
			ExpiresAt:  time.Now().UTC(),
			ApprovedBy: s.seller.ID,
		}},
	})
	tradeStagePath := model.TradeStagePath{Tid: s.trade.ID, StageIdx: 0}
	_, err = dal.UpdateTrade(s.buyer.Ctx, s.db, s.trade)
	c.Check(err, IsNil)
	rawTx, err := mr.MkTradeStageCloseTx(s.buyer.Ctx, tradeStagePath, model.ApprovalPending)
	c.Assert(err, IsNil)
	rawTxSigned, err := testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser1Seed,
	)
	c.Assert(err, IsNil)
	approveReq, err := mr.TradeStageCloseReq(s.buyer.Ctx, tradeStagePath, rawTxSigned, validReason)
	c.Check(err, IsNil)
	c.Check(approveReq.Status, Equals, model.ApprovalPending)
	c.Check(approveReq.ReqBy, Equals, s.buyer.ID)
	c.Check(approveReq.ReqReason, Equals, validReason)
	c.Check(approveReq.ReqActor, Equals, model.TradeActorB)

	notifications, err := s.noopResolver.Query().NotificationsTrade(s.seller.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.seller.ID)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stages:0", "closeReqs:0"))
	c.Check(notifications[0].Action, Equals, model.ApprovalPending)

	// test for trade stage close request rejecting
	rawTx, err = mr.MkTradeStageCloseTx(s.seller.Ctx, tradeStagePath, model.ApprovalRejected)
	c.Assert(err, IsNil)
	rawTxSigned, err = testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser2Seed,
	)
	c.Assert(err, IsNil)
	_, err = mr.TradeStageCloseReqReject(s.seller.Ctx, tradeStagePath, rawTxSigned, validReason)
	c.Assert(err, IsNil)

	notifications, err = s.noopResolver.Query().NotificationsTrade(s.buyer.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.buyer.ID)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stages:0", "closeReqs:0"))
	c.Check(notifications[0].Action, Equals, model.ApprovalRejected)

	// negative test for trade stage close request with only rejected documents
	s.trade.Stages = append([]model.TradeStage{}, model.TradeStage{
		Name:        "testStage",
		Description: sampleDesc,
		Owner:       model.TradeActorB,
		Docs: []model.TradeStageDoc{model.TradeStageDoc{
			DocID:      "aaa1",
			Status:     model.ApprovalRejected,
			ExpiresAt:  time.Now().UTC(),
			ApprovedBy: s.seller.ID,
		}, model.TradeStageDoc{
			DocID:      "aaa2",
			Status:     model.ApprovalRejected,
			ExpiresAt:  time.Now().UTC(),
			ApprovedBy: s.seller.ID,
		}},
	})
	_, err = dal.UpdateTrade(s.buyer.Ctx, s.db, s.trade)
	c.Assert(err, IsNil)
	stagePath := model.TradeStagePath{
		Tid:      s.trade.ID,
		StageIdx: 0,
	}
	rawTx, err = mr.MkTradeStageCloseTx(s.buyer.Ctx, stagePath, model.ApprovalApproved)
	c.Assert(err, IsNil)
	rawTxSigned, err = testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser1Seed,
	)
	c.Assert(err, IsNil)
	approveReq, err = mr.TradeStageCloseReq(s.buyer.Ctx, tradeStagePath, rawTxSigned, validReason)
	c.Check(err, NotNil)
	c.Check(approveReq, IsNil)
}

func (s *TradeIntegrationSuite) TestCloseStageWithConfirmation(c *C) {
	var err error
	s.trade.Stages = append([]model.TradeStage{}, model.TradeStage{
		Name:        "testStage",
		Description: sampleDesc,
		Owner:       model.TradeActorB,
		Docs: []model.TradeStageDoc{model.TradeStageDoc{
			DocID:      "aaa1",
			Status:     model.ApprovalApproved,
			ExpiresAt:  time.Now().UTC(),
			ApprovedBy: s.seller.ID,
		}},
	})
	tradeStagePath := model.TradeStagePath{Tid: s.trade.ID, StageIdx: 0}
	_, err = dal.UpdateTrade(s.buyer.Ctx, s.db, s.trade)
	c.Check(err, IsNil)

	// test for trade stage close request
	mr := s.noopResolver.Mutation()
	rawTx, err := mr.MkTradeStageCloseTx(s.buyer.Ctx, tradeStagePath, model.ApprovalPending)
	c.Assert(err, IsNil)
	rawTxSigned, err := testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser1Seed,
	)
	c.Assert(err, IsNil)
	approveReq, err := mr.TradeStageCloseReq(s.buyer.Ctx, tradeStagePath, rawTxSigned, validReason)
	c.Assert(err, IsNil)
	c.Check(approveReq.Status, Equals, model.ApprovalPending)
	c.Check(approveReq.ReqBy, Equals, s.buyer.ID)
	c.Check(approveReq.ReqReason, Equals, validReason)
	c.Check(approveReq.ReqActor, Equals, model.TradeActorB)

	// test for trade stage close request rejecting
	rawTx, err = mr.MkTradeStageCloseTx(s.seller.Ctx, tradeStagePath, model.ApprovalRejected)
	c.Assert(err, IsNil)
	rawTxSigned, err = testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser2Seed,
	)
	c.Assert(err, IsNil)
	_, err = mr.TradeStageCloseReqReject(s.seller.Ctx, tradeStagePath, rawTxSigned, validReason)
	c.Check(err, IsNil)

	// test for trade stage close request
	rawTx, err = mr.MkTradeStageCloseTx(s.buyer.Ctx, tradeStagePath, model.ApprovalPending)
	c.Assert(err, IsNil)
	rawTxSigned, err = testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser1Seed,
	)
	c.Assert(err, IsNil)
	_, err = mr.TradeStageCloseReq(s.buyer.Ctx, tradeStagePath, rawTxSigned, validReason)
	c.Check(err, IsNil)

	// test for trade stage close request approving
	rawTx, err = mr.MkTradeStageCloseTx(s.seller.Ctx, tradeStagePath, model.ApprovalApproved)
	c.Assert(err, IsNil)
	rawTxSigned, err = testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser2Seed,
	)
	c.Assert(err, IsNil)
	approveReq, err = mr.TradeStageCloseReqApprove(s.seller.Ctx, tradeStagePath, rawTxSigned)
	c.Assert(err, IsNil)
	c.Check(approveReq.Status, Equals, model.ApprovalApproved)
	c.Check(approveReq.ApprovedBy, Equals, s.seller.ID)
}
