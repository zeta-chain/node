#!/bin/bash

DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

## Must Compile Solidity Contract before running this script
npx hardhat compile

## Runs each deployment script and pass the network name as an argument
## Each deployment script will save the contract addresses into the ../addresses folder
npx hardhat run scripts/single-chain/deploy_zeta.js --network ethLocalNet
npx hardhat run scripts/single-chain/deploy_zeta.js --network bscLocalNet
npx hardhat run scripts/single-chain/deploy_zeta.js --network polygonLocalNet

# Load the contract addresses from the addresses folder
# The Network names in your hardhat config must match the network names in the addresses folder

eth_zetaMPI=$(cat ../../localnet-addresses/ethLocalNet-zetaMPI-address)
eth_zeta=$(cat ../../localnet-addresses/ethLocalNet-zeta-address)
bsc_zetaMPI=$(cat ../../localnet-addresses/bscLocalNet-zetaMPI-address) 
bsc_zeta=$(cat ../../localnet-addresses/bscLocalNet-zeta-address)
polygon_zetaMPI=$(cat ../../localnet-addresses/polygonLocalNet-zetaMPI-address)
polygon_zeta=$(cat ../../localnet-addresses/polygonLocalNet-zeta-address)

echo "-----------------------------------------------"
echo ""
echo "Contracts Deployed"
echo ""
echo "---Ethererum---"
echo "ZetaMPI Address ${eth_zetaMPI}"
echo "Zeta Coin Address ${eth_zeta}"
echo ""
echo "---Binance Smart Chain (BSC)---"
echo "ZetaMPI Address ${bsc_zetaMPI}"
echo "Zeta Coin Address ${bsc_zeta}"
echo ""
echo "---Polygon---"
echo "ZetaMPI Address ${polygon_zetaMPI}"
echo "Zeta Coin Address ${polygon_zeta}"



