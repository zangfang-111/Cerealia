package resolver

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/utils"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
	bat "github.com/robert-zaremba/go-bat"
)

func tadeCreateNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User) (*model.Notification, errstack.E) {
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = t.FullID2() + "/"
	n.Msg = fmt.Sprintf("New trade '%s' has been created by %s %s", t.Name, u.FirstName, u.LastName)
	n.Action = model.ApprovalApproved
	return n, dal.InsertNotification(ctx, db, n)
}

func tradeStageAddReqNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User, withApproval bool) (*model.Notification, errstack.E) {
	idx := strconv.Itoa(len(t.StageAddReqs) - 1)
	action := model.ApprovalSubmitted
	msg := fmt.Sprintf("New stage has been added by %s %s", u.FirstName, u.LastName)
	if withApproval {
		action = model.ApprovalPending
		msg = fmt.Sprintf("New stage add request has been created by %s %s", u.FirstName, u.LastName)
	}
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = bat.StrJoin("/", t.FullID2(), "stageAddReqs:"+idx)
	n.Msg = msg
	n.Action = action
	return n, dal.InsertNotification(ctx, db, n)
}

func tradeStageAddApprovalNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User,
	id model.TradeStagePath, isApprove bool) (*model.Notification, errstack.E) {
	msg := fmt.Sprintf("New stage add request has been approved by %s %s", u.FirstName, u.LastName)
	action := model.ApprovalApproved
	if !isApprove {
		msg = fmt.Sprintf("New stage add request has been rejected by %s %s", u.FirstName, u.LastName)
		action = model.ApprovalRejected
	}
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = bat.StrJoin("/", t.FullID2(), "stageAddReqs:"+utils.UintToString(id.StageIdx))
	n.Msg = msg
	n.Action = action
	return n, dal.InsertNotification(ctx, db, n)
}

func tradeStageDelReqNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User,
	id model.TradeStagePath) (*model.Notification, errstack.E) {
	idx := strconv.Itoa(len(t.Stages[id.StageIdx].DelReqs) - 1)
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = bat.StrJoin("/", t.FullID2(), "stages:"+utils.UintToString(id.StageIdx), "delReqs:"+idx)
	n.Msg = fmt.Sprintf("New stage delete request has been created by %s %s", u.FirstName, u.LastName)
	n.Action = model.ApprovalPending
	return n, dal.InsertNotification(ctx, db, n)
}

func tradeStageDeleteApprovalNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User,
	id model.TradeStagePath, isApprove bool) (*model.Notification, errstack.E) {
	idx := strconv.Itoa(len(t.Stages[id.StageIdx].DelReqs) - 1)
	msg := fmt.Sprintf("New stage delete request has been approved by %s %s", u.FirstName, u.LastName)
	action := model.ApprovalApproved
	if !isApprove {
		msg = fmt.Sprintf("New stage delete request has been rejected by %s %s", u.FirstName, u.LastName)
		action = model.ApprovalRejected
	}
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = bat.StrJoin("/", t.FullID2(), "stages:"+utils.UintToString(id.StageIdx), "delReqs:"+idx)
	n.Msg = msg
	n.Action = action
	return n, dal.InsertNotification(ctx, db, n)
}

func tradeStageDocApprovalNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User,
	id model.TradeStageDocPath, isApprove bool) (*model.Notification, errstack.E) {
	msg := fmt.Sprintf("New stage doc add request has been approved by %s %s", u.FirstName, u.LastName)
	action := model.ApprovalApproved
	if !isApprove {
		msg = fmt.Sprintf("New stage doc add request has been rejected by %s %s", u.FirstName, u.LastName)
		action = model.ApprovalRejected
	}
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = bat.StrJoin("/", t.FullID2(), "stages:"+utils.UintToString(id.StageIdx), "docs:"+utils.UintToString(id.StageDocIdx))
	n.Msg = msg
	n.Action = action
	return n, dal.InsertNotification(ctx, db, n)
}

func tradeStageCloseReqNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User,
	id model.TradeStagePath) (*model.Notification, errstack.E) {
	idx := strconv.Itoa(len(t.Stages[id.StageIdx].CloseReqs) - 1)
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = bat.StrJoin("/", t.FullID2(), "stages:"+utils.UintToString(id.StageIdx), "closeReqs:"+idx)
	n.Msg = fmt.Sprintf("New stage close request has been created by %s %s", u.FirstName, u.LastName)
	n.Action = model.ApprovalPending
	return n, dal.InsertNotification(ctx, db, n)
}

func tradeStageCloseApprovalNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User,
	id model.TradeStagePath, isApprove bool) (*model.Notification, errstack.E) {
	idx := strconv.Itoa(len(t.Stages[id.StageIdx].CloseReqs) - 1)
	msg := fmt.Sprintf("New stage close request has been approved by %s %s", u.FirstName, u.LastName)
	action := model.ApprovalApproved
	if !isApprove {
		msg = fmt.Sprintf("New stage close request has been rejected by %s %s", u.FirstName, u.LastName)
		action = model.ApprovalRejected
	}
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = bat.StrJoin("/", t.FullID2(), "stages:"+utils.UintToString(id.StageIdx), "closeReqs:"+idx)
	n.Msg = msg
	n.Action = action
	return n, dal.InsertNotification(ctx, db, n)
}

func tradeCloseReqNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User) (*model.Notification, errstack.E) {
	idx := strconv.Itoa(len(t.CloseReqs) - 1)
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = bat.StrJoin("/", t.FullID2(), "closeReqs:"+idx)
	n.Msg = fmt.Sprintf("New trade close request has been created by %s %s", u.FirstName, u.LastName)
	n.Action = model.ApprovalPending
	return n, dal.InsertNotification(ctx, db, n)
}

func tradeCloseApprovalNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User,
	isApprove bool) (*model.Notification, errstack.E) {
	idx := strconv.Itoa(len(t.CloseReqs) - 1)
	msg := fmt.Sprintf("New trade close request has been approved by %s %s", u.FirstName, u.LastName)
	action := model.ApprovalApproved
	if !isApprove {
		msg = fmt.Sprintf("New trade close request has been rejected by %s %s", u.FirstName, u.LastName)
		action = model.ApprovalRejected
	}
	n := mkBasicNotification(ctx, db, t, u)
	n.EntityID = bat.StrJoin("/", t.FullID2(), "closeReqs:"+idx)
	n.Msg = msg
	n.Action = action
	return n, dal.InsertNotification(ctx, db, n)
}

// func tradeStageSetExpireNotif(ctx context.Context, db driver.Database, t *model.Trade, u *model.User,
// 	id model.TradeStagePath) (*model.Notification, errstack.E) {
// 	n := mkBasicNotification(ctx, db, t, u)
// 	n.EntityID = bat.StrJoin("/", t.FullID2(), "stages:"+utils.UintToString(id.StageIdx), "expireTime:0")
// 	n.Msg = fmt.Sprintf("New stage expire time has been set by %s %s", u.FirstName, u.LastName)
// 	n.Action = model.ApprovalApproved
// 	return n, dal.InsertNotification(ctx, db, n)
// }

func mkBasicNotification(ctx context.Context, db driver.Database, t *model.Trade, u *model.User) *model.Notification {
	receiver, _ := t.NotificationReceiver(u)
	n := model.Notification{
		CreatedAt:   time.Now().UTC(),
		Receiver:    receiver,
		TriggeredBy: u.ID,
		Type:        model.NotifTypeAction,
		Dismissed:   []string{},
	}
	return &n
}
