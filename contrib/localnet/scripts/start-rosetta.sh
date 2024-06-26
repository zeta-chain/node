#!/bin/bash

# This script is used to start the Rosetta API server for the Zetacore network.

echo "Waiting for network to start producing blocks"
CURRENT_HEIGHT=0
WAIT_HEIGHT=1
while [[ $CURRENT_HEIGHT -lt $WAIT_HEIGHT ]]
do
    CURRENT_HEIGHT=$(curl -s zetacore0:26657/status | jq '.result.sync_info.latest_block_height' | tr -d '"')
    sleep 5
done

zetacored rosetta --tendermint zetacore0:26657 --grpc zetacore0:9090 --network athens_101-1 --blockchain zetacore