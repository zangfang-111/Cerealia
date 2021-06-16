package dal

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"

	"bitbucket.org/cerealia/apps/go-lib/model"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

// GetNotification gets a notification object by its id
func GetNotification(ctx context.Context, db driver.Database, id string) (*model.Notification, errstack.E) {
	var n model.Notification
	return &n, DBGetOneFromColl(ctx, &n, id, dbconst.ColNotifications, db)
}

// InsertNotification inserts new notification
func InsertNotification(ctx context.Context, db driver.Database, n *model.Notification) errstack.E {
	_, errs := insertHasID(ctx, dbconst.ColNotifications, n, db)
	return errs
}

// GetUserNotifacations returns 30 active notifications starting at @from position.
func GetUserNotifacations(ctx context.Context, db driver.Database, uid string, from uint) ([]model.Notification, errstack.E) {
	var ns []model.Notification
	query := `for d in notifications
	 filter @uid in d.receiver && @uid not in d.dismissed
	 sort d.createdAt desc limit @start, 30 return d`
	bindVars := map[string]interface{}{
		"start": from * 30,
		"uid":   uid,
	}
	err := DBQueryMany(ctx, &ns, query, bindVars, db)
	return ns, err
}

// GetNotifacationsByTrade gets all notifications for specific trade
func GetNotifacationsByTrade(ctx context.Context, db driver.Database, uid, tid string) ([]model.Notification, errstack.E) {
	var ns []model.Notification
	query := `for d in notifications
	 filter d.entityID like @tid && @uid in d.receiver && @uid not in d.dismissed
	 sort d.createdAt desc return d`
	bindVars := map[string]interface{}{
		"uid": uid,
		"tid": "trades:" + tid + "/%",
	}
	err := DBQueryMany(ctx, &ns, query, bindVars, db)
	return ns, err
}

// NotificationDismiss dismisses the notification
func NotificationDismiss(ctx context.Context, db driver.Database, uid, notifID string) errstack.E {
	query := `for d in notifications update { _key:@notifID, dismissed: append(d.dismissed, @uid)} in notifications`
	bindVars := map[string]interface{}{
		"uid":     uid,
		"notifID": notifID,
	}
	_, err := db.Query(ctx, query, bindVars)
	return errstack.WrapAsInf(err, "Failed to dismiss notification in db")
}
