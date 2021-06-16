package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/cerealia/apps/cmd/websrv/config"
	"bitbucket.org/cerealia/apps/cmd/websrv/trades"
	"bitbucket.org/cerealia/apps/cmd/websrv/users"
	"bitbucket.org/cerealia/apps/go-lib/gql"
	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/resolver"
	"bitbucket.org/cerealia/apps/go-lib/setup"
	dbs "bitbucket.org/cerealia/apps/go-lib/setup/arangodb"
	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource/txsourceimpl"
	"github.com/99designs/gqlgen/handler"
	driver "github.com/arangodb/go-driver"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/log15"
	"github.com/robert-zaremba/log15/rollbar"
)

var logger = log15.Root()
var db driver.Database

func main() {
	defer rollbar.WaitForRollbar(logger)
	setup.FlagSimpleInit("websrv", "", &config.F)
	ctx := context.Background()
	var err error
	db, err = dbs.GetDb(ctx)
	if err != nil {
		logger.Fatal("Can't connect database", err)
	}
	stellarDriver, err := stellar.NewDriver(*config.F.StellarNetwork)
	if err != nil {
		logger.Fatal("Can't build stellar.Driver", err)
	}
	lockDriver := txsourceimpl.NewDriver(db, time.Duration(*config.F.SCAddrLockDuration)*time.Second)
	router, err := buildRouter(stellarDriver, lockDriver)
	if err != nil {
		logger.Fatal("Can't build router", err)
	}
	http.Handle("/", router)
	logger.Info(fmt.Sprintf("connect to http://localhost:%s/graphiql/ for GraphQL playground", *config.F.Port))
	logger.Fatal("Server stopped",
		http.ListenAndServe(":"+*config.F.Port, nil))
}

func buildRouter(stellarDriver *stellar.Driver, txSourceDriver txsource.Driver) (http.Handler, error) {
	recovery := handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
		logger.Crit("Unhandled exception", err)
		return errstack.NewInf("Internal server error")
	})
	graphQLLogging := handler.ErrorPresenter(middleware.GraphQLError)
	gqlconfig := gql.Config{
		Resolvers: resolver.NewResolver(db, stellarDriver, txSourceDriver),
	}
	router, rgroup := middleware.StdRouter(db, *config.F.Production)
	trades.SetTradeRoutes(rgroup.Group("/v1/trades"), stellarDriver, txSourceDriver)
	trades.SetTradeOfferRoutes(rgroup.Group("/v1/trade-offers"))
	users.SetUserRoutes(rgroup.Group("/v1/users"))
	const gqlEndpoint = "/query"
	rgroup.Any(gqlEndpoint, routing.HTTPHandlerFunc(
		handler.GraphQL(gql.NewExecutableSchema(gqlconfig), recovery, graphQLLogging)))
	rgroup.Get("/graphiql", routing.HTTPHandlerFunc(handler.Playground("GraphQL playground", gqlEndpoint)))
	SetFrontendRoutes(router)
	return router, nil
}
