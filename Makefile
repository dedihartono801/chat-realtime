GOPATH ?= $(HOME)/go

run-api:
	go run cmd/http/main.go 

run-worker:
	go run cmd/worker/main.go

run-websocket:
	go run cmd/websocket/main.go

format:
	go fmt ./...

api:
	DC_APP_ENV=dev $(GOPATH)/bin/reflex -s -r '\.go$$' make format run-api

worker:
	DC_APP_ENV=dev $(GOPATH)/bin/reflex -s -r '\.go$$' make format run-worker

websocket:
	DC_APP_ENV=dev $(GOPATH)/bin/reflex -s -r '\.go$$' make format run-websocket

migration-up:
	migrate -database "postgresql://postgres:postgres@localhost:5432/chat?sslmode=disable" -path migrations up

migration-down:
	migrate -database "postgresql://postgres:postgres@localhost:5432/chat?sslmode=disable" -path migrations down

migration $$(enter):
	@read -p "Migration name:" migration_name; \
	migrate create -ext sql -dir migrations $$migration_name

test-cov:
	go test -coverprofile=cover.out ./... && go tool cover -html=cover.out -o cover.html