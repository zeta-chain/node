#!/bin/bash

SMOKETEST_CMD=$1

echo "waiting for geth RPC to start..."
sleep 6
echo "funding deployer address 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", value: web3.toWei(100,"ether")})' attach http://eth:8545
echo "funding TSS address 0xF421292cb0d3c97b90EEEADfcD660B893592c6A2 with 1 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xF421292cb0d3c97b90EEEADfcD660B893592c6A2", value: web3.toWei(100,"ether")})' attach http://eth:8545

echo "waiting for 6s for the transaction to be mined"
sleep 6
echo "the new balance of the deployer addrees:"
curl -sS http://eth:8545 \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{"method":"eth_getBalance","params":["0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", "latest"],"id":1,"jsonrpc":"2.0"}'
curl -sS http://eth:8545 \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{"method":"eth_getBalance","params":["0xF421292cb0d3c97b90EEEADfcD660B893592c6A2", "latest"],"id":1,"jsonrpc":"2.0"}'
echo "running smoketest..."
smoketest "$SMOKETEST_CMD"
SMOKETEST_EXIT_CODE=$?
if [ $SMOKETEST_EXIT_CODE -eq 0 ]; then
  echo "smoketest passed"
  exit 0
else
  echo "smoketest failed"
  exit 1
fi
