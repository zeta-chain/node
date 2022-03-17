#!/bin/env bash

if [[ -z "$PRIVKEY" ]]; then
  echo "PRIVKEY not provided. Using Test Key"
  export PRIVKEY="2082bc9775d6ee5a05ef221a9d1c00b3cc3ecb274a4317acc0a182bc1e05d1bb"
fi

export BSC_ENDPOINT=wss://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/bsc/testnet/archive/ws
export ETH_ENDPOINT=wss://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/eth/goerli/archive/ws
export POLYGON_ENDPOINT=wss://speedy-nodes-nyc.moralis.io/9555dbf7bdae477b335c2a5d/polygon/mumbai/archive/ws

while true; do
  go run ./cmd/mockmpi "$@"
  sleep 1
  echo -e "\nRestarting..."
done
