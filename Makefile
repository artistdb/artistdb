DC ?= docker-compose
GO := go

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
	$(GO) test -v -race -short ./...

.PHONY: build
build: clean
	$(GO) build -o bin/api

.PHONY: clean
clean:
	rm -f bin/*