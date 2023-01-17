# ZetaChain Local Net Development & Testing Environment
This directory contains a set of scripts to help you quickly set up a 
local ZetaChain network for development and testing purposes. 
The scripts are based on the Docker 
and [Docker Compose](https://docs.docker.com/compose/).

This is primarily tested on a recent Linux distribution such
as Ubuntu 22.04 LTS. 

The docker-compose.yml file defines a network with:

* 2 zetacore nodes
* 2 zetaclient nodes
* 1 go-ethereum private net node
* 1 bitcoin core private node

## Prerequisites
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/doc/install)
- [jq](https://stedolan.github.io/jq/download/)

## Steps

### Build zetanode & smoketest docker image
```bash
# in zeta-node/
$ docker build -t zetanode .
$ docker build -t smoketest -f Dockerfile.smoketest .
```

### Smoke Test Dev & Test Cyccle
The smoke test is in the directory /contrib/localnet/orchestrator/smoketest. 
It's a Go program that performs various operations on the localnet.

When you update/add tests to the smoke test, you need to rebuild the smoketest
image: 

```bash
# in zeta-node/
$ docker build -t smoketest -f Dockerfile.smoketest .
```

and then rebuild the orchestrator image:

```bash
# in zeta-node/contrib/localnet/orchestrator
$ docker build -t orchestrator .
```




## References
[Setup testnet reference](https://www.notion.so/zetachain/Set-up-athens-1-like-testnet-to-test-your-PRs-ac523eb5dd5d4e73902072ab7d85fa2f)

