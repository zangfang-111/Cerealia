package trades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/resolver/testutil"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txvalidation"
	routing "github.com/go-ozzo/ozzo-routing"
	. "github.com/robert-zaremba/checkers"
	bat "github.com/robert-zaremba/go-bat"
	. "gopkg.in/check.v1"
)

const testDocHash = "0de5620066bd089d06fbe45dc3bd80959502a75c865f915321bbbfb78f9d8f08"

var sampleDesc = "123456789012345678901234567890"

func (s *TradeIntegrationSuite) TestCloseStageReqApprove(c *C) {
	mr := s.noopResolver.Mutation()
	// New doc
	docHash := "f308fc02ce9172ad02a7d75800ecfc027109bc67987ea32aba9b8dcc7b10150e"
	docPath := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     0,
		StageDocIdx:  0,
		StageDocHash: docHash,
	}
	newDocTx, err := s.noopResolver.Mutation().MkTradeStageDocTx(s.seller.Ctx, docPath, model.ApprovalPending, &s.sampleExpireTime)
	c.Assert(err, IsNil)
	signedNewDocTx, err := testutil.SignTx(*s.noopDriver, newDocTx, testutil.SampleUser2Seed)
	c.Assert(err, IsNil)
	_, err = UploadDoc(
		s.seller.Ctx,
		s.noopDocHandler,
		UploadDocInput{
			StageIdx:     0,
			Data:         "test",
			TradeID:      s.trade.ID,
			ExpiresAt:    s.sampleExpireTimeStr,
			SignedTX:     signedNewDocTx,
			DocHash:      docHash,
			WithApproval: true,
		})
	c.Assert(err, IsNil)

	notifications, err := s.noopResolver.Query().NotificationsTrade(s.buyer.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.buyer.ID)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stages:0", "docs:0"))
	c.Check(notifications[0].Action, Equals, model.ApprovalPending)

	// Approve doc
	docApprovalTx, err := mr.MkTradeStageDocTx(s.buyer.Ctx, docPath, model.ApprovalApproved, nil)
	c.Assert(err, IsNil)
	signedDocApproveTx, err := testutil.SignTx(*s.noopDriver, docApprovalTx, testutil.SampleUser1Seed)
	c.Assert(err, IsNil)
	c.Assert(signedDocApproveTx, NotNil)
	tradeStageDoc, err := testutil.ApproveDoc(s.buyer.Ctx, s.noopResolver, docPath, signedDocApproveTx)
	c.Assert(err, IsNil)
	c.Assert(tradeStageDoc, NotNil)

	// Let other user approve
	approveInput := model.TradeStagePath{
		Tid:      s.trade.ID,
		StageIdx: 0,
	}
	rawStageTx, err := mr.MkTradeStageAddTx(s.seller.Ctx, approveInput, model.ApprovalRejected)
	c.Assert(err, IsNil)
	rawStageTxSigned, err := testutil.SignTx(*s.noopDriver, rawStageTx, testutil.SampleUser2Seed)
	c.Assert(err, IsNil)

	// Let's try approving a non-existent close request of the stage
	_, err = mr.TradeStageCloseReqApprove(s.seller.Ctx, approveInput, rawStageTxSigned)
	c.Check(err, ErrorContains, "No close request in this stage")

	// Create close request
	rawStageCloseTx, err := mr.MkTradeStageCloseTx(s.seller.Ctx, approveInput, model.ApprovalPending)
	c.Assert(err, IsNil)
	rawStageCloseTxSigned2, err := testutil.SignTx(*s.noopDriver, rawStageCloseTx, testutil.SampleUser2Seed)
	c.Assert(err, IsNil)
	_, err = mr.TradeStageCloseReq(s.seller.Ctx, approveInput, rawStageCloseTxSigned2, "reason 12345")
	c.Assert(err, IsNil)

	// Attempt to close the stage by the same user
	rawStageCloseApproveTx, err := mr.MkTradeStageCloseTx(s.buyer.Ctx, approveInput, model.ApprovalApproved)
	c.Assert(err, IsNil)
	rawStageCloseApproveTxSigned2, err := testutil.SignTx(*s.noopDriver, rawStageCloseApproveTx, testutil.SampleUser2Seed)
	c.Assert(err, IsNil)
	_, err = mr.TradeStageCloseReqApprove(s.seller.Ctx, approveInput, rawStageCloseApproveTxSigned2)
	c.Assert(err, ErrorContains, "You can't approve your own request")

	// Approve closing of the stage
	rawStageCloseTxSigned1, err := testutil.SignTx(*s.noopDriver, rawStageCloseApproveTx, testutil.SampleUser1Seed)
	c.Assert(err, IsNil)
	_, err = mr.TradeStageCloseReqApprove(s.buyer.Ctx, approveInput, rawStageCloseTxSigned1)
	c.Assert(err, IsNil)
}

