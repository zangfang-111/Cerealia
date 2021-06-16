// Package users contains rest api services for users
package users

import (
	"path/filepath"

	"bitbucket.org/cerealia/apps/go-lib/model/dal"

	"bitbucket.org/cerealia/apps/cmd/websrv/config"
	"bitbucket.org/cerealia/apps/go-lib/fstore"
	"bitbucket.org/cerealia/apps/go-lib/middleware"
	dbs "bitbucket.org/cerealia/apps/go-lib/setup/arangodb"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/log15"
)

const maxAvatarSize int64 = 5000000 // 5 MB
const avatarDir = "public/user/avatars"

var logger = log15.Root()

// UserHandler is a route object for users
type UserHandler struct{}

// HandlePostUserAvatar is to upload user avatar file to server
func (h UserHandler) HandlePostUserAvatar(c *routing.Context) error {
	ctx := c.Request.Context()
	u, errs := middleware.GetAuthUser(ctx)
	if errs != nil {
		return errs
	}
	err := c.Request.ParseMultipartForm(maxAvatarSize)
	if err != nil {
		return errstack.WrapAsReq(err, "Can't Parse the form data")
	}
	db, errs := dbs.GetDb(ctx)
	if errs != nil {
		return errs
	}
	fileHeader := c.Request.MultipartForm.File["avatarFile"][0]
	file, err := fileHeader.Open()
	defer errstack.CallAndLog(logger, file.Close)
	if err != nil {
		return errstack.WrapAsReq(err, "Can't get file data from request")
	}
	avatarStoreDir := filepath.Join(config.F.FileStorageDir.String(), avatarDir)
	fileName, errs := fstore.SaveAvatar(file, fileHeader.Filename, avatarStoreDir)
	if errs != nil {
		return errs
	}
	return dal.UpdateAvatar(ctx, db, u.ID, fileName)
}

// SetUserRoutes sets user routes
func SetUserRoutes(routerG *routing.RouteGroup) {
	u := UserHandler{}
	routerG.Post("/avatar", u.HandlePostUserAvatar)
}
