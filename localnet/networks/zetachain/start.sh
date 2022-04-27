#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

docker compose up -d

# Get TSS Address
echo "Waiting for TSS Address... This may take a few minutes"
until  [ ! -z "$TSS_ADDR" ]
do
    RESPONSE=$(curl -s http://localhost:1317/zeta-chain/zetacore/zetacore/tSS/ETH | jq -r .TSS.address)
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

cd ../../hardhat
ts-node scripts/update-localnet-tss-address.ts ${TSS_ADDR}
echo "TSS Address Is: ${TSS_ADDR}"
