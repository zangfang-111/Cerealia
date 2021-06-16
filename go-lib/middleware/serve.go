package middleware

import (
	"net/http"
	"time"

	"bitbucket.org/cerealia/apps/go-lib/setup"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/robert-zaremba/errstack"
)

// CreateServer creates default server to be used in our code.
func CreateServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 120 * time.Second, // client internet connection can be slow and we still want to send him data
	}
}

// ServeHTTP build and starts the http server.
func ServeHTTP(addr string, handler http.Handler) {
	logger.Info("Starting HTTP server at http://"+addr, "version", setup.GitVersion)
	srv := CreateServer(addr, handler)
	errstack.Log(logger, gracehttp.Serve(srv))
}
