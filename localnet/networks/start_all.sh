#!/bin/bash

LOCALNET_DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $LOCALNET_DIR

docker network create localnet --subnet 172.24.0.0/16

for d in $(ls -d */); do 
    if  [ $d != "zetachain/" ]; then
        echo "Starting $d"
        cd $d
        ./start.sh
        cd ..
    fi
done

echo "Pausing for LocalNet Nodes to start -- Please Wait..."
sleep 30

echo "Deploying ZetaChain Contracts"
../hardhat/scripts/multi-chain/deploy_multi_chain.sh

cd $LOCALNET_DIR

ETH_MPI_ADDRESS=$(cat ../hardhat/localnet-addresses/ethLocalNet-zetaMPI-address)
# ETH_ZETA=$(cat ../hardhat/localnet-addresses/ethLocalNet-zeta-address) # Commented out because is isn't needed yet
BSC_MPI_ADDRESS=$(cat ../hardhat/localnet-addresses/bscLocalNet-zetaMPI-address) 
# BSC_ZETA=$(cat ../hardhat/localnet-addresses/bscLocalNet-zeta-address)
POLYGON_MPI_ADDRESS=$(cat ../hardhat/localnet-addresses/polygonLocalNet-zetaMPI-address)
# POLYGON_ZETA=$(cat ../hardhat/localnet-addresses/polygonLocalNet-zeta-address)

cd zetachain 
cp env_vars .env
echo "Added Contract Addresses to ZetaClient Environment Variables"
echo "ETH_MPI_ADDRESS=$ETH_MPI_ADDRESS" >> .env
echo "BSC_MPI_ADDRESS=$BSC_MPI_ADDRESS" >> .env
echo "POLYGON_MPI_ADDRESS=$POLYGON_MPI_ADDRESS" >> .env
echo "Launching Zetachain Nodes"
./start.sh
cd ..


for d in $(ls -d */); do 
    cd $d
    source .env
    echo ""
    echo "----------"
    echo "Network Name: ${NETWORK_NAME}"
    echo "RPC Port: ${RPC_PORT} "
    echo "Chain Id: ${NETWORK_ID}"
    cd ..
done

echo "WARNING - BSC LocalNet is known to have issues and may not be working properly at this time"