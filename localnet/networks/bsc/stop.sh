#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

source .env
# cd bsc-docker
docker-compose -f docker-compose.simple.yml stop
docker container prune -f