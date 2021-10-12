DC ?= docker-compose
GO := go
TEST_DB_CONN_STRING ?= postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable

go_packages = $(GO) list ./... | grep -v /test | xargs

.PHONY: start
start: stop
	$(DC) up --build

.PHONY: stop
stop:
	$(DC) down

.PHONY: start-db
start-db: stop
	$(DC) up db

.PHONY: start-api
start-api: stop
	$(DC) up api --build

.PHONY: test
test:
	$(GO) test -v -race -short $(shell $(call go_packages))
.PHONY: build
build: clean
	$(GO) build -o bin/api

.PHONY: clean
clean:
	rm -f bin/*

.PHONY: test-integration
test-integration:
	TEST_DB_CONN_STRING="$(TEST_DB_CONN_STRING)" \
	$(GO) test -count=1 -v ./test/integration

.PHONY: test-local
test-local: stop test
	$(DC) up -d db
	sleep 5 # Wait for DB
	make test-integration
	$(DC) down