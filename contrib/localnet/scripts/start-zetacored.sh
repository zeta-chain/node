#!/bin/bash

HOSTNAME=$(hostname)
cp -r "/root/zetacored/zetacored_$HOSTNAME" /root/.zetacored

if [ -f "/root/.zetacored/data/priv_validator_state.json" ]; then
  echo "priv_validator_state.json already exists"
else
  echo "priv_validator_state.json does not exist; creating an empty one"
  mkdir -p /root/.zetacored/data
  cp -r "/root/zetacored/zetacored_$HOSTNAME/priv_validator_state.json" /root/.zetacored/data
fi

/root/zetacored-proposal.sh &

zetacored start --pruning=nothing   --minimum-gas-prices=0.0001azeta --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable --home /root/.zetacored