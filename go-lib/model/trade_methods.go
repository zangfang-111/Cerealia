package model

import (
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	"bitbucket.org/cerealia/apps/go-lib/validation"
	"github.com/robert-zaremba/errstack"
)

const buyerIsSeller = "validation.buyer-is-seller"
const creatorIsNotSellerOrBuyer = "validation.creator-is-not-seller-or-buyer"

// Field names copied from models_gql_gen.go
const nameField = "name"
const sellerIDField = "sellerID"
const buyerIDField = "buyerID"

type tradeValidatorBuilder struct {
	validation.Builder
}

func (vb *tradeValidatorBuilder) sellerShouldNotBeBuyer(buyerID, sellerID string) {
	if buyerID == sellerID {
		vb.Append(sellerIDField, buyerIsSeller)
		vb.Append(buyerIDField, buyerIsSeller)
	}
}

func (vb *tradeValidatorBuilder) creatorIsSellerOrBuyer(buyerID, sellerID, submitterID string) {
	if buyerID != submitterID && sellerID != submitterID {
		vb.Append(sellerIDField, creatorIsNotSellerOrBuyer)
		vb.Append(buyerIDField, creatorIsNotSellerOrBuyer)
	}
}

// Validate checks the create trade request with its all parameters.
func (input NewTradeInput) Validate(templateID string, u *User) errstack.Builder {
	vb := tradeValidatorBuilder{}
	vb.MinLength(nameField, input.Name, 5)
	vb.Required(sellerIDField, input.SellerID)
	vb.Required(buyerIDField, input.BuyerID)
	vb.creatorIsSellerOrBuyer(input.BuyerID, input.SellerID, u.ID)
	vb.sellerShouldNotBeBuyer(input.BuyerID, input.SellerID)
	return vb.ToErrstackBuilder()
}

// SetID implements dal.HasID interface
func (t *Trade) SetID(id string) {
	t.ID = id
}

// GetStage returns stage if it exists
func (t Trade) GetStage(idx uint) (*TradeStage, errstack.E) {

	if int(idx) >= len(t.Stages) {
		return nil, errstack.NewReq("Stage Index out of range")
	}
	return &t.Stages[idx], nil
}

// GetStageAddReq returns stage add request if it exists
func (t Trade) GetStageAddReq(idx uint) (*TradeStageAddReq, errstack.E) {
	if int(idx) >= len(t.StageAddReqs) {
		return nil, errstack.NewReq("Stage AddReq index out of range")
	}
	return &t.StageAddReqs[idx], nil
}

// GetStageDoc returns stage and StageDoc if it exists
func (t Trade) GetStageDoc(stageIdx uint, docIdx uint) (*TradeStage, *TradeStageDoc, errstack.E) {
	s, err := t.GetStage(stageIdx)
	if err != nil {
		return nil, nil, err
	}
	d, err := s.GetDoc(docIdx)
	return s, d, err
}

// Requester checks the requester info for auth and returns as tradeActor
func (t Trade) Requester(u *User) (TradeActor, errstack.E) {
	if u == nil || u.ID == "" {
		return TradeActorB, ErrUnauthenticated
	}
	if t.Buyer.UserID == u.ID {
		return TradeActorB, nil
	} else if t.Seller.UserID == u.ID {
		return TradeActorS, nil
	} else if u.IsModerator() {
		return TradeActorM, nil
	}
	return TradeActorB, ErrUnauthorized
}

// NotificationReceiver checks the requester info for notification action and returns receiver as string(user id) array
func (t Trade) NotificationReceiver(u *User) ([]string, errstack.E) {
	if t.Buyer.UserID == u.ID {
		return []string{t.Seller.UserID}, nil
	} else if t.Seller.UserID == u.ID {
		return []string{t.Buyer.UserID}, nil
	} else if u.IsModerator() {
		return []string{t.Buyer.UserID, t.Seller.UserID}, nil
	}
	return nil, ErrUnauthorized
}

// CanBeModifiedBy checks whether a trade can be modified by the given user
func (t *Trade) CanBeModifiedBy(user *User) errstack.E {
	if user.ID != t.Buyer.UserID && user.ID != t.Seller.UserID && !user.IsModerator() {
		return ErrUnauthorized
	}
	return nil
}

// CanMakeCloseReq returns the possibility if a trade close request can be made.
func (t *Trade) CanMakeCloseReq() bool {
	if len(t.CloseReqs) > 0 && t.CloseReqs[len(t.CloseReqs)-1].Status != ApprovalRejected {
		return false
	}
	for _, s := range t.Stages {
		if !(len(s.DelReqs) > 0 && s.DelReqs[len(s.DelReqs)-1].Status == ApprovalApproved) &&
			!(len(s.CloseReqs) > 0 && s.CloseReqs[len(s.CloseReqs)-1].Status == ApprovalApproved) {
			return false
		}
	}
	return true
}

// CheckTradeClosed checks if the trade has already been closed or not
func (t *Trade) CheckTradeClosed() bool {
	return len(t.CloseReqs) > 0 && t.CloseReqs[len(t.CloseReqs)-1].Status == ApprovalApproved
}

// FullID returns full ID as in ArangoDB
func (t Trade) FullID() string {
	return dbconst.ColTrades.FullID(t.ID)
}

// FullID2 returns splited string from trade's FullID
func (t Trade) FullID2() string {
	return string(dbconst.ColTrades) + ":" + t.ID
}

