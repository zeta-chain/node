#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

docker-compose -f docker-compose.simple.yml up -d

# echo "Go to http://localhost:3010/ verify all nodes gradually connected to each other and block start counting."
# echo "It will take a few minutes for the miner to start"

echo "RPC Port: $RPC_PORT"