#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

cp env_vars .env
docker build . -t geth-client --build-arg ACCOUNT_PASSWORD=${ACCOUNT_PASSWORD}