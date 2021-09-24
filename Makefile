SHELL := /bin/bash # this is for sourcing env.sh
.PHONY: build


PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=meta-chaind \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=meta-chaind \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=metacli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: go.sum
		@echo "--> Installing meta-chaind & meta-noded"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/metacored
#		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/meta-noded
#		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/config

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

reset:
	sudo rm -rf ~/.metachain
	./build/scripts/make-localnet.sh

down:
	docker-compose -f build/mocknet/metachain.yml down --remove-orphans

up:
	source build/scripts/env.sh && docker-compose -f build/mocknet/metachain.yml up --remove-orphans -d


build:
	docker-compose -f build/mocknet/metachain.yml build