func (s *TradeIntegrationSuite) TestUploadDocNoAuth(c *C) {
	stageDoc, err := UploadDoc(context.Background(), s.noopDocHandler, UploadDocInput{WithApproval: true})
	c.Assert(err, ErrorContains, "Authentication required")
	c.Check(stageDoc, IsNil)
}

func (s *TradeIntegrationSuite) TestUploadDocWithAuth(c *C) {
	var err error
	u, err := testutil.Login(s.noopResolver, testutil.SampleUser1)
	c.Assert(err, IsNil)
	stageDoc, err := UploadDoc(u.Ctx, s.noopDocHandler, UploadDocInput{WithApproval: true})
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "(?s).*tid.*validation.required.*")
	c.Check(err, ErrorMatches, "(?s).*expiresAt.*")
	c.Check(err, ErrorMatches, "(?s).*signedTx.*validation.required.*")
	c.Check(stageDoc, IsNil)
}

func (s *TradeIntegrationSuite) TestDownloadDocAuth(c *C) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/anything/", nil)
	ctx := routing.NewContext(res, req, s.noopDocHandler.HandleGetDocByID)
	err := ctx.Next()
	c.Assert(err, ErrorContains, "Authentication required")
}

func (s *TradeIntegrationSuite) TestCloseStageReqReject(c *C) {
	mr := s.noopResolver.Mutation()
	docHash := "f308fc02ce9172ad02a7d75800ecfc027109bc67987ea32aba9b8dcc7b10150e"

	// New doc
	docPath := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     0,
		StageDocIdx:  0,
		StageDocHash: docHash,
	}
	expireTimeExpected, err := time.Parse(time.RFC3339, "2096-01-02T15:04:05+00:00")
	c.Assert(err, IsNil)
	newDocTx, err := mr.MkTradeStageDocTx(s.seller.Ctx, docPath, model.ApprovalPending, &expireTimeExpected)
	c.Assert(err, IsNil)
	signedNewDocTx, err := testutil.SignTx(*s.noopDriver, newDocTx, testutil.SampleUser2Seed)
	c.Assert(err, IsNil)
	stageDoc, err := UploadDoc(
		s.seller.Ctx,
		s.noopDocHandler,
		UploadDocInput{
			StageIdx:     0,
			Data:         "test",
			TradeID:      s.trade.ID,
			ExpiresAt:    "2096-01-02T15:04:05+00:00",
			SignedTX:     signedNewDocTx,
			DocHash:      docHash,
			WithApproval: true,
		})
	c.Assert(err, IsNil)
	c.Assert(stageDoc, NotNil)
	expected := model.TradeStageDoc{
		Status:       "pending",
		ApprovedTx:   "",
		ApprovedBy:   "",
		ApprovedAt:   (*time.Time)(nil),
		ExpiresAt:    expireTimeExpected.UTC(),
		RejectReason: "",
	}
	expected.ReqTx = stageDoc.ReqTx
	expected.DocID = stageDoc.DocID
	c.Check(
		*stageDoc,
		DeepEquals,
		expected,
	)

	// Approve doc
	docApprovalTx, err := mr.MkTradeStageDocTx(s.buyer.Ctx, docPath, model.ApprovalApproved, &s.sampleExpireTime)
	c.Assert(err, IsNil)
	signedDocApproveTx, err := testutil.SignTx(*s.noopDriver, docApprovalTx, testutil.SampleUser1Seed)
	c.Assert(err, IsNil)
	c.Assert(signedDocApproveTx, NotNil)
	tradeStageDoc, err := testutil.ApproveDoc(s.buyer.Ctx, s.noopResolver, docPath, signedDocApproveTx)
	c.Assert(err, IsNil)
	c.Assert(tradeStageDoc, NotNil)

	// Try to approve stage (user2 is close initiator)
	stagePath := model.TradeStagePath{
		Tid:      s.trade.ID,
		StageIdx: 0,
	}
	stageCloseReqTx, err := mr.MkTradeStageCloseTx(s.seller.Ctx, stagePath, model.ApprovalPending)
	c.Assert(err, IsNil)
	signedStageCloseReqTx, err := testutil.SignTx(*s.noopDriver, stageCloseReqTx, testutil.SampleUser2Seed)
	c.Assert(err, IsNil)
	closeStageReq, err := testutil.CloseStage(s.seller.Ctx, s.noopResolver, stagePath, signedStageCloseReqTx, "My long reason")
	c.Assert(err, IsNil)
	c.Assert(closeStageReq, NotNil)
	// Let's see trade changes
	trade, err := testutil.GetTrade(s.buyer.Ctx, s.noopResolver, s.trade.ID)
	c.Assert(err, IsNil)
	c.Assert(trade, NotNil)
	c.Check(len(trade.Stages[0].CloseReqs), Equals, 1)
	c.Check(trade.Stages[0].CloseReqs[0].Status, Equals, model.ApprovalPending)
	c.Check(trade.Stages[0].CloseReqs[0].ReqTx, Matches, "noop-driver-[0-9]+")
	c.Check(trade.Stages[0].CloseReqs[0].ApprovedTx, Equals, "")
	// Reject stage approval
	stageCloseRejectedTx, err := mr.MkTradeStageCloseTx(s.buyer.Ctx, stagePath, model.ApprovalRejected)
	c.Assert(err, IsNil)
	signedStageCloseRejectedTx, err := testutil.SignTx(*s.noopDriver, stageCloseRejectedTx, testutil.SampleUser1Seed)
	c.Assert(err, IsNil)
	closeStageInteger, err := testutil.CloseStageReqReject(s.buyer.Ctx, s.noopResolver, stagePath, signedStageCloseRejectedTx, "My long reason")
	c.Assert(err, IsNil)
	c.Assert(closeStageReq, NotNil)
	// Magic number, not intuitive
	c.Assert(closeStageInteger, IsNil)

	// Let's see trade changes with rejected stage close request
	trade, err = testutil.GetTrade(s.buyer.Ctx, s.noopResolver, trade.ID)
	c.Assert(err, IsNil)
	c.Assert(trade, NotNil)
	c.Check(len(trade.Stages[0].CloseReqs), Equals, 1)
	c.Check(trade.Stages[0].CloseReqs[0].Status, Equals, model.ApprovalRejected)
	c.Check(trade.Stages[0].CloseReqs[0].ReqTx, Matches, "noop-driver-[0-9]+")
	c.Check(trade.Stages[0].CloseReqs[0].ApprovedTx, Matches, "noop-driver-[0-9]+")
	c.Check(trade.Stages[0].CloseReqs[0].ApprovedTx, Not(Equals), trade.Stages[0].CloseReqs[0].ReqTx)
}

