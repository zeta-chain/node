.PHONY: build

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

BUILD_FLAGS := -ldflags '$(ldflags)' -tags pebbledb,ledger

TEST_DIR?="./..."
TEST_BUILD_FLAGS := -tags pebbledb,ledger
HSM_BUILD_FLAGS := -tags pebbledb,ledger,hsm_test

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

lint-cosmos-gosec:
	@bash ./scripts/cosmos-gosec.sh

gosec:
	gosec  -exclude-dir=localnet ./...

###############################################################################
###                           Generation commands  		                    ###
###############################################################################

proto:
	@echo "--> Removing old Go types "
	@find . -name '*.pb.go' -type f -delete
	@echo "--> Generating new Go types from protocol buffer files"
	@bash ./scripts/protoc-gen-go.sh
	@buf format -w
.PHONY: proto

typescript:
	@echo "--> Generating TypeScript bindings"
	@bash ./scripts/protoc-gen-typescript.sh
.PHONY: typescript

proto-format:
	@bash ./scripts/proto-format.sh

openapi:
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

generate: proto openapi specs typescript docs-zetacored
.PHONY: generate

###############################################################################
###                            E2E tests                                    ###
###############################################################################

zetanode:
	@echo "Building zetanode"
	$(DOCKER) build -t zetanode -f ./Dockerfile .
	$(DOCKER) build -t orchestrator -f contrib/localnet/orchestrator/Dockerfile.fastbuild .
.PHONY: zetanode

install-zetae2e: go.sum
	@echo "--> Installing zetae2e"
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/zetae2e
.PHONY: install-zetae2e

start-e2etest:
	@echo "--> Starting e2e test"
	cd contrib/localnet/ && $(DOCKER) compose up -d

stop-e2etest:
	@echo "--> Stopping e2e test"
	cd contrib/localnet/ && $(DOCKER) compose down --remove-orphans

start-e2etest-p2p-diag:
	@echo "--> Starting e2e test in p2p diagnostic mode"
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-p2p-diag.yml up -d

stop-e2etest-p2p-diag:
	@echo "--> Stopping e2e test in p2p diagnostic mode"
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-p2p-diag.yml down --remove-orphans

stress-test: zetanode
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-stresstest.yml up -d

stop-stress-test:
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-stresstest.yml down --remove-orphans

start-upgrade-test:
	@echo "--> Starting upgrade test"
	$(DOCKER) build --build-arg old_version=v12.2.1 -t zetanode -f ./Dockerfile-upgrade .
	$(DOCKER) build -t orchestrator -f contrib/localnet/orchestrator/Dockerfile.fastbuild .
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-upgrade.yml up -d

stop-upgrade-test:
	cd contrib/localnet/ && $(DOCKER) compose -f docker-compose-stateful.yml down --remove-orphans

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
GOLANG_CROSS_VERSION  ?= v1.20
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

mainnet-zetarpc-node:
	cd contrib/local-mainnet/zetacored && docker-compose up

mainnet-bitcoind-node:
	cd contrib/local-mainnet/bitcoind && docker-compose up