// CanBeApproved check and returns error if a trade stage add request
// can't be approved or rejected
func (sr TradeStageAddReq) CanBeApproved() errstack.E {
	if sr.ApprovedBy == sr.ReqBy {
		return errstack.NewReq("You have made this request, so you can't approve this.")
	}
	return sr.Status.ShouldBePending()
}

// CanMakeDelReq returns the possibility if a stage delete request can be made.
func (s *TradeStage) CanMakeDelReq() bool {
	if len(s.Docs) > 0 {
		return false
	}
	l := len(s.DelReqs)
	return l == 0 || s.DelReqs[l-1].Status != ApprovalPending // last delreq is not pending, then can delete
}

// GetLastDeletionRequest checks and returns error if a trade stage del request can not
// be approved and executed.
func (s TradeStage) GetLastDeletionRequest() (*ApproveReq, errstack.E) {
	l := len(s.DelReqs)
	if l == 0 {
		return nil, errstack.NewReq("There is no delete request to approve")
	}
	lastReq := &s.DelReqs[l-1]
	return lastReq, nil
}

// CanMakeCloseReq returns the possibility if a stage close request can be made.
func (s *TradeStage) CanMakeCloseReq() bool {
	for _, doc := range s.Docs {
		if doc.Status != ApprovalRejected {
			return true
		}
	}
	return false
}

// GetLastClosingRequest checks and returns error if a trade stage close request can not
// be approved and executed.
func (s TradeStage) GetLastClosingRequest() (*ApproveReq, errstack.E) {
	l := len(s.CloseReqs)
	if l == 0 {
		return nil, errstack.NewReq("No close request in this stage")
	}
	lastReq := &s.CloseReqs[l-1]
	return lastReq, nil
}

// GetDoc returns stage doc if it exists
func (s TradeStage) GetDoc(idx uint) (*TradeStageDoc, errstack.E) {
	if int(idx) >= len(s.Docs) {
		return nil, errstack.NewReq("Stage Doc index out of range")
	}
	return &s.Docs[idx], nil
}

// AssertOwnedBy checks whether a user is the owner of the stage or not, returns permission error if not owner
func (s TradeStage) AssertOwnedBy(t *Trade, uid string) errstack.E {
	if (s.Owner == TradeActorN && t.Buyer.UserID != uid && t.Seller.UserID != uid) ||
		(s.Owner == TradeActorB && t.Buyer.UserID != uid) ||
		(s.Owner == TradeActorS && t.Seller.UserID != uid) {
		return errstack.NewReq("trade.stage edit.not-authorized")
	}
	return nil
}

// ShouldBePending checks if status has expected value
func (s Approval) ShouldBePending() errstack.E {
	if s == ApprovalPending {
		return nil
	}
	return errstack.NewReq("You can only approve/reject objects with pending status. Currently status=" + s.String())
}

// IsDeletedOrClosed checks if the stage has been deleted or closed
func (s *TradeStage) IsDeletedOrClosed() bool {
	delReq, errs := s.GetLastDeletionRequest()
	if errs == nil && delReq.Status == ApprovalApproved {
		return true
	}
	closeReq, errs := s.GetLastClosingRequest()
	return errs == nil && closeReq.Status == ApprovalApproved
}

// NilOrPending check if the obj is pending status
func (obj ApproveReq) NilOrPending() bool {
	return obj.ApprovedBy == "" || obj.ApprovedBy == ApprovalNil.String() || obj.ApprovedBy == ApprovalPending.String()
}

// SetApprovedBy mutates the request and sets the status to "approved"
func (obj *ApproveReq) SetApprovedBy(userID string) {
	now := time.Now().UTC()
	obj.ApprovedBy = userID
	obj.ApprovedAt = &now
	obj.Status = ApprovalApproved
}

// NilOrPending check if the obj is pending status
func (obj TradeStageDoc) NilOrPending() bool {
	return obj.ApprovedBy == "" || obj.ApprovedBy == ApprovalNil.String() || obj.ApprovedBy == ApprovalPending.String()
}

// SetID implements dal.HasID interface
func (tt *TradeTemplate) SetID(id string) {
	tt.ID = id
}

// FindParticipant looks up for the user in the trade
func (t *Trade) FindParticipant(user *User) (*TradeParticipant, errstack.E) {
	if t.Buyer.UserID == user.ID {
		return &t.Buyer, nil
	}
	if t.Seller.UserID == user.ID {
		return &t.Seller, nil
	}
	if user.IsModerator() {
		// HACK! HD wallets won't work
		// TODO: Remove it during HD wallet refactoring
		sw, ok := user.StaticWallets[user.DefaultWalletID]
		if !ok {
			// Not telling that it's a moderator user
			return nil, errstack.NewReqF("User '%s' doesn't have a default wallet", user.ID)
		}
		return &TradeParticipant{
			UserID: user.ID,
			//KeyDerivationPath: "under construction",
			//WalletID:          "under construction",
			PubKey: sw.PubKey,
		}, nil
	}
	return nil, errstack.NewReqF("User '%s' is not a participant of trade '%s'", user.ID, t.ID)
}

// FindPubKey finds user's key for a trade
func (t Trade) FindPubKey(user *User) (pubKey SCAddr, error errstack.E) {
	tp, err := t.FindParticipant(user)
	if err != nil {
		return "", errstack.WrapAsInfF(err, "Key not found")
	}
	return SCAddr(tp.PubKey), nil
}
