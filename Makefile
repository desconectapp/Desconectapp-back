SQLC_CONFIG=src/sqlc.yaml

COMPOSE_FILE=docker-compose-dev.yaml
MAIN_SERVICE=app 

dev: sqlc
	docker compose -f $(COMPOSE_FILE) up -d --build

sqlc:
	docker run --rm -v ./src:/src -w /src sqlc/sqlc generate -f /src/sqlc.yaml

logs:
	docker compose -f $(COMPOSE_FILE) logs -f $(MAIN_SERVICE)

down:
	docker compose -f $(COMPOSE_FILE) down -t1

reset:
	docker compose -f $(COMPOSE_FILE) down -v
