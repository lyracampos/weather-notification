GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
SWAGGER := docker run --rm -e GOPATH=$$(go env GOPATH):/go -v $$(pwd):$$(pwd) -w $$(pwd) quay.io/goswagger/swagger:v0.31.0

MIGRATE := go run -tags='postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1
MIGRATIONS_PATH ?= './internal/gateways/database/postgres/migrations'
PG_CONNECTION_STRING ?= 'postgres://user:password@localhost:5432/weather_notification?sslmode=disable'

lint:
	$(GOLANGCI_LINT) run --fix

build:
	go build ./...

start:
	make start-infra && sleep 10 && make start-app

start-infra:
	docker compose up -d postgres rabbitmq
	until docker exec postgres pg_isready; do echo 'Waiting for postgres server...' && sleep 1; done
	make migration/up

start-app:
	docker compose up -d api worker websocket
	until docker exec postgres pg_isready; do echo 'Waiting for postgres server...' && sleep 1; done
	make migration/up

stop-infra:
	docker compose down postgres rabbitmq

stop-app:
	docker compose down api worker rabbitmq websocket

run/api:
	go run cmd/main.go -e api -c ./configs/config.yaml

run/worker-web:
	go run cmd/main.go -e worker -t web -c ./configs/config.yaml

run/websocket-client:
	go run cmd/main.go -e websocket -c ./configs/config.yaml

create/migration:
	$(MIGRATE) create -seq -ext sql -dir $(MIGRATIONS_PATH) $(MIGRATION_NAME)

migration/up:
	$(MIGRATE) -path $(MIGRATIONS_PATH) --database $(PG_CONNECTION_STRING) up

swagger/generate:
	$(GOSWAGGER) generate spec -o ./swagger.yaml --scan-models
# fmt