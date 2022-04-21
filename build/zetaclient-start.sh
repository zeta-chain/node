#!/bin/bash

echo $1 $2

NODE_NUMBER=$1
SEED_NODE=$2

echo "Starting ZetaClient Node $NODE_NUMBER"
source /etc/environment
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin
export IDX=$NODE_NUMBER
export TSSPATH=/root/.tssnew


if [ -z ${MYIP} ]; then
    # If MYIP is not set, use the private IP of the host
    export MYIP=$(hostname -i)
fi

MYIP=$(hostname -i)
echo "Hostname: $(hostname)"
echo "MYIP: $MYIP"

if (($NODE_NUMBER == 0)); then
    sleep 5 # Wait for Zetacored to start
    yes | zetaclientd -val val 2>&1 | tee ~/.zetaclient/zetaclient.log
else
    until [ -f SEED_NODE_ID ]; do
        echo "Waiting for Seed Node Validator ID"
        sleep 10
        curl -s ${SEED_NODE}:8123/p2p -o SEED_NODE_ID

    done
    SEED_NODE_ID=$(cat SEED_NODE_ID)
    echo "SEED_NODE_ID=${SEED_NODE_ID}"

    yes | zetaclientd -val val \
        --peer /dns/${SEED_NODE}/tcp/6668/p2p/${SEED_NODE_ID} \
        2>&1 | tee ~/.zetaclient/zetaclient.log

fi
