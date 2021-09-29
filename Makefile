DC ?= docker-compose

.PHONY: start
start: stop
	$(DC) up --build

.PHONY: stop
stop:
	$(DC) down

.PHONY: start-db
start-db: stop
	$(DC) up db

start-api: stop
	$(DC) up api --build