SHELL=/bin/bash

CHECK=' \033[32m✔\033[39m'
HR=\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#
IBLACK='\e[1;30m'
GREEN='\e[0;32m'
NC='\e[0m'              # No Color

VERSION=`git log  -n 1 --format="%h-%cI"`  #`git describe`
# GOLIB=$(shell find go-lib -type f -name \*.go)
GOCLEAN=go clean

.PHONY: all setup-dev install-deps \
	vendor-check \
	go-generate gqlgen \
	lint lint-go lint-go-mega lint-graphql \
	build build-example build-websrv build-seed \
	db-seed db-migrate \
	test clean

all: install-deps build

CONFIG_DIR="config/config.ini"
FILE_STORAGE_PATH=`grep file-storage-path ${CONFIG_DIR} | awk '{print $$2}'`
CONTRACTDATA_STORAGE_PATH="${FILE_STORAGE_PATH}/trade-docs"
AVATARS_STORAGE_PATH="${FILE_STORAGE_PATH}/public/user/avatars"

setup-dev-githooks:
	@test ! -e .git/hooks/pre-commit && \
		echo INFO: Copying pre-commit hook to .git/hooks/pre-commit && \
		ln -vs ../../scripts/dev/git-hooks/pre-commit .git/hooks/pre-commit || \
		echo INFO: git pre-commit hook already exists. No action.



# validate-go-generate-is-uptodate:
# 	@./test/validate_go_generate_is_up_to_date.sh || exit 1

gqlgen:
	@cd api && gorunpkg github.com/99designs/gqlgen -v generate

###############################
# linting

lint: lint-go lint-graphql

lint-graphql:
	@graphql lint -p webserv


define golintx
	find $(1) -maxdepth 1 -name '[^.#]*.go' ! -name '*_test.go' ! -name '*_string.go' ! -name 'benchmark_*.go'  ! -name '*_mock.go' ! -name '*_gen.go' | xargs golint -set_exit_status;
endef

# lint-go: GODIRS = $(shell find go-lib cmd -name '*.go' ! -name '*mock.go' -exec dirname '{}' \; | sort -u)
lint-go:
	@echo -e $(IBLACK)linting go libs and apps... $(NC)
# 	@for d in $(GODIRS); do \
# 		$(call golintx,$$d) done;
# The command below requires bash shell
	@revive go-lib/... cmd/... | tee >(test -s /dev/stdout)

	@echo -e $(IBLACK)go vet go-lib... cmd...$(NC)
	@go tool vet -all go-lib
	@echo -e $(IBLACK)errcheck go-lib... cmd... $(NC)
	@errcheck -ignoregenerated -ignore 'fmt:[FS]?[Pp]rint*' ./go-lib/... ./cmd/...


lint-browser:
	@(cd browser && make lint)

###############################
# testing

test:
	@go test ./go-lib/... ./cmd/...



###############################

define _build
	@go generate ./go-lib/...
# GOOS=darwin GOOS=windows GOARCH=amd64
	GOBIN=`pwd`/bin go install -v \
		-ldflags "-X bitbucket.org/cerealia/apps/go-lib/setup.GitVersion=$(VERSION) -w" \
		./cmd/$(1)
	@echo -e "> build completed" $(CHECK)
endef

build:
	@$(call _build,"...")

build-websrv:
	@$(call _build,"websrv")
	./bin/websrv

build-example:
	@$(call _build,"_example")

build-stellar-cleanup:
	@$(call _build,"stellar_cleanup")


# db migration and seeding

db-migrate:
	@$(call _build,"migration")
	./bin/migration

build-seed:
	@$(call _build,"seed")

db-seed: build-seed
	@echo "Copying fixture contracts to:" $(CONTRACTDATA_STORAGE_PATH)
	@mkdir -vp $(CONTRACTDATA_STORAGE_PATH)
	@cp -R fixtures/fixture-docs/* $(CONTRACTDATA_STORAGE_PATH)
	@echo "Copying avatars to:" $(AVATARS_STORAGE_PATH)
	@mkdir -vp $(AVATARS_STORAGE_PATH)
	@cp -R fixtures/avatars/* $(AVATARS_STORAGE_PATH)
	./bin/seed

# run example:
#    ./bin/contracts
