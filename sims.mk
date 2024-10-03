#!/usr/bin/make -f

########################################
### Simulations

BINDIR ?= $(GOPATH)/bin
SIMAPP = ./tests/simulation


# Run sim is a cosmos tool which helps us to run multiple simulations in parallel.
runsim: $(BINDIR)/runsim
$(BINDIR)/runsim:
	@echo 'Installing runsim...'
	@TEMP_DIR=$$(mktemp -d) && \
	cd $$TEMP_DIR && \
	go install github.com/cosmos/tools/cmd/runsim@v1.0.0 && \
	rm -rf $$TEMP_DIR || (echo 'Failed to install runsim' && exit 1)
	@echo 'runsim installed successfully'


define run-sim-test
	@echo "Running $(1)..."
	@go test -mod=readonly $(SIMAPP) -run $(2) -Enabled=true \
		-NumBlocks=$(3) -BlockSize=$(4) -Commit=true -Period=0 -v -timeout $(5)
endef

test-sim-nondeterminism:
	$(call run-sim-test,"non-determinism test",TestAppStateDeterminism,100,200,2h)

test-sim-fullappsimulation:
	$(call run-sim-test,"TestFullAppSimulation",TestFullAppSimulation,100,200,2h)

test-sim-multi-seed-long: runsim
	@echo "Running long multi-seed application simulation."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 500 50 TestFullAppSimulation

test-sim-multi-seed-short: runsim
	@echo "Running short multi-seed application simulation."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 10 TestFullAppSimulation



.PHONY: \
test-sim-nondeterminism \
test-sim-fullappsimulation \
test-sim-multi-seed-long \
test-sim-multi-seed-short

