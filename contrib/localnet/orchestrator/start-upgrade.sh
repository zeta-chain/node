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


echo "running E2E command to setup the networks and populate the state..."

zetae2e "$ZETAE2E_CMD" --config-out deployed.yml

ZETAE2E_EXIT_CODE=$?
if [ $ZETAE2E_EXIT_CODE -ne 0 ]; then
  echo "E2E setup failed"
  exit 1
fi

echo "E2E setup passed, waiting for upgrade height..."

# Restart zetaclients at upgrade height
/work/restart-zetaclientd.sh -u 180 -n 2

echo "waiting 10 seconds for node to restart..."

sleep 10

echo "running E2E command to test the network after upgrade..."

zetae2e "$ZETAE2E_CMD" --skip-setup --config deployed.yml

ZETAE2E_EXIT_CODE=$?
if [ $ZETAE2E_EXIT_CODE -eq 0 ]; then
  echo "E2E passed after upgrade"
  exit 0
else
  echo "E2E failed after upgrade"
  exit 1
fi