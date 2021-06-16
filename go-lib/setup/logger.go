package setup

import (
	"github.com/robert-zaremba/log15"
	"github.com/robert-zaremba/log15/log15setup"
)

var logger = log15.Root()

// MustLogger setups logger
func MustLogger(appname, rollbartoken string) {
	log15setup.MustLogger(envName, appname, GitVersion, rollbartoken, "sec",
		"DEBUG", *flagLogColored)
}
