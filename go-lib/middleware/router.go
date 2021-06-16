package middleware

import (
	"fmt"
	"net/http"

	driver "github.com/arangodb/go-driver"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	"github.com/scale-it/go-web/ozzohandlers"
)

func panicLogger(format string, a ...interface{}) {
	logger.Crit(fmt.Sprintf("Panic occured. "+format, a...))
}

// StdRouter defines standard, default router
// prefix is the routing group which will be use in the service. The leading `/` is added
// automatically if absent in the prefix.
func StdRouter(db driver.Database, isProduction bool) (*routing.Router, *routing.RouteGroup) {
	router := routing.New()
	logTrace := ozzohandlers.LogTrace{Logger: logger}
	router.Use(
		fault.PanicHandler(panicLogger),
		logTrace.LogTrace,
		restError)
	if !isProduction {
		logger.Info("CORS: allowing all origins")
		router.Use(cors.Handler(cors.AllowAll))
	}
	r := router.Group("")
	r.Use(
		WithAuth(db),
		content.TypeNegotiator(content.HTML, content.JSON),
		slash.Remover(http.StatusMovedPermanently))
	r.Get("/health-check", func(ctx *routing.Context) error {
		return ctx.Write("OK")
	})
	return router, r
}
