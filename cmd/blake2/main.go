package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"bitbucket.org/cerealia/apps/go-lib/setup"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/log15"
	"golang.org/x/crypto/blake2s"
)

var logger = log15.Root()

func main() {
	setup.FlagSimpleInit("blake2s", "filename")
	flag.Parse()

	fmt.Println("argument:", flag.Arg(0))
	fmt.Println("hash:", compute(flag.Arg(0)))
}

func compute(filename string) string {
	h, err := blake2s.New256(nil)
	if err != nil {
		logger.Error("can't create a hasher", err)
		return ""
	}
	f, err := os.Open(filename)
	defer errstack.CallAndLog(logger, f.Close)
	if err != nil {
		logger.Error("can't open file", err)
		return ""
	}
	if _, err = io.Copy(h, f); err != nil {
		logger.Error("can't read file content", err)
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
