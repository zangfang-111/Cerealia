# apps

This is a mono-repo of all Cerealia applications.


## Structure

+ `/browser` – the browser web application (React).
+ `/cmd/<app-name>` – Go main packages (applications). For web service it should contain all routers and handlers,
+ [optional] `cmd/<app-name>/<app-package>` – some package specific to the app (like some set of controllers which we don’t like to put into domain packages.
+ [for future] `cmd/<app-name>/<grpc-interface>` – common interface for grpc protocol.
+ `go-lib/mode/<domain_name>` (eg `user`) – domain packages
+ `/go-lib`– Go libraries. Some generic libraries in the future can be moved into a separate repository under Cerealia or the main author and will require individual assessment.
+ `/submodules` - linked git submodules.

Specific README files:

+ ./go-lib/model/README.md


## Setting up

You have to set-up the following environment:

1. Install `go >= 1.10` - please use package from your package manager or [PPA](https://github.com/golang/go/wiki/Ubuntu) (only if needed).
    * Make sure that `go` executable is in your `PATH`.
    * Make sure your `GOBIN` (usually `$HOME/go/bin`) is in your `PATH` - this is the place where GO binnaries are installed by default.
1. Install `yarn >= 1.10` and `make`.
1. Install and setup the `ArangoDB >= 3.4`, [instructions](https://docs.arangodb.com/3.4/Manual/GettingStarted/Installation.html).
1. copy the config file and edit it according to your local configuraiton:

        cd config
        cp config.ini.example config.ini
        $EDITOR config.ini


### Backend

We use `make` to organize build and setup commands.

1. For static compilation: `libc-static` (eg: `glibc-static`).
1. Make sure you have an access to this repository.
1. Clone this repository into your `GOPATH` (`go env`). Install dependencies:

	git clone git@bitbucket.org:cerealia/apps.git $GOPATH/src/bitbucket.org/cerealia/apps
	cd ~/go/src/bitbucket.org/cerealia/apps
	make install-deps

1. (alternative). Use `go get` (this will work only if the repo is open). Then, to be able to push, update your git config (`.git/config`) to use the git over ssh, rather https.

    go get bitbucket.org/cerealia/apps

#### Development:

This is not required for building application. You should make this step if you are a
developer and want to update this project.

	make setup-dev install-deps setup-dev-githooks

Before commiting any changes run:

	make lint

### Frontend

Please refere to `browser/README.md` file for details.
