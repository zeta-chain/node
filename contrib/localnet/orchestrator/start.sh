#!/bin/bash

SMOKETEST_CMD=$1

echo "waiting for geth RPC to start..."
sleep 6

# unlock the deployer account
echo "funding deployer address 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock erc20 tester accounts
echo "funding deployer address 0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6 with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock zeta tester accounts
echo "funding deployer address 0x5cC2fBb200A929B372e3016F1925DcF988E081fd with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x5cC2fBb200A929B372e3016F1925DcF988E081fd", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock the TSS account
echo "funding TSS address 0xF421292cb0d3c97b90EEEADfcD660B893592c6A2 with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xF421292cb0d3c97b90EEEADfcD660B893592c6A2", value: web3.toWei(100,"ether")})' attach http://eth:8545

# wait for the transaction to be mined
echo "waiting for 6s for the transaction to be mined"
sleep 6

# note: uncomment the following lines to print the balance of the deployer address if debugging is needed
#echo "the new balance of the deployer address:"
#curl -sS http://eth:8545 \
#  -X POST \
#  -H "Content-Type: application/json" \
#  --data '{"method":"eth_getBalance","params":["0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", "latest"],"id":1,"jsonrpc":"2.0"}'

# run smoketest
echo "running smoketest..."
smoketest "$SMOKETEST_CMD"
SMOKETEST_EXIT_CODE=$?

# if smoketest passed, exit with 0, otherwise exit with 1
if [ $SMOKETEST_EXIT_CODE -eq 0 ]; then
  echo "smoketest passed"
  exit 0
else
  echo "smoketest failed"
  exit 1
fi
