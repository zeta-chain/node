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
    
# FILE="/root/.tssnew/address_book.seed"
# if [ ! -f "$FILE" ]; then
#     echo "$FILE does not exist - Copying from /zetashared/node$NODE_NUMBER/"
#     mkdir -p /root/.tssnew/
#     # cp -rf /zetashared/node$NODE_NUMBER/tssnew/* /root/.tssnew/
# fi

if  (( $NODE_NUMBER == 0 )); then
    yes | zetaclientd -val val 2>&1 | tee ~/.zetaclient/zetaclient.log
else
    export SEED=172.24.0.220
    yes | zetaclientd -val val \
        --peer /dns/${NODE_0_DNS}/tcp/6668/p2p/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp \
        2>&1 | tee ~/.zetaclient/zetaclient.log
fi



