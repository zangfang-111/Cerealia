package main

import (
	"net/http"

	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/setup"
	"github.com/robert-zaremba/log15"
	"github.com/robert-zaremba/log15/rollbar"
)

var flags = setup.NewSrvFlags()
var logger = log15.Root()

const serviceName = "direct-pledge"

func setupFlags() {
	setup.FlagSimpleInit(serviceName, "", flags)
}

func main() {
	setupFlags()
	defer rollbar.WaitForRollbar(logger)

	handler, _ := middleware.StdRouter(serviceName)
	logger.Info("Example app listening at", "port", *flags.Port)
	if err := http.ListenAndServe(":"+*flags.Port, handler); err != nil {
		logger.Error("Can't initiate HTTP service", err)
	}
}
