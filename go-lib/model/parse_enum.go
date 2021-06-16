package model

import (
	"errors"

	"github.com/robert-zaremba/errstack"
)

// ParseTradeActor converts string to TradeActor value
func ParseTradeActor(s string, errb errstack.Builder) TradeActor {
	switch TradeActor(s) {
	case TradeActorB, TradeActorS, TradeActorM:
		return TradeActor(s)
	}
	errb.Put("owner", "wrong tradeActor value")
	return TradeActorS
}

// ParseApproved converts string to Approved value
func ParseApproved(s string) (Approval, error) {
	switch Approval(s) {
	case ApprovalNil, ApprovalApproved, ApprovalPending, ApprovalRejected:
		return Approval(s), nil
	}
	return ApprovalNil, errors.New("wrong approval value")
}
