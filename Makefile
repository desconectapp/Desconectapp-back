SQLC_CONFIG=src/sqlc.yaml

COMPOSE_FILE=docker-compose-dev.yaml
MAIN_SERVICE=app 

dev: sqlc
	docker compose -f $(COMPOSE_FILE) up -d --build

sqlc:
	sqlc generate -f $(SQLC_CONFIG)

logs:
	docker compose -f $(COMPOSE_FILE) logs -f $(MAIN_SERVICE)

down:
	docker compose -f $(COMPOSE_FILE) down -t1

reset:
	docker compose -f $(COMPOSE_FILE) down -v
