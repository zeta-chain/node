#!/bin/bash

ZETAE2E_CMD=$1

echo "waiting for geth RPC to start..."
sleep 2

# unlock the deployer account
echo "funding deployer address 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock erc20 tester accounts
echo "funding deployer address 0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6 with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock zeta tester accounts
echo "funding deployer address 0x5cC2fBb200A929B372e3016F1925DcF988E081fd with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x5cC2fBb200A929B372e3016F1925DcF988E081fd", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock bitcoin tester accounts
echo "funding deployer address 0x283d810090EdF4043E75247eAeBcE848806237fD with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x283d810090EdF4043E75247eAeBcE848806237fD", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock ethers tester accounts
echo "funding deployer address 0x8D47Db7390AC4D3D449Cc20D799ce4748F97619A with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x8D47Db7390AC4D3D449Cc20D799ce4748F97619A", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock miscellaneous tests accounts
echo "funding deployer address 0x90126d02E41c9eB2a10cfc43aAb3BD3460523Cdf with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0x90126d02E41c9eB2a10cfc43aAb3BD3460523Cdf", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock advanced erc20 tests accounts
echo "funding deployer address 0xcC8487562AAc220ea4406196Ee902C7c076966af with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xcC8487562AAc220ea4406196Ee902C7c076966af", value: web3.toWei(100,"ether")})' attach http://eth:8545

# unlock the TSS account
echo "funding TSS address 0xF421292cb0d3c97b90EEEADfcD660B893592c6A2 with 100 Ether"
geth --exec 'eth.sendTransaction({from: eth.coinbase, to: "0xF421292cb0d3c97b90EEEADfcD660B893592c6A2", value: web3.toWei(100,"ether")})' attach http://eth:8545

# run e2e tests
echo "running e2e tests..."
zetae2e "$ZETAE2E_CMD" --setup-only --config-out deployed.yml
ZETAE2E_EXIT_CODE=$?

# if e2e passed, exit with 0, otherwise exit with 1
if [ $ZETAE2E_EXIT_CODE -eq 0 ]; then
  echo "starting tests with skip setup..."

  zetae2e "$ZETAE2E_CMD" --skip-setup --config deployed.yml

  echo "e2e passed"
  exit 0
else
  echo "e2e failed"
  exit 1
fi
