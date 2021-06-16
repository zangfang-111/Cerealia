// Package config contains flag data for web server
package config

import (
	"bitbucket.org/cerealia/apps/go-lib/setup"
	"github.com/robert-zaremba/flag"
)

// AppFlags is a set of websrv configuration flags
type AppFlags struct {
	setup.SrvFlags
	FileStorageDir setup.PathFlag
}

// F is the only official AppFlags instance
var F = AppFlags{
	setup.NewSrvFlags(),
	setup.PathFlag{Path: "/tmp/cerealia-files"},
}

func init() {
	flag.Var(&F.FileStorageDir, "file-storage-path",
		"path to store trade related and other files")
}

// Check validates the flags. Implements `flag.Checker` interface.
func (af *AppFlags) Check() error {
	return setup.FlagCheckMany(af.SrvFlags, af.FileStorageDir)
}
