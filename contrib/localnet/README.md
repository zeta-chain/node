# ZetaChain Local Net Development & Testing Environment
This directory contains a set of scripts to help you quickly set up a 
local ZetaChain network for development and testing purposes. 
The scripts are based on the Docker 
and Docker Compose.

As a development testing environment, the setup aims to be
flexible, close to real world, and with fast turnaround
between edit code -> compile -> test results.

The `docker-compose.yml` file defines a network with:

* 2 zetacore nodes
* 2 zetaclient nodes
* 1 go-ethereum private net node (act as GOERLI testnet, with chainid 1337)
* 1 bitcoin core private node (planned; not yet done)
* 1 orchestrator node which coordinates E2E tests. 

The following Docker compose files can extend the default localnet setup:

- `docker-compose-stresstest.yml`: Spin up more nodes and clients for testing performance of the network.
- `docker-compose-upgrade.yml`: Spin up a network with a upgrade proposal defined at a specific height.

Finally, `docker-compose-monitoring.yml` can be run separately to spin up a local grafana and prometheus setup to monitor the network.

## Running Localnet

Running the localnet requires `zetanode` Docker image. The following command should be run at the root of the repo:

```
make zetanode
```

Localnet can be started with Docker Compose:

```
docker-compose up -d
```

To stop the localnet:

```
docker-compose down
```

## Orchestrator

The `orchestrator` directory contains the orchestrator node which coordinates E2E tests. The orchestrator is responsible for:

- Initializing accounts on the local Ethereum network.
- Using `zetae2e` CLI to run the tests.
- Restarting ZetaClient during upgrade tests.

## Scripts

The `scripts` directory mainly contains the following scripts:

- `start-zetacored.sh`: Used by zetacore images to bootstrap genesis and start the nodes.
- `start-zetaclientd.sh`: Used by zetaclient images to setup TSS and start the clients.

## Prerequisites

The following are required to run the localnet:

- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/doc/install)
- [jq](https://stedolan.github.io/jq/download/)
