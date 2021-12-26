DC ?= docker-compose
GO := go
FN := frontend/
TEST_DB_CONN_STRING ?= postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable
GIT_REF    := $(shell git describe --all | sed  -e  's%tags/%%g'  -e 's%/%.%g' )
GIT_COMMIT := $(shell git rev-parse --short HEAD)

go_packages = $(GO) list ./... | grep -v /test | xargs

export GO_MODULE=$(shell head -1 go.mod | cut -d' ' -f 2)

ifndef GITHUB_REF
	DATE := ${shell date +%s}
	GITHUB_REF := ${GIT_REF}-${GIT_COMMIT}-$(DATE)
endif

.PHONY: lint
lint:
	$(GO) vet

.PHONY: start
start: stop
	GOPATH=$$(go env GOPATH) $(DC) up

.PHONY: stop
stop:
	$(DC) down

.PHONY: start-db
start-db: stop
	GOPATH=$$(go env GOPATH) $(DC) up db

.PHONY: start-api
start-api: stop
	GOPATH=$$(go env GOPATH) $(DC) up api

.PHONY: start-frontend
start-frontend: stop
	cd $(FN) && ng serve

.PHONY: gen-graph
gen-graph:
	$(GO) run github.com/99designs/gqlgen generate

.PHONY: test
test:
	$(GO) test -v -race -short $(shell $(call go_packages))

.PHONY: build
build: clean
	CGO_ENABLED=0  $(GO) build -o bin/api -a -ldflags '-X $(GO_MODULE)/internal.Version=$(GITHUB_REF)'

.PHONY: clean
clean:
	rm -f bin/*

.PHONY: test-integration
test-integration:
	TEST_DB_CONN_STRING="$(TEST_DB_CONN_STRING)" \
	$(GO) test -count=1 -v ./test/integration

.PHONY: test-e2e
test-e2e:
	$(GO) test -count=1 -v ./test/e2e

.PHONY: test-frontend
test-frontend:
	cd $(FN) && ng test

.PHONY: test-local
test-local: stop test
	GOPATH=$$(go env GOPATH) $(DC) up -d db
	make test-integration
	$(DC) down
	$(DC) up -d
	make test-e2e
	$(DC) down