#!/bin/bash
# TODO - Clean this up, no longer needs to treat node 0 differently than node 1+
echo "Starting Zetacore"

NODE_NUMBER=$1
NODE_0_DNS=$2
 
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin
export MYIP=$(hostname -i)

FILE="/root/.zetacore/config/app.toml"
if  [[ -f "$FILE" ]]; then
    echo "This is Node $NODE_NUMBER"
    echo "$FILE already exists."
    echo "Skipping Config Copy From /zetashared/node$NODE_NUMBER/"
    zetacored start --rpc.laddr "tcp://0.0.0.0:26657" \
        --proxy_app "tcp://0.0.0.0:26658" \
        --rpc.pprof_laddr "0.0.0.0:6060" 2>&1 | tee /root/.zetacore/zetacored.log

else
    echo "This is Node $NODE_NUMBER"
    echo "Copying Config From /zetashared/node$NODE_NUMBER/"
    cp -rf /zetashared/node$NODE_NUMBER/* /root/.zetacore/
    zetacored start --rpc.laddr "tcp://0.0.0.0:26657" \
        --proxy_app "tcp://0.0.0.0:26658" \
        --rpc.pprof_laddr "0.0.0.0:6060" 2>&1 | tee /root/.zetacore/zetacored.log
fi

