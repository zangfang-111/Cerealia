// Package fstore contains file store functions for document upload
package fstore

import (
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pborman/uuid"
	"github.com/robert-zaremba/errstack"
	bat "github.com/robert-zaremba/go-bat"
	"github.com/robert-zaremba/log15"
	"golang.org/x/crypto/blake2s"
)

var logger = log15.Root()

var docTypes = []string{".pdf", ".doc", ".docx", ".odt", ".gzip", ".zip"}
var avatarTypes = []string{".jpg", ".jpeg", ".gif", ".png"}

// check file type, only pdf, doc, docx, odt, gzip, zip files are permitted
func getFileExtension(filename string, fileType []string) (string, error) {
	fx := strings.ToLower(filepath.Ext(filename))
	if fx == "" {
		return fx, errstack.NewReq(filename + " this file has no extension and is invalid for upload")
	}
	if bat.StrSliceIdx(fileType, fx) < 0 {
		return fx, errstack.NewReq(filename + " this file type doesn't match")
	}
	return fx, nil
}

// SaveDoc saves document file to disk and calculates it's hash
func SaveDoc(src io.Reader, filename string, storageDir string) (string, string, errstack.E) {
	copy, hasher := readerHasher(src)
	nakedFilename, err := Save(copy, filename, storageDir, docTypes)
	return nakedFilename, hex.EncodeToString(hasher.Sum(nil)), err
}

// SaveAvatar saves avatar image file to disk
func SaveAvatar(src io.Reader, filename string, storageDir string) (string, errstack.E) {
	return Save(src, filename, storageDir, avatarTypes)
}

// Save saves file to the filesystem.
// Returns saved filename without storage path and file hash
func Save(src io.Reader, filename string, storageDir string, fileTypes []string) (string, errstack.E) {
	if srcClose, ok := src.(io.ReadCloser); ok {
		defer errstack.CallAndLog(logger, srcClose.Close)
	}
	fileExtension, err := getFileExtension(filename, fileTypes)
	if err != nil {
		return "", errstack.WrapAsInf(err, "Bad file extension")
	}
	nakedFilename := uuid.NewUUID().String() + fileExtension
	absoluteFilename := filepath.Join(storageDir, nakedFilename)
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		return "", errstack.WrapAsInf(err, "Can't create upload directory")
	}
	dest, err := os.Create(absoluteFilename)
	if err != nil {
		return "", errstack.WrapAsInf(err, "Can't create file")
	}
	defer errstack.CallAndLog(logger, dest.Close) // idempotent, okay to call twice
	_, err = io.Copy(dest, src)
	return nakedFilename, errstack.WrapAsInf(err, "Can't write to file")
}

func readerHasher(input io.Reader) (io.Reader, hash.Hash) {
	h, err := blake2s.New256(nil)
	if err != nil {
		logger.Error("Impossible error, New256(nil) shouldn't return error", err)
	}
	srcTee := io.TeeReader(input, h)
	return srcTee, h
}
