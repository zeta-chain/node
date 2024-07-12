.PHONY: build

VERSION := $(shell ./version.sh)
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
	-X github.com/zeta-chain/zetacore/pkg/constant.Name=zetacored \
	-X github.com/zeta-chain/zetacore/pkg/constant.Version=$(VERSION) \
	-X github.com/zeta-chain/zetacore/pkg/constant.CommitHash=$(COMMIT) \
	-X github.com/zeta-chain/zetacore/pkg/constant.BuildTime=$(BUILDTIME) \
	-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb

BUILD_FLAGS := -ldflags '$(ldflags)' -tags pebbledb,ledger

TEST_DIR?="./..."
TEST_BUILD_FLAGS := -tags pebbledb,ledger
HSM_BUILD_FLAGS := -tags pebbledb,ledger,hsm_test

export DOCKER_BUILDKIT := 1

# parameters for localnet docker compose files
# set defaults to empty to prevent docker warning
export E2E_ARGS := $(E2E_ARGS)

clean: clean-binaries clean-dir clean-test-dir clean-coverage

clean-binaries:
	@rm -rf ${GOBIN}/zetacored
	@rm -rf ${GOBIN}/zetaclientd

clean-dir:
	@rm -rf ~/.zetacored
	@rm -rf ~/.zetacore

all: install

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

###############################################################################
###                             Test commands                               ###
###############################################################################

run-test:
	@go test ${TEST_BUILD_FLAGS} ${TEST_DIR}

test :clean-test-dir run-test

test-hsm:
	@go test ${HSM_BUILD_FLAGS} ${TEST_DIR}

# Generate the test coverage
# "|| exit 1" is used to return a non-zero exit code if the tests fail
test-coverage:
	@go test ${TEST_BUILD_FLAGS} -coverprofile coverage.out ${TEST_DIR} || exit 1

coverage-report: test-coverage
	@go tool cover -html=coverage.out -o coverage.html

clean-coverage:
	@rm -f coverage.out
	@rm -f coverage.html

clean-test-dir:
	@rm -rf x/crosschain/client/integrationtests/.zetacored
	@rm -rf x/crosschain/client/querytests/.zetacored
	@rm -rf x/observer/client/querytests/.zetacored

###############################################################################
###                          Install commands                               ###
###############################################################################

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
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd-supervisor

install-zetaclient: go.sum
		@echo "--> Installing zetaclientd"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

install-zetacore: go.sum
		@echo "--> Installing zetacored"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetacored

# running with race detector on will be slow
install-zetaclient-race-test-only-build: go.sum
		@echo "--> Installing zetaclientd"
		@go install -race -mod=readonly $(BUILD_FLAGS) ./cmd/zetaclientd

install-zetatool: go.sum
		@echo "--> Installing zetatool"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetatool

###############################################################################
###                             Local network                               ###
###############################################################################

init:
	./standalone-network/init.sh

run:
	./standalone-network/run.sh

chain-init: clean install-zetacore init
chain-run: clean install-zetacore init run
chain-stop:
	@killall zetacored
	@killall tail

test-cctx:
	./standalone-network/cctx-creator.sh

###############################################################################
###                                 Linting            	                    ###
###############################################################################

lint-pre:
	@test -z $(gofmt -l .)
	@GOFLAGS=$(GOFLAGS) go mod verify

lint: lint-pre
	@golangci-lint run

lint-gosec:
	@bash ./scripts/gosec.sh

gosec:
	gosec  -exclude-dir=localnet ./...

###############################################################################
###                           		Formatting			                    ###
###############################################################################

fmt:
	@bash ./scripts/fmt.sh

###############################################################################
###                           Generation commands  		                    ###
###############################################################################

protoVer=0.13.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user $(shell id -u):$(shell id -g) $(protoImageName)

proto-format:
	@echo "--> Formatting Protobuf files"
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;
.PHONY: proto-format

typescript: proto-format
	@echo "--> Generating TypeScript bindings"
	@bash ./scripts/protoc-gen-typescript.sh
.PHONY: typescript

proto-gen: proto-format
	@echo "--> Removing old Go types "
	@find . -name '*.pb.go' -type f -delete
	@echo "--> Generating Protobuf files"
	@$(protoImage) sh ./scripts/protoc-gen-go.sh

openapi: proto-format
	@echo "--> Generating OpenAPI specs"
	@bash ./scripts/protoc-gen-openapi.sh
.PHONY: openapi

specs:
	@echo "--> Generating module documentation"
	@go run ./scripts/gen-spec.go
.PHONY: specs

