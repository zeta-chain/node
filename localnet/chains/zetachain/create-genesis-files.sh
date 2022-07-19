#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

rm -rf ./config/genesis/*  ./config/*/data/* ./config/*/config/* ./config/*/keyring-test/* || true

mkdir -p ./config/genesis
mkdir -p ./config/node{0,1,2,3}/{data,keyring-test}
mkdir -p ./config/node{0,1,2,3}/config/gentx

docker compose -f docker-compose.build-config.yml up && docker compose -f docker-compose.build-config.yml rm -fsv
