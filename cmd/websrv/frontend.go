package main

import (
	"net/http"

	"bitbucket.org/cerealia/apps/cmd/websrv/config"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/file"
)

const frontendDir = "./browser/build/"

// SetFrontendRoutes serves the app
func SetFrontendRoutes(r *routing.Router) {
	var fileStorage = config.F.FileStorageDir.String()
	r.Get("/", redirectToHome)
	r.Get("/view/*", file.Content(frontendDir+"index.html"))
	r.Get("/assets/public/*", file.Server(
		file.PathMap{"/assets": ""},
		file.ServerOptions{RootPath: fileStorage}))
	// this must be the last handler
	r.Get("/*", file.Server(
		file.PathMap{"": ""},
		file.ServerOptions{RootPath: frontendDir}))
}

func redirectToHome(c *routing.Context) error {
	http.Redirect(c.Response, c.Request, "/view/home", http.StatusMovedPermanently)
	return nil
}
