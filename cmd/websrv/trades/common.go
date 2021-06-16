package trades

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"bitbucket.org/cerealia/apps/cmd/websrv/config"
	"bitbucket.org/cerealia/apps/go-lib/fstore"
	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	dbs "bitbucket.org/cerealia/apps/go-lib/setup/arangodb"
	driver "github.com/arangodb/go-driver"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/log15"
)

var logger = log15.Root()

const singleFile = 1
const maxDocSize int64 = 5000000 // 5 MB
const maxDocSizeMB = maxDocSize / 1000000
const (
	tradeDocDir      = "trade-docs"
	tradeOfferDocDir = "tradeoffer-docs"
)

func readInput(c *routing.Context, uploaderUserID string) (*TradeStageDocInputP, errstack.E) {
	err := c.Request.ParseMultipartForm(maxDocSize * 10)
	if err != nil {
		return nil, errstack.WrapAsReq(err, "Can't Parse the form data")
	}
	if len(c.Request.MultipartForm.File["formfile"]) != singleFile {
		return nil, errstack.NewReqF("Expecting %d file", singleFile)
	}
	var input = TradeStageDocInput{Uploader: uploaderUserID}
	err = c.Read(&input)
	if err != nil {
		return nil, errstack.WrapAsReq(err, "Can't read the input")
	}
	parsed, errb := input.parseFields()
	return parsed, errb.ToReqErr()
}

func storeDocFile(r *http.Request, idx int, destDir string) (model.FileInfo, errstack.E) {
	fileHeader := r.MultipartForm.File["formfile"][idx]
	if fileHeader.Size > maxDocSize {
		return model.FileInfo{}, errstack.NewReq(
			fmt.Sprint(fileHeader.Filename, " file is too big. Max size: ", maxDocSizeMB))
	}
	file, err := fileHeader.Open()
	if err != nil {
		return model.FileInfo{}, errstack.WrapAsReq(err, "Can't get file data from request")
	}
	defer errstack.CallAndLog(logger, file.Close)

	// get the query data
	dir := filepath.Join(config.F.FileStorageDir.String(), destDir)
	storedName, hash, errs := fstore.SaveDoc(file, fileHeader.Filename, dir)
	return model.FileInfo{
		FileName: fileHeader.Filename,
		Hash:     hash,
		URL:      storedName,
	}, errs
}

func validateStageOwnership(ctx context.Context, db driver.Database, u *model.User, tradeID string, stageIdx uint) (*model.Trade, *model.TradeStage, errstack.E) {
	t, err := dal.GetTrade(ctx, db, tradeID)
	if err != nil {
		return nil, nil, err
	}
	if int(stageIdx) >= len(t.Stages) {
		return nil, nil, errstack.NewReqF("Stage '%d' does not exist", stageIdx)
	}
	s := t.Stages[stageIdx]
	return t, &s, s.AssertOwnedBy(t, u.ID)
}

func readAndValidate(ctx context.Context, c *routing.Context, db driver.Database, u *model.User, now time.Time) (*TradeStageDocInputP, *model.Trade, *model.TradeStage, errstack.E) {
	input, err := readInput(c, u.ID)
	if err != nil {
		return input, nil, nil, errstack.WrapAsReqF(err, "Can't read document input")
	}
	t, s, err := validateStageOwnership(ctx, db, u, input.Tid, input.StageIdx)
	if now.After(input.ExpiresAtTime) && input.WithApproval {
		return input, t, s, errstack.NewReq("Approval expired")
	}
	return input, t, s, err
}

func getAndCheckAuthUser(c *routing.Context) (context.Context, driver.Database, *model.User, errstack.E) {
	ctx := c.Request.Context()
	u, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	db, err := dbs.GetDb(ctx)
	return ctx, db, u, err
}

func serveDocFile(ctx context.Context, c *routing.Context, db driver.Database, docID string) errstack.E {
	doc, err := dal.GetDoc(ctx, db, docID)
	if err != nil {
		return err
	}
	// Content-Type is set by header negotiation middleware (forced to application/json)
	// For correct web browser preview http.ServeFile has to choose the header
	c.Response.Header().Del("Content-Type")
	http.ServeFile(c.Response, c.Request,
		filepath.Join(config.F.FileStorageDir.String(), tradeDocDir, doc.URL))
	return nil
}
