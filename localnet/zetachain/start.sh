#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

VARS_CONFIG_FILE=../vars.config
echo "Sourcing Environment Variables from $VARS_CONFIG_FILE"
source $VARS_CONFIG_FILE

ETH_MPI_ADDRESS=$(jq -r .'"eth-localnet"'.connector $CONTRACT_ADDRESS_FILE)
BSC_MPI_ADDRESS=$(jq -r .'"bsc-localnet"'.connector $CONTRACT_ADDRESS_FILE)
POLYGON_MPI_ADDRESS=$(jq -r .'"matic-localnet"'.connector $CONTRACT_ADDRESS_FILE)

echo "Adding Contract Addresses & Endpoints to $(pwd)/.env file"
cp env_vars .env
echo "ETH_MPI_ADDRESS=$ETH_MPI_ADDRESS" >> .env
echo "BSC_MPI_ADDRESS=$BSC_MPI_ADDRESS" >> .env
echo "POLYGON_MPI_ADDRESS=$POLYGON_MPI_ADDRESS" >> .env

if [ $USE_GANACHE == true ]; then
    echo "ETH_ENDPOINT=http://ganache-eth:8545" >> .env
    echo "BSC_ENDPOINT=http://ganache-bsc:8545" >> .env
    echo "POLYGON_ENDPOINT=http://ganache-polygon:8545" >> .env
else 
    echo "ETH_ENDPOINT=http://ethereum-geth-rpc-endpoint-1.localnet:8545" >> .env
    echo "BSC_ENDPOINT=http://bsc-rpc:8545" >> .env
    echo "POLYGON_ENDPOINT=http://bor.localnet:8545" >> .env
fi

docker compose up -d

# Get TSS Address
echo "Waiting for TSS Address... This may take a few minutes"
until  [ ! -z "$TSS_ADDR" ]
do
    RESPONSE=$(curl -s http://localhost:1317/zeta-chain/zetacore/TSS/ETH | jq -r .TSS.address)
    CHARACTER_COUNT=$(echo $RESPONSE | wc -m)
    # Uses Character Count to determine if returned value is an address or not
    if [ $CHARACTER_COUNT = 43  ]; then 
        TSS_ADDR=$RESPONSE
        echo "TSS Address Is: ${TSS_ADDR}"
        break
    fi
    echo "Waiting for TSS Address..."
    sleep 10
done
echo "TSS Address Is: ${TSS_ADDR}"

# Save TSS Address
echo "Saving TSS Address to address file"

jq --arg tss "$TSS_ADDR" .'"eth-localnet".tss = $tss' $CONTRACT_ADDRESS_FILE > tmp1.json
jq --arg tss "$TSS_ADDR" .'"bsc-localnet".tss = $tss' tmp1.json > tmp2.json
jq --arg tss "$TSS_ADDR" .'"matic-localnet".tss = $tss' tmp2.json > tmp3.json
mv tmp3.json $CONTRACT_ADDRESS_FILE
rm tmp*.json

# Update TSS Address On Contracts
cd $ZETA_CONTRACTS_PATH/packages/protocol-contracts/
npx hardhat run scripts/set-zeta-token-addresses.ts --network eth-localnet
npx hardhat run scripts/set-zeta-token-addresses.ts --network bsc-localnet
npx hardhat run scripts/set-zeta-token-addresses.ts --network polygon-localnet

# Send Gas
npx hardhat run scripts/send-tss-gas.ts --network eth-localnet
npx hardhat run scripts/send-tss-gas.ts --network bsc-localnet
npx hardhat run scripts/send-tss-gas.ts --network polygon-localnet


