#!/bin/bash

echo "funding deployer address 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC with 100ETH"
geth attach http://eth:8545 --exec 'loadScript("/cmd/setup.js")'
echo "waiting for 10s for the transaction to be mined"
sleep 5
echo "the new balance of the deployer addrees:"
curl http://eth:8545 \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{"method":"eth_getBalance","params":["0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", "latest"],"id":1,"jsonrpc":"2.0"}'

echo "running smoketest..."
smoketest
echo "smoketest done"