func (s *TradeIntegrationSuite) TestCloseStageReqRejectNegative(c *C) {
	mr := s.noopResolver.Mutation()
	docHash := "f308fc02ce9172ad02a7d75800ecfc027109bc67987ea32aba9b8dcc7b10150e"

	// New doc
	docPath := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     0,
		StageDocIdx:  0,
		StageDocHash: docHash,
	}
	newDocTx, err := mr.MkTradeStageDocTx(s.seller.Ctx, docPath, model.ApprovalPending, &s.sampleExpireTime)
	c.Assert(err, IsNil)
	signedNewDocTx, err := testutil.SignTx(*s.noopDriver, newDocTx, testutil.SampleUser2Seed)
	c.Assert(err, IsNil)

	stageDoc, err := UploadDoc(
		s.seller.Ctx,
		s.noopDocHandler,
		UploadDocInput{
			StageIdx:     0,
			Data:         "test",
			TradeID:      s.trade.ID,
			ExpiresAt:    s.sampleExpireTimeStr,
			SignedTX:     signedNewDocTx,
			DocHash:      docHash,
			WithApproval: true,
		})
	c.Assert(err, IsNil)
	c.Assert(stageDoc, NotNil)
	c.Check(stageDoc.DocID, NotNil)
	c.Check(stageDoc.Status, Equals, model.ApprovalPending)
	c.Check(stageDoc.ReqTx, NotNil)
	c.Check(stageDoc.ApprovedTx, Equals, "")
	c.Check(stageDoc.ApprovedBy, Equals, "")
	c.Check(stageDoc.ApprovedAt, IsNil)
	c.Check(stageDoc.ExpiresAt, NotNil)
	c.Check(stageDoc.RejectReason, Equals, "")
	// Approve doc
	docApprovalTx, err := mr.MkTradeStageDocTx(s.buyer.Ctx, docPath, model.ApprovalApproved, &s.sampleExpireTime)
	c.Assert(err, IsNil)
	signedDocApproveTx, err := testutil.SignTx(*s.noopDriver, docApprovalTx, testutil.SampleUser1Seed)
	c.Assert(err, IsNil)
	c.Assert(signedDocApproveTx, NotNil)
	tradeStageDoc, err := testutil.ApproveDoc(s.buyer.Ctx, s.noopResolver, docPath, signedDocApproveTx)
	c.Assert(err, IsNil)
	c.Assert(tradeStageDoc, NotNil)

	// Try to approve stage (user2 is close initiator)
	stagePath := model.TradeStagePath{
		Tid:      s.trade.ID,
		StageIdx: 0,
	}
	stageCloseReqTx, err := mr.MkTradeStageCloseTx(s.seller.Ctx, stagePath, model.ApprovalPending)
	c.Assert(err, IsNil)
	signedStageCloseReqTx, err := testutil.SignTx(*s.noopDriver, stageCloseReqTx, testutil.SampleUser2Seed)
	c.Assert(err, IsNil)
	closeStageReq, err := testutil.CloseStage(s.seller.Ctx, s.noopResolver, stagePath, signedStageCloseReqTx, "My long reason")
	c.Assert(err, IsNil)
	c.Assert(closeStageReq, NotNil)
	// Let's see trade changes
	s.trade, err = testutil.GetTrade(s.buyer.Ctx, s.noopResolver, s.trade.ID)
	c.Assert(err, IsNil)
	c.Assert(s.trade, NotNil)
	c.Check(len(s.trade.Stages[0].CloseReqs), Equals, 1)
	c.Check(s.trade.Stages[0].CloseReqs[0].Status, Equals, model.ApprovalPending)
	c.Check(s.trade.Stages[0].CloseReqs[0].ReqTx, Matches, "noop-driver-[0-9]+")
	c.Check(s.trade.Stages[0].CloseReqs[0].ApprovedTx, Equals, "")
	// Reject stage approval
	stageCloseRejectedTx, err := mr.MkTradeStageCloseTx(s.buyer.Ctx, stagePath, model.ApprovalRejected)
	c.Assert(err, IsNil)
	signedStageCloseRejectedTx, err := testutil.SignTx(*s.noopDriver, stageCloseRejectedTx, testutil.SampleUser1Seed)
	c.Assert(err, IsNil)
	closeStageInteger, err := testutil.CloseStageReqReject(s.buyer.Ctx, s.noopResolver, stagePath, signedStageCloseRejectedTx, "My long reason")
	c.Assert(err, IsNil)
	c.Assert(closeStageReq, NotNil)
	// Magic number, not intuitive
	c.Assert(closeStageInteger, IsNil)

	// Let's see trade changes with rejected stage close request
	s.trade, err = testutil.GetTrade(s.seller.Ctx, s.noopResolver, s.trade.ID)
	c.Assert(err, IsNil)
	c.Assert(s.trade, NotNil)
	c.Check(len(s.trade.Stages[0].CloseReqs), Equals, 1)
	c.Check(s.trade.Stages[0].CloseReqs[0].Status, Equals, model.ApprovalRejected)
	c.Check(s.trade.Stages[0].CloseReqs[0].ReqTx, Matches, "noop-driver-[0-9]+")
	c.Check(s.trade.Stages[0].CloseReqs[0].ApprovedTx, Matches, "noop-driver-[0-9]+")
	c.Check(s.trade.Stages[0].CloseReqs[0].ApprovedTx, Not(Equals), s.trade.Stages[0].CloseReqs[0].ReqTx)
}

