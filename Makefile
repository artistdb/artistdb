DC ?= docker-compose

.PHONY: start
start: stop
	$(DC) up

.PHONY: stop
stop:
	$(DC) down

.PHONY: start-db
start-db: stop
	$(DC) up -d db
