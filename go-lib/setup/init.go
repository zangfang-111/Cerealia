/* Package providings setup / initialization of common structures, drivers, ...
 * for other tools and applications
 */

// Package setup contains all basic settings for projects
package setup

import (
	"os"
	"strings"

	bat "github.com/robert-zaremba/go-bat"
	"github.com/robert-zaremba/log15/log15setup"
)

// GitVersion should be substituted during build time by the git version. This is done
// using go linker flags:
// -ldflags "-X bitbucket.org/cerealia/apps/go-lib/setup.GitVersion=$(git describe)
var GitVersion = "unset:go-lib/setup.GitVersion"

// envName is the application stage environment name (eg production, dev, backstage, ...).
// It is set using `CEREALIA_ENV` environment variable.
var envName string

// RootDir is an absolute path to the project root directory
var RootDir string

// init initializes packages
func init() {
	envName = os.Getenv("CEREALIA_ENV")
	envName = strings.ToLower(envName)
	if envName == "" || strings.HasPrefix(envName, "dev") {
		envName = "development"
	} else if strings.HasPrefix(envName, "prod") {
		envName = "production"
	}
	log15setup.MustAppName(envName, "env name")
	var err error
	RootDir, err = bat.FindRoot()
	if err != nil {
		logger.Fatal("Can't set a Root directory", err)
	}
}
