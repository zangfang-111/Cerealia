CHECK=' \033[32m✔\033[39m'
HR=\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#\#
IBLACK='\e[1;30m'
GREEN='\e[0;32m'
NC='\e[0m'              # No Color

install-deps:
	@echo -en $(IBLACK)fetching API tools ... $(NC)
	@go get -v -u github.com/bukalapak/snowboard

apib-websrv:
	@snowboard html -o websrv.html websrv.apib
