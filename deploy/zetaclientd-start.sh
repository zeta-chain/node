#!/bin/bash

echo "$1" "$2"

NODE_NUMBER=$1
SEED_NODE=$2

echo "Starting ZetaClient Node $NODE_NUMBER"
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:~/go/bin
export IDX=$NODE_NUMBER
export TSSPATH=~/.tssnew

if [ -z "${MYIP}" ]; then
    # If MYIP is not set, use the private IP of the host
    echo "MYIP ENV Variable Not Set -- Setting it automatically using host IP"
    export MYIP=$(hostname -i)
fi
echo "MYIP: $MYIP"

rm -f ~/.tssnew/address_book.seed || true

if (($NODE_NUMBER == 0)); then
    sleep 5 # Wait for Zetacored to start
    exec zetaclientd -val val -enable-chains GOERLI,BSCTESTNET,MUMBAI,ROPSTEN,BAOBAB
        
else
    until [ -f SEED_NODE_ID ]; do
        echo "Waiting for Seed Node Validator ID"
        sleep 10
        curl -s "${SEED_NODE}":8123/p2p -o SEED_NODE_ID
    done
    SEED_NODE_ID=$(cat SEED_NODE_ID)
    echo "SEED_NODE_ID=${SEED_NODE_ID}"

    exec zetaclientd -val val \
        --peer /dns/"${SEED_NODE}"/tcp/6668/p2p/"${SEED_NODE_ID}" \
        -enable-chains GOERLI,BSCTESTNET,MUMBAI,ROPSTEN,BAOBAB
fi
