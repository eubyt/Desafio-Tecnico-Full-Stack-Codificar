ifneq (,$(wildcard ./.env))
    include .env
    export
else ifneq (,$(wildcard ../.env))
    include ../.env
    export
endif

.PHONY: all build run seed test test-race test-e2e sqlc swagger openapi clean

SWAG := $(shell command -v swag 2>/dev/null || echo $(shell go env GOPATH)/bin/swag)

all: test build

openapi:
	@mkdir -p ../api
	$(SWAG) init -g cmd/api/main.go -o ../api -ot yaml --parseDependency --parseInternal --useStructName
	@mv ../api/swagger.yaml ../api/openapi.yaml

sqlc:
	sqlc generate

build:
	@mkdir -p bin
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

seed:
	go run ./cmd/seed

test:
	go test -p 1 -v ./...

test-race:
	go test -p 1 -v -race ./...

test-e2e:
	go test -v -tags=e2e ./test/e2e/...

clean:
	rm -rf bin/
