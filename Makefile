.PHONY: build


PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell git describe --tags)
COMMIT := $(shell [ -z "${COMMIT_ID}" ] && git log -1 --format='%H' || echo ${COMMIT_ID} )
BUILDTIME := $(shell date -u +"%Y%m%d.%H%M%S" )
DOCKER ?= docker
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
GOFLAGS:=""

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=zetacore \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=zetacored \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=zetaclientd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/zeta-chain/zetacore/common.Name=zetacored \
	-X github.com/zeta-chain/zetacore/common.Version=$(VERSION) \
	-X github.com/zeta-chain/zetacore/common.CommitHash=$(COMMIT) \
	-X github.com/zeta-chain/zetacore/common.BuildTime=$(BUILDTIME) \
	-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb

BUILD_FLAGS := -ldflags '$(ldflags)' -tags PRIVNET,pebbledb,ledger
TESTNET_BUILD_FLAGS := -ldflags '$(ldflags)' -tags TESTNET,pebbledb,ledger

TEST_DIR?="./..."
TEST_BUILD_FLAGS := -tags TESTNET,pebbledb,ledger
PRIV_BUILD_FLAGS := -tags PRIVNET,pebbledb,ledger

clean: clean-binaries clean-dir clean-test-dir clean-coverage

clean-binaries:
	@rm -rf ${GOBIN}/zetacored
	@rm -rf ${GOBIN}/zetaclientd

clean-dir:
	@rm -rf ~/.zetacored
	@rm -rf ~/.zetacore

all: install

test-coverage-exclude-core:
	@go test ${TEST_BUILD_FLAGS} -v -coverprofile coverage.out $(go list ./... | grep -v /x/zetacore/)

test-coverage:
	-@go test ${TEST_BUILD_FLAGS} -v -coverprofile coverage.out ${TEST_DIR}

coverage-report: test-coverage
	@go tool cover -html=coverage.out -o coverage.html

clean-coverage:
	@rm -f coverage.out
	@rm -f coverage.html

clean-test-dir:
	@rm -rf x/crosschain/client/integrationtests/.zetacored
	@rm -rf x/crosschain/client/querytests/.zetacored
	@rm -rf x/observer/client/querytests/.zetacored

run-test:
	@go test ${TEST_BUILD_FLAGS} ${TEST_DIR}

test :clean-test-dir run-test

test-priv:
	@go test ${PRIV_BUILD_FLAGS} ${TEST_DIR}

gosec:
	gosec  -exclude-dir=localnet ./...

install-testnet: go.sum
		@echo "--> Installing zetacored & zetaclientd"
		@go install -mod=readonly $(TESTNET_BUILD_FLAGS) ./cmd/zetacored
		@go install -mod=readonly $(TESTNET_BUILD_FLAGS) ./cmd/zetaclientd

build-testnet-ubuntu: go.sum
		docker build -t zetacore-ubuntu --platform linux/amd64 -f ./Dockerfile-athens3-ubuntu .
		docker create --name temp-container zetacore-ubuntu
		docker cp temp-container:/go/bin/zetaclientd .
		docker cp temp-container:/go/bin/zetacored .
		docker rm temp-container

install: go.sum
		@echo "--> Installing zetacored & zetaclientd"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetacored
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

install-zetaclient: go.sum
		@echo "--> Installing zetaclientd"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

# running with race detector on will be slow
install-zetaclient-race-test-only-build: go.sum
		@echo "--> Installing zetaclientd"
		@go install -race -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

install-zetacore: go.sum
		@echo "--> Installing zetacored"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetacored

install-zetacore-testnet: go.sum
		@echo "--> Installing zetacored"
		@go install -mod=readonly $(TESTNET_BUILD_FLAGS) ./cmd/zetacored

install-smoketest: go.sum
		@echo "--> Installing orchestrator"
		@go install -mod=readonly $(BUILD_FLAGS) ./contrib/localnet/orchestrator/smoketest

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test-cctx:
	./standalone-network/cctx-creator.sh

init:
	./standalone-network/init.sh

run:
	./standalone-network/run.sh

chain-init: clean install-zetacore init
chain-run: clean install-zetacore init run
chain-stop:
	@killall zetacored
	@killall tail


chain-init-testnet: clean install-zetacore-testnet init
chain-run-testnet: clean install-zetacore-testnet init run

lint-pre:
	@test -z $(gofmt -l .)
	@GOFLAGS=$(GOFLAGS) go mod verify

lint: lint-pre
	@golangci-lint run

proto:
	@echo "--> Removing old Go types "
	@find . -name '*.pb.go' -type f -delete
	@echo "--> Generating new Go types from protocol buffer files"
	@bash ./scripts/protoc-gen-go.sh
.PHONY: proto

openapi:
	@echo "--> Generating OpenAPI specs"
	@bash ./scripts/protoc-gen-openapi.sh
.PHONY: openapi

specs:
	@echo "--> Generating module documentation"
	@go run ./scripts/gen-spec.go
.PHONY: specs

generate: proto openapi specs
.PHONY: generate

###############################################################################
###                                Docker Images                             ###
###############################################################################

zetanode:
	@echo "Building zetanode"
	$(DOCKER) build -t zetanode -f ./Dockerfile .
	$(DOCKER) build -t orchestrator -f contrib/localnet/orchestrator/Dockerfile.fastbuild .
.PHONY: zetanode

smoketest:
	@echo "DEPRECATED: NO-OP: Building smoketest"

start-smoketest:
	@echo "--> Starting smoketest"
	cd contrib/localnet/ && $(DOCKER) compose up -d

start-smoketest-upgrade:
	@echo "--> Starting smoketest with upgrade proposal"
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-upgrade.yml up -d

start-smoketest-p2p-diag:
	@echo "--> Starting smoketest in p2p diagnostic mode"
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-p2p-diag.yml up -d

stop-smoketest:
	@echo "--> Stopping smoketest"
	cd contrib/localnet/ && $(DOCKER) compose down --remove-orphans

stop-smoketest-p2p-diag:
	@echo "--> Stopping smoketest in p2p diagnostic mode"
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-p2p-diag.yml down --remove-orphans

stress-test: zetanode
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-stresstest.yml up -d

stop-stress-test:
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-stresstest.yml down --remove-orphans