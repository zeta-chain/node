#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd "$DIR" || exit

git clone https://github.com/bnb-chain/bsc-docker
cp env_vars .env
cp .env bsc-docker/.env

cd bsc-docker && git pull

# If using a Mac, force docker to use amd64 images
OS=$(uname -s)
if [ "$OS" = "Darwin" ]; then
    sed -i '' -e '1 s:^FROM golang:FROM --platform=linux/amd64 golang:' Dockerfile.bsc
    sed -i '' -e '1 s:^FROM ethereum/solc:FROM --platform=linux/amd64 ethereum/solc:' Dockerfile.bootstrap
fi

docker-compose -f docker-compose.bsc.yml build

## Not Currently Used
# docker-compose -f docker-compose.simple.bootstrap.yml build
# docker-compose -f docker-compose.simple.yml build

# # Generate genesis.json, validators & bootstrap cluster data
# # Once finished, all cluster bootstrap data are generated at ./storage
# docker-compose -f docker-compose.simple.bootstrap.yml run bootstrap-simple

 
