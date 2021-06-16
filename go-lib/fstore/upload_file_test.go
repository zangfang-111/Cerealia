package fstore

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/cerealia/apps/go-lib/setup"
	routing "github.com/go-ozzo/ozzo-routing"
	. "github.com/robert-zaremba/checkers"
	"github.com/robert-zaremba/errstack"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type LocalStorageSuite struct{}

var _ = Suite(&LocalStorageSuite{})

var (
	filesourcepath = "/test/testuploadfiles/"
	fileuploadpath = "/tmp/contractfiles-test/upload"
	filefieldname  = "formfile"
	maxfilesize    = 5000000
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer errstack.CallAndLog(logger, file.Close)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func (s *LocalStorageSuite) TestUploadDocument(c *C) {
	path := setup.RootDir + filesourcepath

	docpath := path + "docxFile.docx"
	filetypecheck(docpath, "docx", c, true)

	docpath = path + "otdFile.odt"
	filetypecheck(docpath, "odt", c, true)

	docpath = path + "pdfFile.pdf"
	filetypecheck(docpath, "pdf", c, true)

	docpath = path + "zipFile.zip"
	filetypecheck(docpath, "zip", c, true)

	docpath = path + "docFile.doc"
	filetypecheck(docpath, "doc", c, true)

	// Negative tests
	docpath = path + "rtfFile.rtf"
	filetypecheck(docpath, "rtf", c, false)

}

func filetypecheck(filename, typ string, c *C, isPositive bool) {
	// create new request context as for form file upload
	req, err := newfileUploadRequest("", nil, filefieldname, filename)
	c.Check(err, IsNil, Comment("The file info to upload is not correct. please make sure the correct file info"))
	context := routing.NewContext(nil, req, nil)

	// validate the file size as given limit size
	context.Request.Body = http.MaxBytesReader(context.Response, context.Request.Body, int64(maxfilesize))
	err = context.Request.ParseMultipartForm(int64(maxfilesize))
	c.Check(err, IsNil, Comment("Filesize is too big.", err))

	// validate the uploading of the files
	file, handler, err := context.Request.FormFile(filefieldname)
	c.Assert(err, IsNil)
	_, _, errs := SaveDoc(file, handler.Filename, fileuploadpath)
	if isPositive {
		c.Check(errs, IsNil, Comment("Filetype ", typ, " must be supported.", errs))
	} else {
		c.Check(errs, NotNil, Comment("Filetype ", typ, " uploaded unexpectedly. It should not be supported."))
	}
}