docs-zetacored:
	@echo "--> Generating zetacored documentation"
	@bash ./scripts/gen-docs-zetacored.sh
.PHONY: docs-zetacored

mocks:
	@echo "--> Generating mocks"
	@bash ./scripts/mocks-generate.sh
.PHONY: mocks

# generate also includes Go code formatting
generate: proto-gen openapi specs typescript docs-zetacored mocks fmt
.PHONY: generate


###############################################################################
###                         Localnet                          				###
###############################################################################
start-localnet: zetanode start-localnet-skip-build

start-localnet-skip-build:
	@echo "--> Starting localnet"
	export LOCALNET_MODE=setup-only && \
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose.yml up -d

# stop-localnet should include all profiles so other containers are also removed
stop-localnet:
	cd contrib/localnet/ && $(DOCKER) compose --profile all down --remove-orphans

###############################################################################
###                         E2E tests               						###
###############################################################################

zetanode:
	@echo "Building zetanode"
	$(DOCKER) build -t zetanode --target latest-runtime -f ./Dockerfile-localnet .
	$(DOCKER) build -t orchestrator -f contrib/localnet/orchestrator/Dockerfile.fastbuild .
.PHONY: zetanode

install-zetae2e: go.sum
	@echo "--> Installing zetae2e"
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetae2e
.PHONY: install-zetae2e

start-e2e-test: zetanode
	@echo "--> Starting e2e test"
	cd contrib/localnet/ && $(DOCKER) compose up -d

start-e2e-admin-test: zetanode
	@echo "--> Starting e2e admin test"
	export E2E_ARGS="--skip-regular --test-admin" && \
	cd contrib/localnet/ && $(DOCKER) compose --profile eth2 -f docker-compose.yml up -d

start-e2e-performance-test: zetanode
	@echo "--> Starting e2e performance test"
	export E2E_ARGS="--test-performance" && \
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose.yml up -d

start-e2e-import-mainnet-test: zetanode
	@echo "--> Starting e2e import-data test"
	export ZETACORED_IMPORT_GENESIS_DATA=true && \
	export ZETACORED_START_PERIOD=15m && \
	cd contrib/localnet/ && ./scripts/import-data.sh mainnet && $(DOCKER) compose -f docker-compose.yml up -d

start-stress-test: zetanode
	@echo "--> Starting stress test"
	cd contrib/localnet/ && $(DOCKER) compose --profile stress -f docker-compose.yml up -d

###############################################################################
###                         Upgrade Tests              						###
###############################################################################

# build from source only if requested
ifdef UPGRADE_TEST_FROM_SOURCE
zetanode-upgrade: zetanode
	@echo "Building zetanode-upgrade from source"
	$(DOCKER) build -t zetanode:old -f Dockerfile-localnet --target old-runtime-source --build-arg OLD_VERSION='release/v17' .
.PHONY: zetanode-upgrade
else
zetanode-upgrade: zetanode
	@echo "Building zetanode-upgrade from binaries"
	$(DOCKER) build -t zetanode:old -f Dockerfile-localnet --target old-runtime --build-arg OLD_VERSION='https://github.com/zeta-chain/ci-testing-node/releases/download/v17.0.1-internal' .
.PHONY: zetanode-upgrade
endif


start-upgrade-test: zetanode-upgrade
	@echo "--> Starting upgrade test"
	export LOCALNET_MODE=upgrade && \
	export UPGRADE_HEIGHT=225 && \
	cd contrib/localnet/ && $(DOCKER) compose --profile upgrade -f docker-compose.yml -f docker-compose-upgrade.yml up -d

start-upgrade-test-light: zetanode-upgrade
	@echo "--> Starting light upgrade test (no ZetaChain state populating before upgrade)"
	export LOCALNET_MODE=upgrade && \
	export UPGRADE_HEIGHT=90 && \
	cd contrib/localnet/ && $(DOCKER) compose --profile upgrade -f docker-compose.yml -f docker-compose-upgrade.yml up -d


start-upgrade-test-admin: zetanode-upgrade
	@echo "--> Starting admin upgrade test"
	export LOCALNET_MODE=upgrade && \
	export UPGRADE_HEIGHT=90 && \
	export E2E_ARGS="--skip-regular --test-admin" && \
	cd contrib/localnet/ && $(DOCKER) compose --profile upgrade -f docker-compose.yml -f docker-compose-upgrade.yml up -d

