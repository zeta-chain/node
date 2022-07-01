#!/bin/bash

LOCALNET_DIR="$( cd "$( dirname "$0" )" && pwd )"
cd "$LOCALNET_DIR" || exit

VARS_CONFIG_FILE=.env
echo "Sourcing Environment Variables from $VARS_CONFIG_FILE"
source $VARS_CONFIG_FILE

docker network create localnet --subnet 172.24.0.0/16

# Deploy External Nodes
if [ "$USE_GANACHE" == true ]; then
    echo "Launching Ganache Development Networks (Not Forked)"
        cd ganache || exit
        ./start.sh
        cd ..
        sleep 10
else
    for d in $(ls -d */); do 
        if  [ "$d" != "zetachain/" ] && [ "$d" != "ganache/"  ] && [ "$d" != "node_modules/" ]; then
            echo "Starting $d"
            cd "$d" || exit
            ./start.sh
            cd ..
        fi
    done
    echo "Pausing for LocalNet Nodes to start -- Please Wait... (20s)"
    sleep 20
fi

# Deploy Contracts
if [ "$DEPLOY_CONTRACTS" == true ]; then
    echo "Deploying ZetaChain Contracts"

    cd "$ZETA_CONTRACTS_PATH"/packages/protocol-contracts/ || exit
    npx hardhat run scripts/deploy.ts --network eth-localnet
    npx hardhat run scripts/deploy.ts --network bsc-localnet
    npx hardhat run scripts/deploy.ts --network polygon-localnet
fi

# Deploy ZetaChain Nodes
cd "$LOCALNET_DIR" || exit

cd zetachain || exit
echo "Launching Zetachain Nodes"
./start.sh
cd ..

## Output Results
for d in $(ls -d */); do 
    cd "$d" || exit
    source .env
    echo ""
    echo "----------"
    echo "Network Name: ${NETWORK_NAME}"
    echo "RPC Port: ${RPC_PORT} "
    echo "Chain Id: ${NETWORK_ID}"
    cd ..
done
