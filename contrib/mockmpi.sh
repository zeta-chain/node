#!/bin/env bash

if [[ -z "$PRIVKEY" ]]; then
  echo "Must provide PRIVKEY in environment" 1>&2
  exit 1
fi

export BSC_ENDPOINT=wss://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/bsc/testnet/archive/ws
export ETH_ENDPOINT=wss://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/eth/goerli/archive/ws

while true; do
  go run ./cmd/mockmpi
  sleep 1
done
