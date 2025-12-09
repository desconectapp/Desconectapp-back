SQLC_CONFIG=src/sqlc.yaml

COMPOSE_FILE=docker-compose-dev.yaml
MAIN_SERVICE=app 

TEST_COMPOSE_FILE=docker-compose-test.yaml

dev: sqlc
	docker compose -f $(COMPOSE_FILE) up -d --build

sqlc:
	docker run --rm -u $(shell id -u):$(shell id -g) -v ./src:/src -w /src sqlc/sqlc generate -f /src/sqlc.yaml

logs:
	docker compose -f $(COMPOSE_FILE) logs -f $(MAIN_SERVICE)

down:
	docker compose -f $(COMPOSE_FILE) down -t1

reset:
	docker compose -f $(COMPOSE_FILE) down -v

test: sqlc
	sudo rm -rf database_test
	docker compose -f $(TEST_COMPOSE_FILE) up --build --abort-on-container-exit --exit-code-from test

down_test:
	docker compose -f $(TEST_COMPOSE_FILE) down -v -t1

seed:
	python3 seed_desconectapp.py