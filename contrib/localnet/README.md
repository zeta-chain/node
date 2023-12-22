# ZetaChain Local Net Development & Testing Environment
This directory contains a set of scripts to help you quickly set up a 
local ZetaChain network for development and testing purposes. 
The scripts are based on the Docker 
and [Docker Compose](https://docs.docker.com/compose/).

As a smoke test (sanity integration tests), the setup aims
at fully automatic, only requiring a few image building steps
and docker compose launch. 

As a development testing environment, the setup aims to be
flexible, close to real world, and with fast turnaround
between edit code -> compile -> test results. 

This is primarily tested on a recent Linux distribution such
as Ubuntu 22.04 LTS, though macOS should also work (not tested). 

The docker-compose.yml file defines a network with:

* 2 zetacore nodes
* 2 zetaclient nodes
* 1 go-ethereum private net node (act as GOERLI testnet, with chainid 1337)
* 1 bitcoin core private node (planned; not yet done)
* 1 orchestrator node which coordinates smoke tests. 

## Prerequisites
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/doc/install)
- [jq](https://stedolan.github.io/jq/download/)

## Steps

### Build zetanode 
```bash
$ make zetanode
```

This Makefile rule builds the zetanode image. **Rebuild if zetacored/zetaclientd code is updated**.  
```bash
# in zeta-node/
$ docker build -t zetanode .
```

### Smoke Test Dev & Test Cycle
The smoke test is in the directory /contrib/localnet/orchestrator/smoketest. 
It's a Go program that performs various operations on the localnet.

The above `make zetanode` should already produced the orchestrator image.

### Run smoke test

Now we have built all the docker images; we can run the smoke test with make command:
```bash
# in zeta-node/
make start-smoketest
```
which does the following docker compose command:
```bash
# in zeta-node/contrib/localnet/orchestrator
$ docker compose up -d
```

The most straightforward log to observe is the orchestrator log.
If everything works fine, it should finish without panic, and with
a message "smoketest done". 

To stop the tests, 
```bash
# in zeta-node/
make stop-smoketest
```
which does the following docker compose command:
```bash
# in zeta-node/contrib/localnet/orchestrator
$ docker compose down --remove-orphans
```
### Run monitoring setup
Before starting the monitoring setup, make sure you have added the eth addresses of interest to the grafana/addresses.txt file
```bash
# in zeta-node/
make start-monitoring
```
which does the following docker compose command:
```bash
# in zeta-node/contrib/localnet/
$ docker compose -f docker-compose-monitoring.yml up -d
```
To spin down the monitoring setup:
```bash
# in zeta-node/
make stop-monitoring
```
which does the following docker compose command:
```bash
# in zeta-node/contrib/localnet/
$ docker compose -f docker-compose-monitoring.yml down --remove-orphans
```
### Run blockscout for localnet
```bash
# in zeta-node/
make start-blockscout
```
which does the following docker compose command:
```bash
# in zeta-node/contrib/localnet/blockscout
$ docker compose -f compose -f docker-compose.yml up -d
```
To spin down the blockscout setup:
```bash
# in zeta-node/
make stop-blockscout
```
which does the following docker compose command:
```bash
# in zeta-node/contrib/localnet/blockscout
$ docker compose -f docker-compose.yml down --remove-orphans
```

## Useful data

- On GOERLI (private ETH net), the deployer account is pre-funded with Ether. 
[Deployer Address and Private Key](orchestrator/smoketest/main.go)

- TSS Address (on ETH): 0xF421292cb0d3c97b90EEEADfcD660B893592c6A2



## Add more smoke tests
The smoke test (integration tests) are located in the
orchestrator/smoketest directory. The orchestrator is a Go program.


## References
[Setup testnet reference](https://www.notion.so/zetachain/Set-up-athens-1-like-testnet-to-test-your-PRs-ac523eb5dd5d4e73902072ab7d85fa2f)

