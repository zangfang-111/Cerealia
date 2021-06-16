package trades

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	"bitbucket.org/cerealia/apps/go-lib/model"
	routing "github.com/go-ozzo/ozzo-routing"
)

// UploadDocInput is an object that wraps parameters for document upload
type UploadDocInput struct {
	StageIdx     uint
	Data         string
	TradeID      string
	ExpiresAt    string
	SignedTX     string
	DocHash      string
	WithApproval bool
}

func UploadDoc(userCtx context.Context, docHandler DocHandler, input UploadDocInput) (*model.TradeStageDoc, error) {
	yrl, err := url.Parse("/v1/trades/stage-docs")
	if err != nil {
		return nil, err
	}
	parameters := url.Values{}
	parameters.Add("tid", input.TradeID)
	parameters.Add("stageIdx", fmt.Sprint(input.StageIdx))
	parameters.Add("expiresAt", input.ExpiresAt)
	parameters.Add("signedTx", input.SignedTX)
	parameters.Add("hash", input.DocHash)
	parameters.Add("withApproval", strconv.FormatBool(input.WithApproval))
	yrl.RawQuery = parameters.Encode()
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		yrl.String(),
		bytes.NewReader([]byte(("--11111\nContent-Disposition: form-data; name=\"formfile\"; filename=\"aaaa-test.txt.pdf\"\nContent-Type: application/pdf\n\n" + input.Data + "\n--11111--\n"))))
	req.Header.Add("Content-Type", `multipart/form-data; boundary=11111`)
	reqUCtx := req.WithContext(userCtx)
	ctx := routing.NewContext(res, reqUCtx, docHandler.HandlePostTradeStageDoc)
	err = ctx.Next()
	if err != nil {
		return nil, err
	}
	doc := model.TradeStageDoc{}
	return &doc, json.Unmarshal(res.Body.Bytes(), &doc)
}
