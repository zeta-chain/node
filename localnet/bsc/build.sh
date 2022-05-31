#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

git clone https://github.com/bnb-chain/bsc-docker
cp env_vars .env
cp .env bsc-docker/.env

cd bsc-docker && git pull

docker-compose -f docker-compose.bsc.yml build

# docker-compose -f docker-compose.simple.bootstrap.yml build
# docker-compose -f docker-compose.simple.yml build

# # Generate genesis.json, validators & bootstrap cluster data
# # Once finished, all cluster bootstrap data are generated at ./storage
# docker-compose -f docker-compose.simple.bootstrap.yml run bootstrap-simple

 