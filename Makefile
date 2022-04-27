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
		@echo "--> Installing zetacored & zetaclientd"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetacored
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

install-zetaclient: go.sum
		@echo "--> Installing zetaclientd"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

install-zetacore: go.sum
		@echo "--> Installing zetacored"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetacored

install-mockmpi:
	@echo "--> Installing MockMPI"
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/mockmpi


go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify


