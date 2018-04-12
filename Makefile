#!/usr/bin/env bash

default: format
	make start
	make log

gofmtvalidate:
	scripts/gofmt_validate.sh;
	scripts/gotest.sh;

format:
	goimports -w -local github.com/modeneis/coind ./src
	gofmt -s -w ./src

start: stop
	echo "will compile app";
	godep go build -o cryptobot main.go;
	echo "will start app with dev credentials";
	MONGODB_HOST=localhost MONGODB_PASSWORD= PORT=3030 WATCH_MARKETS=true WATCH_TIME=10 DEBUG=true \
	nohup ./cryptobot &
	make log;

stop:
	echo "will stop dev app"
	pidof cryptobot |awk '{print $1}'| xargs kill | true;

killall:
	lsof -i tcp:3000 | awk 'NR!=1 {print $2}' | xargs kill -9 | true;
	echo "will sleep 5 secs ";

status:
	ps -ef |grep cryptobot

log:
	tail -f ./nohup.out

test:
	go test ./src/... -timeout=5m -cover

test-race: ## Run tests with -race. Note: expected to fail, but look for "DATA RACE" failures specifically
	go test ./src/... -timeout=5m -race

mongodb:
	go test -v ./mongodb/mongodb_index_test.go;

cover: ## Runs tests on ./src/ with HTML code coverage
	@echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	go get -u github.com/alecthomas/gometalinter
	gometalinter --vendored-linters --install

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	gometalinter --deadline=3m -j 2 --disable-all --tests --vendor \
		-E deadcode \
		-E errcheck \
		-E gas \
		-E goconst \
		-E gofmt \
		-E goimports \
		-E golint \
		-E ineffassign \
		-E interfacer \
		-E maligned \
		-E megacheck \
		-E misspell \
		-E nakedret \
		-E structcheck \
		-E unconvert \
		-E unparam \
		-E varcheck \
		-E vet \
		./...

.PHONY: all test mongodb check permission
