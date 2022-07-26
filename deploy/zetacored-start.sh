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

FILE="~/.zetacore/config/app.toml"
if  [[ ! -f "$FILE" ]]; then
    echo "Copying Config From /zetashared/node$NODE_NUMBER/"
    cp -rf /zetashared/node"$NODE_NUMBER"/* ~/.zetacore/
fi

zetacored start \
    --rpc.laddr "tcp://0.0.0.0:26657" \
    --rpc.pprof_laddr "0.0.0.0:6060"  \
    --address "tcp://$MYIP:26658" \
    --log_format json \
    --moniker "node$NODE_NUMBER"