func (s *TradeIntegrationSuite) TestMkTxUnknownUser(c *C) {
	mr := s.noopResolver.Mutation()
	tx, err := mr.MkTradeCloseTx(s.third.Ctx, s.trade.ID, model.ApprovalApproved)
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "(?si).*permissions.trade.modify.denied.*")
	c.Assert(tx, Equals, "")

	stagePath := model.TradeStagePath{
		Tid:      s.trade.ID,
		StageIdx: 0,
	}
	tx, err = mr.MkTradeStageCloseTx(s.third.Ctx, stagePath, model.ApprovalApproved)
	c.Assert(err, NotNil)
	c.Assert(tx, Equals, "")
	c.Check(err, ErrorMatches, "(?si).*permissions.trade.modify.denied.*")

	docPath := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     0,
		StageDocIdx:  123,
		StageDocHash: "213",
	}
	tx, err = mr.MkTradeStageDocTx(s.third.Ctx, docPath, model.ApprovalApproved, &s.sampleExpireTime)
	c.Assert(err, NotNil)
	c.Assert(tx, Equals, "")
	c.Check(err, ErrorMatches, "(?si).*permissions.trade.modify.denied.*")

	newStagePath := model.TradeStagePath{
		Tid:      s.trade.ID,
		StageIdx: uint(len(s.trade.Stages)),
	}
	tx, err = mr.MkTradeStageAddTx(s.third.Ctx, newStagePath, model.ApprovalApproved)
	c.Assert(err, NotNil)
	c.Assert(tx, Equals, "")
	c.Check(err, ErrorMatches, "(?si).*permissions.trade.modify.denied.*")
}

