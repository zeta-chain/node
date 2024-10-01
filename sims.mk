#!/usr/bin/make -f

########################################
### Simulations

BINDIR ?= $(GOPATH)/bin
SIMAPP = ./app

runsim: $(BINDIR)/runsim
$(BINDIR)/runsim:
	@echo "Installing runsim..."
	@(cd /tmp && go install github.com/cosmos/tools/cmd/runsim@v1.0.0)


test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=10 -BlockSize=20 -Commit=true -Period=0 -v -timeout 24h


test-sim-fullappsimulation:
	@echo "Running TestFullAppSimulation."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Enabled=true \
		-NumBlocks=10 -BlockSize=20 -Commit=true -Period=0 -v -timeout 24h
test-sim-multi-seed-long: runsim
	@echo "Running long multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 500 50 TestFullAppSimulation

test-sim-multi-seed-short: runsim
	@echo "Running short multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 10 TestFullAppSimulation



.PHONY: \
test-sim-nondeterminism \
test-sim-fullappsimulation \
test-sim-multi-seed-long \
test-sim-multi-seed-short

