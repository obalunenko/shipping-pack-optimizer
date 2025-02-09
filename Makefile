APP_NAME?=server
SHELL := env APP_NAME=$(APP_NAME) $(SHELL)

IMAGE_NAME?=shipping-pack-optimizer-$(APP_NAME)
SHELL := env IMAGE_NAME=$(IMAGE_NAME) $(SHELL)

BIN_DIR?=$(CURDIR)/bin


COMMIT_MAKE?=$(shell git rev-parse HEAD)
SHELL := env COMMIT_MAKE=$(COMMIT_MAKE) $(SHELL)
SHORTCOMMIT_MAKE?=$(shell git rev-parse --short HEAD)
SHELL := env SHORTCOMMIT_MAKE=$(SHORTCOMMIT_MAKE) $(SHELL)
DATE_MAKE?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
SHELL := env DATE_MAKE=$(DATE_MAKE) $(SHELL)
VERSION_MAKE?=$(shell git tag | sort -V | tail -1)
SHELL := env VERSION_MAKE=$(VERSION_MAKE) $(SHELL)

GOVERSION:=1.23

TEST_DISCARD_LOG?=false
SHELL := env TEST_DISCARD_LOG=$(TEST_DISCARD_LOG) $(SHELL)

format-code: swagger-fmt fmt goimports
.PHONY: format-code

fmt:
	@echo "Formatting code..."
	./scripts/style/fmt.sh
.PHONY: fmt

goimports:
	@echo "Formatting code..."
	./scripts/style/fix-imports.sh
.PHONY: goimports

vet:
	@echo "Vetting code..."
	@go vet ./...
	@echo "Done"
.PHONY: vet

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Done"
.PHONY: test

build:
	@echo "Building..."
	@./scripts/build/app.sh
	@echo "Done"
.PHONY: build

run: vendor build
	@echo "Running..."
	@${BIN_DIR}/$(APP_NAME)
	@echo "Done"
.PHONY: run

vendor:
	@echo "Vendoring..."
	@go mod tidy && go mod vendor
	@echo "Done"
.PHONY: vendor

docker-build: vendor
	@echo "Building docker image..."
	@IMAGE_DESCRIPTION="$$(cat README.md)" docker buildx bake
	@echo "Done"
.PHONY: docker-build

docker-run: docker-build
	@echo "Running docker image..."
	@docker compose -f ./deployments/compose/compose.yaml up
	@echo "Done"
.PHONY: docker-run

docker-stop:
	@echo "Stopping docker image..."
	@docker compose -f ./deployments/compose/compose.yaml down
	@echo "Done"
.PHONY: docker-stop

## Release
release:
	./scripts/release/release.sh
.PHONY: release

## Release local snapshot
release-local-snapshot:
	./scripts/release/local-snapshot-release.sh
.PHONY: release-local-snapshot

## Check goreleaser config.
check-releaser:
	./scripts/release/check.sh
.PHONY: check-releaser

## Issue new release.
new-version: vet test build docker-build
	./scripts/release/new-version.sh
.PHONY: new-release

## Bump go version
bump-go-version:
	./scripts/bump-go.sh $(GOVERSION)
.PHONY: bump-go-version

## Generate swagger docs
swagger-gen:
	./scripts/swagger-docs.sh
.PHONY: swagger-gen

## Format swagger annotations
swagger-fmt:
	./scripts/style/swagger-fmt.sh
.PHONY: swagger-fmt

## Copy .env file from .env.example
copy_env:
	@cp .env.example .env
.PHONY: copy_env