func (s *TradeIntegrationSuite) TestUploadDocBothParties(c *C) {
	mr := s.noopResolver.Mutation()
	docHash := "f308fc02ce9172ad02a7d75800ecfc027109bc67987ea32aba9b8dcc7b10150e"
	docPath := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     1,
		StageDocIdx:  0,
		StageDocHash: docHash,
	}
	// the second stage's owner is seller
	// new doc from buyer should not be accepted
	newDocTx, err := mr.MkTradeStageDocTx(s.buyer.Ctx, docPath, model.ApprovalPending, &s.sampleExpireTime)
	c.Assert(err, IsNil)
	signedNewDocTx, err := testutil.SignTx(*s.noopDriver, newDocTx, testutil.SampleUser1Seed)
	c.Assert(err, IsNil)
	_, err = UploadDoc(
		s.buyer.Ctx,
		s.noopDocHandler,
		UploadDocInput{
			StageIdx:     1,
			Data:         "test",
			TradeID:      s.trade.ID,
			ExpiresAt:    "2096-01-02T15:04:05+07:00",
			SignedTX:     signedNewDocTx,
			DocHash:      docHash,
			WithApproval: true,
		})
	c.Assert(err, ErrorMatches, "(?si).*permissions.trade.modify.denied.*")

	// New doc from seller should be accepted
	docPath2 := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     1,
		StageDocIdx:  0,
		StageDocHash: docHash,
	}
	newSellerDocTx, err := mr.MkTradeStageDocTx(s.seller.Ctx, docPath2, model.ApprovalPending, &s.sampleExpireTime)
	c.Assert(err, IsNil)
	signedSellerDocTx, err := testutil.SignTx(*s.noopDriver, newSellerDocTx, testutil.SampleUser2Seed)
	c.Assert(err, IsNil)
	sellerDoc, err := UploadDoc(
		s.seller.Ctx,
		s.noopDocHandler,
		UploadDocInput{
			StageIdx:     1,
			Data:         "test",
			TradeID:      s.trade.ID,
			ExpiresAt:    s.sampleExpireTimeStr,
			SignedTX:     signedSellerDocTx,
			DocHash:      docHash,
			WithApproval: true,
		})
	c.Assert(err, IsNil)
	c.Assert(sellerDoc.DocID, Not(Equals), "")
}

