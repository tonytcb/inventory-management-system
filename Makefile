SHELL := /bin/bash

.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## up: starts the application
up:
	docker-compose up app
	docker-compose down --remove-orphans

## up: starts the application exposing its HTTP port
down:
	docker-compose down --remove-orphans

## clean: clean up all docker containers
clean: down
	docker ps -aq | xargs docker stop | xargs docker rm
	rm -rf ./dbdata

## lint: runs linter for all packages
lint:
	@ docker run  --rm -v "`pwd`:/workspace:cached" -w "/workspace/." golangci/golangci-lint:v1.57-alpine golangci-lint run  --timeout 5m

generate-mocks:
	mockery --all --dir ./internal --output ./test/mocks --exported --case=underscore

## tests-ci: runs all tests, with no Docker
tests:
	# flag ldflags is used to avoid warnings related to race flag: https://github.com/golang/go/issues/61229
	go test -p 1 -race -ldflags=-extldflags=-Wl,-ld_classic -tags=integration ./...

## Creates a migration file
## Usage: make create-migration FILE=create_table_test
create-migration:
	migrate create -ext sql -dir db/migrations -seq $(FILE)
