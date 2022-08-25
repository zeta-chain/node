#!/bin/bash

NODE_NUMBER=$1
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:~/go/bin

if [ -z "${MYIP}" ]; then 
    # If MYIP is not set, use the private IP of the host
    echo "MYIP ENV Variable Not Set -- Setting it automatically using host IP"
    export MYIP=$(hostname -i)
fi

echo "Starting Zetacore Node $NODE_NUMBER"

FILE="~/.zetacored/config/app.toml"
if  [[ ! -f "$FILE" ]]; then
    echo "Copying Config From /zetashared/node$NODE_NUMBER/"
    cp -rf /zetashared/node"$NODE_NUMBER"/* ~/.zetacored/
fi

zetacored start --trace \
    --home ~/.zetacored \
    --address "tcp://$MYIP:26658" \
    --rpc.laddr "tcp://0.0.0.0:26657" \
    --rpc.pprof_laddr "0.0.0.0:6060"  \
    --moniker "node$NODE_NUMBER" \
    --log_format json \
    --pruning custom \
    --pruning-keep-recent 3 \
    --pruning-keep-every 100  \
    --pruning-interval 10 \
    --state-sync.snapshot-interval 1000 \
    --state-sync.snapshot-keep-recent 1 

