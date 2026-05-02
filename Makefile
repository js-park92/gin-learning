.PHONY: run build test lint fmt tidy migrate-up migrate-down migrate-version migrate-create docker-up docker-down docker-logs

run:
	go run .

build:
	go build -o bin/server .

test:
	go test ./... -v

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

tidy:
	go mod tidy

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down $(steps)

migrate-version:
	go run ./cmd/migrate version

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=<migration_name>"; exit 1; fi
	@n=$$(find migrations -name '*.up.sql' 2>/dev/null | wc -l | tr -d ' '); \
	seq=$$(printf "%06d" $$((n + 1))); \
	touch migrations/$${seq}_$(name).up.sql migrations/$${seq}_$(name).down.sql; \
	echo "Created migrations/$${seq}_$(name).up.sql and migrations/$${seq}_$(name).down.sql"

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-build:
	docker compose build --no-cache
