SHELL := /bin/bash # this is for sourcing env.sh
.PHONY: build


PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=metachain \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=metacored \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=metaclientd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'
TEST_DIR?="./..."
TEST_BUILD_FLAGS :=  -tags mocknet



all: install


test-coverage:
	@go test ${TEST_BUILD_FLAGS} -v -coverprofile coverage.out ${TEST_DIR}

coverage-report: test-coverage
	@go tool cover -html=cover.txt

test:
	@go test ${TEST_BUILD_FLAGS} ${TEST_DIR}



install: go.sum
		@echo "--> Installing metacored & metaclientd"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/metacored
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/metaclientd

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify


reset:
	sudo rm -rf ~/.metachain
	./build/scripts/make-localnet.sh

down:
	docker-compose -f build/mocknet/metachain.yml down --remove-orphans

up:
	source build/scripts/env.sh && docker-compose -f build/mocknet/metachain.yml up --remove-orphans -d


build:
	docker-compose -f build/mocknet/metachain.yml build