start-upgrade-import-mainnet-test: zetanode-upgrade
	@echo "--> Starting import-data upgrade test"
	export LOCALNET_MODE=upgrade && \
	export ZETACORED_IMPORT_GENESIS_DATA=true && \
	export ZETACORED_START_PERIOD=15m && \
	export UPGRADE_HEIGHT=225 && \
	cd contrib/localnet/ && ./scripts/import-data.sh mainnet && $(DOCKER) compose --profile upgrade -f docker-compose.yml -f docker-compose-upgrade.yml up -d

###############################################################################
###                              Monitoring                                 ###
###############################################################################

start-monitoring:
	@echo "Starting monitoring services"
	cd contrib/localnet/grafana/ && ./get-tss-address.sh
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-monitoring.yml up -d

stop-monitoring:
	@echo "Stopping monitoring services"
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-monitoring.yml down

###############################################################################
###                                GoReleaser  		                        ###
###############################################################################

PACKAGE_NAME          := github.com/zeta-chain/node
GOLANG_CROSS_VERSION  ?= v1.20.7
GOPATH ?= '$(HOME)/go'
release-dry-run:
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip-validate --skip-publish --snapshot

release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-e "GITHUB_TOKEN=${GITHUB_TOKEN}" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip-validate

###############################################################################
###                     Local Mainnet Development                           ###
###############################################################################

#BTC
start-bitcoin-node-mainnet:
	cd contrib/rpc/bitcoind-mainnet && DOCKER_TAG=$(DOCKER_TAG) docker-compose up

stop-bitcoin-node-mainnet:
	cd contrib/rpc/bitcoind-mainnet && DOCKER_TAG=$(DOCKER_TAG) docker-compose down

clean-bitcoin-node-mainnet:
	cd contrib/rpc/bitcoind-mainnet && DOCKER_TAG=$(DOCKER_TAG) docker-compose down -v

#ETHEREUM
start-eth-node-mainnet:
	cd contrib/rpc/ethereum && DOCKER_TAG=$(DOCKER_TAG) docker-compose up

stop-eth-node-mainnet:
	cd contrib/rpc/ethereum && DOCKER_TAG=$(DOCKER_TAG) docker-compose down

clean-eth-node-mainnet:
	cd contrib/rpc/ethereum && DOCKER_TAG=$(DOCKER_TAG) docker-compose down -v

#ZETA

#FULL-NODE-RPC-FROM-BUILT-IMAGE
start-zetacored-rpc-mainnet:
	cd contrib/rpc/zetacored && bash init_docker_compose.sh mainnet image $(DOCKER_TAG)

stop-zetacored-rpc-mainnet:
	cd contrib/rpc/zetacored && bash kill_docker_compose.sh mainnet false

clean-zetacored-rpc-mainnet:
	cd contrib/rpc/zetacored && bash kill_docker_compose.sh mainnet true

#FULL-NODE-RPC-FROM-BUILT-IMAGE
start-zetacored-rpc-testnet:
	cd contrib/rpc/zetacored && bash init_docker_compose.sh athens3 image $(DOCKER_TAG)

stop-zetacored-rpc-testnet:
	cd contrib/rpc/zetacored && bash kill_docker_compose.sh athens3 false

clean-zetacored-rpc-testnet:
	cd contrib/rpc/zetacored && bash kill_docker_compose.sh athens3 true

#FULL-NODE-RPC-FROM-LOCAL-BUILD
start-zetacored-rpc-mainnet-localbuild:
	cd contrib/rpc/zetacored && bash init_docker_compose.sh mainnet localbuild $(DOCKER_TAG)

stop-zetacored-rpc-mainnet-localbuild:
	cd contrib/rpc/zetacored && bash kill_docker_compose.sh mainnet false

clean-zetacored-rpc-mainnet-localbuild:
	cd contrib/rpc/zetacored && bash kill_docker_compose.sh mainnet true

#FULL-NODE-RPC-FROM-LOCAL-BUILD
start-zetacored-rpc-testnet-localbuild:
	cd contrib/rpc/zetacored && bash init_docker_compose.sh athens3 localbuild $(DOCKER_TAG)

stop-zetacored-rpc-testnet-localbuild:
	cd contrib/rpc/zetacored && bash kill_docker_compose.sh athens3 false

clean-zetacored-rpc-testnet-localbuild:
	cd contrib/rpc/zetacored && bash kill_docker_compose.sh athens3 true


###############################################################################
###                               Debug Tools                               ###
###############################################################################

filter-missed-btc: install-zetatool
	zetatool filterdeposit btc --config ./tool/filter_missed_deposits/zetatool_config.json

filter-missed-eth: install-zetatool
	zetatool filterdeposit eth \
		--config ./tool/filter_missed_deposits/zetatool_config.json \
		--evm-max-range 1000 \
		--evm-start-block 19464041