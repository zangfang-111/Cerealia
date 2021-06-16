
.PHONY: all \
	install-deps \
	lint check-flow-headers \


install-deps:
	@yarn install --dev --prefer-offline

check-flow-headers:
	@awk "FNR == 1 && ! /^\/\/ @flow$$/ {print FILENAME}; {nextfile}" src/**/*.js* src/*.js*

lint: check-flow-headers
	@yarn run lint

lint-fix:
	@yarn run lint-fix

check-all: lint
	@yarn run flow
