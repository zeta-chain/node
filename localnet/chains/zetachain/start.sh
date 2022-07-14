#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd "$DIR" || exit

ZETA_MONOREPO_PATH="$DIR/../../zetachain-monorepo/"
CONTRACT_ADDRESS_FILE=$ZETA_MONOREPO_PATH/packages/addresses/src/addresses/addresses.troy.json

VARS_CONFIG_FILE=../../.env
echo "Sourcing Environment Variables from $VARS_CONFIG_FILE"
source $VARS_CONFIG_FILE

ETH_CONNECTOR_ADDRESS=$(jq -r .'"eth-localnet"'.connector "$CONTRACT_ADDRESS_FILE")
BSC_CONNECTOR_ADDRESS=$(jq -r .'"bsc-localnet"'.connector "$CONTRACT_ADDRESS_FILE")
POLYGON_CONNECTOR_ADDRESS=$(jq -r .'"polygon-localnet"'.connector "$CONTRACT_ADDRESS_FILE")

echo "Adding Contract Addresses & Endpoints to $(pwd)/.env file"
cp env_vars .env
echo "ETH_CONNECTOR_ADDRESS=$ETH_CONNECTOR_ADDRESS" >> .env
echo "BSC_CONNECTOR_ADDRESS=$BSC_CONNECTOR_ADDRESS" >> .env
echo "POLYGON_CONNECTOR_ADDRESS=$POLYGON_CONNECTOR_ADDRESS" >> .env

if [ "$USE_GANACHE" == true ]; then
    echo "ETH_ENDPOINT=http://ganache-eth:8545" >> .env
    echo "BSC_ENDPOINT=http://ganache-bsc:8545" >> .env
    echo "POLYGON_ENDPOINT=http://ganache-polygon:8545" >> .env
else 
    echo "ETH_ENDPOINT=http://ethereum-geth-rpc-endpoint-1.localnet:8545" >> .env
    echo "BSC_ENDPOINT=http://bsc.localnet:8545" >> .env
    echo "POLYGON_ENDPOINT=http://bor.localnet:8545" >> .env
fi

docker compose up -d

# Get TSS Address
echo "Waiting for TSS Address... This may take a few minutes"
until  [ ! -z "$TSS_ADDR" ]
do
    RESPONSE=$(curl -s http://localhost:1317/zeta-chain/zetacore/TSS | jq -r '.TSS[0]'.address)
    CHARACTER_COUNT=$(echo "${RESPONSE}" | wc -m | xargs)
    # echo "$RESPONSE"
    # Uses Character Count to determine if returned value is an address or not
    if [ "$CHARACTER_COUNT" = 43 ] && [ "$RESPONSE" != "0x0000000000000000000000000000000000000000" ]; then 
        TSS_ADDR=$RESPONSE
        echo "TSS Address Is: ${TSS_ADDR}"
        echo "Ending Loop"
        break
    fi
    echo "Waiting for TSS Address..."
    sleep 10
done
echo "TSS Address Is: ${TSS_ADDR}"

# Save TSS Address
echo "Saving TSS Address to address file"

jq --arg tss "$TSS_ADDR" .'"eth-localnet".tss = $tss' "$CONTRACT_ADDRESS_FILE" > tmp1.json
jq --arg tss "$TSS_ADDR" .'"bsc-localnet".tss = $tss' tmp1.json > tmp2.json
jq --arg tss "$TSS_ADDR" .'"polygon-localnet".tss = $tss' tmp2.json > tmp3.json
mv tmp3.json "$CONTRACT_ADDRESS_FILE"
rm tmp*.json

# Update TSS Address On Contracts
cd "$ZETA_MONOREPO_PATH"/packages/protocol-contracts/ || exit
npx hardhat run scripts/set-zeta-token-addresses.ts --network eth-localnet
npx hardhat run scripts/set-zeta-token-addresses.ts --network bsc-localnet
#npx hardhat run scripts/set-zeta-token-addresses.ts --network polygon-localnet

# Send Gas
npx hardhat run scripts/send-tss-gas.ts --network eth-localnet
npx hardhat run scripts/send-tss-gas.ts --network bsc-localnet
npx hardhat run scripts/send-tss-gas.ts --network polygon-localnet

# Approve Connector contract to spend Tokens
npx hardhat run scripts/token-approval.ts --network eth-localnet
npx hardhat run scripts/token-approval.ts --network bsc-localnet
npx hardhat run scripts/token-approval.ts --network polygon-localnet
