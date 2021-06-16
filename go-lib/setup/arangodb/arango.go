// Package arangodb contains aragnodb settings and functions
package arangodb

import (
	"context"
	"strings"

	"bitbucket.org/cerealia/apps/go-lib/setup"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/flag"
)

var dbObj driver.Database
var dbURL setup.URLFlag

func init() {
	flag.Var(&dbURL, "arangodb-url", "//username:password@host:port/database")
	driver.Cause = errstack.Cause
}

// OpenDB initializes a ArangoDB client and connects to the DB.
func OpenDB(ctx context.Context) (driver.Database, errstack.E) {
	client, dbname, err := MkClient()
	if err != nil {
		return nil, err
	}
	db, err2 := client.Database(ctx, dbname)
	return db, errstack.WrapAsInf(err2,
		"Can't open a connection to an existing arangoDB database")
}

// MkClient creates an ArangoDB client
func MkClient() (driver.Client, string, errstack.E) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://" + dbURL.Host}})
	if err != nil {
		return nil, "", errstack.WrapAsInf(err, "Can't create DB network connection")
	}

	pwd, _ := dbURL.User.Password()
	client, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
		Authentication: driver.BasicAuthentication(
			dbURL.User.Username(), pwd)})
	return client, strings.Replace(dbURL.Path, "/", "", -1), errstack.WrapAsInf(err,
		"Can't create ArangoDB client")
}

// GetDb gets a arangodb instance connected
func GetDb(ctx context.Context) (driver.Database, errstack.E) {
	var err errstack.E
	if dbObj == nil {
		dbObj, err = OpenDB(ctx)
	}
	return dbObj, err
}