func (s *TradeIntegrationSuite) TestUploadDocByOther(c *C) {
	mr := s.noopResolver.Mutation()
	docHash := "f308fc02ce9172ad02a7d75800ecfc027109bc67987ea32aba9b8dcc7b10150e"
	docPath := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     0,
		StageDocIdx:  0,
		StageDocHash: docHash,
	}
	// New doc from anyone who is not trade owner should be denied
	newDocTx, err := mr.MkTradeStageDocTx(s.buyer.Ctx, docPath, model.ApprovalPending, &s.sampleExpireTime)
	c.Assert(err, IsNil)
	signedNewDocTx, err := testutil.SignTx(*s.noopDriver, newDocTx, testutil.SampleUser1Seed)
	c.Assert(err, IsNil)
	_, err = UploadDoc(
		s.third.Ctx,
		s.noopDocHandler,
		UploadDocInput{
			StageIdx:     0,
			Data:         "test",
			TradeID:      s.trade.ID,
			ExpiresAt:    s.sampleExpireTimeStr,
			SignedTX:     signedNewDocTx,
			DocHash:      docHash,
			WithApproval: true,
		})
	c.Assert(err, ErrorMatches, "(?si).*permissions.trade.modify.denied.*")
}

func (s *TradeIntegrationSuite) TestSCLockConcurrent(c *C) {
	mr := s.noopResolver.Mutation()
	tx1, err := mr.MkTradeCloseTx(s.buyer.Ctx, s.trade.ID, model.ApprovalApproved)
	c.Assert(err, IsNil)
	s1, _, err := txvalidation.Simplify(tx1)
	c.Assert(err, IsNil)
	tx2, err := mr.MkTradeCloseTx(s.seller.Ctx, s.trade.ID, model.ApprovalApproved)
	s2, _, err := txvalidation.Simplify(tx2)
	c.Assert(err, IsNil)
	c.Check(s1.SourceAccount, Not(Equals), s2.SourceAccount, Comment("Two concurrent txs should originate from different source accounts"))
}

func (s *TradeIntegrationSuite) TestSCLockSameUserTwice(c *C) {
	mr := s.noopResolver.Mutation()
	_, err := mr.MkTradeCloseTx(s.buyer.Ctx, s.trade.ID, model.ApprovalApproved)
	c.Assert(err, IsNil)
	_, err = mr.MkTradeCloseTx(s.buyer.Ctx, s.trade.ID, model.ApprovalApproved)
	c.Assert(err, IsNil)
}

func (s *TradeIntegrationSuite) TestStageDocWithConfirmation(c *C) {
	mr := s.noopResolver.Mutation()
	s.trade.Stages = append([]model.TradeStage{}, model.TradeStage{
		Name:        "testStage",
		Description: sampleDesc,
		Owner:       model.TradeActorB,
		Docs: []model.TradeStageDoc{model.TradeStageDoc{
			DocID:     "aaa1",
			Status:    model.ApprovalPending,
			ExpiresAt: time.Now().UTC().Add(time.Hour),
			ReqTx:     "rawTxSigned",
		}},
	})
	tradeStageDocPath := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     0,
		StageDocHash: testDocHash,
		StageDocIdx:  0,
	}
	docPath := model.TradeStageDocPath{
		Tid:          s.trade.ID,
		StageIdx:     0,
		StageDocIdx:  0,
		StageDocHash: testDocHash,
	}
	var err error
	_, err = dal.UpdateTrade(s.buyer.Ctx, s.db, s.trade)
	c.Check(err, IsNil)

	// test for trade stage doc approving
	rawTx, err := mr.MkTradeStageDocTx(s.seller.Ctx, docPath, model.ApprovalApproved, &s.sampleExpireTime)
	c.Assert(err, IsNil)
	rawTxSigned, err := testutil.SignTx(
		*s.noopDriver,
		rawTx,
		testutil.SampleUser2Seed,
	)
	c.Assert(err, IsNil)
	stageDoc, err := mr.TradeStageDocApprove(s.seller.Ctx, tradeStageDocPath, rawTxSigned)
	c.Check(err, IsNil)
	c.Check(stageDoc.Status, Equals, model.ApprovalApproved)
	c.Check(stageDoc.ApprovedBy, Equals, s.seller.ID)
	c.Check(stageDoc.DocID, Equals, s.trade.Stages[0].Docs[0].DocID)

	notifications, err := s.noopResolver.Query().NotificationsTrade(s.buyer.Ctx, s.trade.ID)
	c.Assert(err, IsNil)
	c.Check(notifications[0].Receiver, Contains, s.buyer.ID)
	c.Check(notifications[0].EntityID, Equals, bat.StrJoin("/", s.trade.FullID2(), "stages:0", "docs:0"))
	c.Check(notifications[0].Action, Equals, model.ApprovalApproved)
}
