#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR || exit

git clone https://github.com/bnb-chain/bsc-docker
cp env_vars .env
cp .env bsc-docker/.env

cd bsc-docker && git pull

docker-compose -f docker-compose.bsc.yml build
