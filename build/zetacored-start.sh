#!/bin/bash

NODE_NUMBER=$1
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin
export MYIP=$(hostname -i)

echo "Starting Zetacore Node $NODE_NUMBER"

FILE="/root/.zetacore/config/app.toml"
if  [[ ! -f "$FILE" ]]; then
    echo "Copying Config From /zetashared/node$NODE_NUMBER/"
    cp -rf /zetashared/node$NODE_NUMBER/* /root/.zetacore/else
fi

zetacored start --rpc.laddr "tcp://0.0.0.0:26657" \
    --proxy_app "tcp://0.0.0.0:26658" \
    --rpc.pprof_laddr "0.0.0.0:6060" 2>&1 | tee /root/.zetacore/zetacored.log
