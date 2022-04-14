#!/bin/bash

echo "Starting ZetaClient"
echo $1 $2 $3

NODE_NUMBER=$1
NODE_0_DNS=$2

echo "This is Node $NODE_NUMBER"
 
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin
export MYIP=$(hostname -i)
export IDX=$NODE_NUMBER 
export TSSPATH=/root/.tssnew 

NODE_0_ID=$(cat /zetashared/node0/config/NODE_VALIDATOR_ID)
mkdir -p /root/.tssnew/

sleep 5 # Waiting for Zetacored to boot

if  (( $NODE_NUMBER == 0 )); then
    yes | zetaclientd -val val 2>&1 | tee ~/.zetaclient/zetaclient.log
else
    export SEED=172.24.0.220
        until [ -f NODE_0_TSS_ID ]
        do
            echo "Waiting for Node 0 Validator ID"
            sleep 5
            curl ${SEED}:8123/p2p -o NODE_0_TSS_ID
        done
    NODE_0_TSS_ID=$(cat NODE_0_TSS_ID)
    echo "NODE_0_TSS_ID=${NODE_0_TSS_ID}"
    yes | zetaclientd -val val \
        --peer /dns/${NODE_0_DNS}/tcp/6668/p2p/${NODE_0_TSS_ID} \
        2>&1 | tee ~/.zetaclient/zetaclient.log

fi



