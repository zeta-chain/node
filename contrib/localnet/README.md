# ZetaChain Local Net Development & Testing Environment
This directory contains a set of scripts to help you quickly set up a 
local ZetaChain network for development and testing purposes. 
The scripts are based on the Docker 
and [Docker Compose](https://docs.docker.com/compose/).

This is primarily tested on a recent Linux distribution such
as Ubuntu 22.04 LTS. 

The docker-compose.yml file defines a network with:

* 4 zetacore nodes
* 4 zetaclient nodes
* 1 go-ethereum private net node
* 1 bitcoin core private node

## Prerequisites
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/doc/install)
- [jq](https://stedolan.github.io/jq/download/)

## Create a local testnet of ZetaChain

### To setup the network nodes

### Create .pem file

```
ssh-keygen -P "" -t rsa -b 4096 -m pem -f z.pem
```

### Clone zeta-node repo

```
git clone git@github.com:zeta-chain/zeta-node.git
```


### Create docker image

```
docker build -t zetanode .
```

### Launch nodes

```
docker compose up
```

### Terminals

At this point you should launch one terminal per node, can be [tmux](https://github.com/tmux/tmux/wiki), or whatever you prefer.
Available nodes are `node0`,`node1`,`node2` and `node3`.
All the scripts into the containers are prepared to detect in which node are
executing, so for example if you launch the reset script on all nodes, it will only work on node0. This is useful on tmux because you can safely run the commands below in order on all the terminals at the same time.


### Connect to nodes 

```
./node.sh <Node number>
```

E.g.
```
./node.sh 0
./node.sh 1
./node.sh 2
./node.sh 3
```


### Reset (node0)

```
./reset-testnet.sh 
```

### Start zetacored

```
./start.sh (all nodes)
```

## Zetaclientd

### Seed tss (all nodes)
```
./seed.sh
```

### Keygen (all nodes)

You need to pass a block number in the future as a parameter.

Usually 10-15 blocks are enough ( depends on how fast you're to launch the next step on all the virtuals )

First check the current height on zetacored:

In the screenshot current height is 1660. 
![height](docs/height.png)


e.g. using 1700. ( wait until it reachs that block number)

```
./keygen.sh 1700
```

At the end of this step you will have a TSS address, remember to write it down somewhere, because you will need it when deploying contracts.

You should run keygen only one time per testnet initialization.

### Env vars

Customize env vars as you like...

```
./env.sh
```

### Launch clients

```
./start_client.sh
```

### Notes

* From now you can stop the services on each container and start again with :

`start.sh` for zetacored
`start_client.sh` for zetaclientd

* If you want to persist the configuration you can backup `/home/alpine` folder on each node.


### Now you need to deploy contracts 
[Setup testnet reference](https://www.notion.so/zetachain/Set-up-athens-1-like-testnet-to-test-your-PRs-ac523eb5dd5d4e73902072ab7d85fa2f)

