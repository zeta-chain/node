#!/bin/bash

#DISCOVERED_HOSTNAME=$(hostname)
DISCOVERED_HOSTNAME=$(nslookup $(hostname -i) | grep '=' | awk -F'= ' '{split($2, a, "."); print a[1]}')
DISCOVERED_NETWORK=$(echo $DISCOVERED_HOSTNAME |  awk -F'-' '{split($1, a, "-"); print a[1]}')
cp -r "/root/zetacored/zetacored_$DISCOVERED_HOSTNAME" /root/.zetacored

if [ -f "/root/.zetacored/data/priv_validator_state.json" ]; then
  echo "priv_validator_state.json already exists"
else
  echo "priv_validator_state.json does not exist; creating an empty one"
  mkdir -p /root/.zetacored/data
  cp -r "/root/zetacored/zetacored_$DISCOVERED_HOSTNAME/priv_validator_state.json" /root/.zetacored/data
fi

/root/zetacored-proposal.sh &

zetacored start --pruning=nothing   --minimum-gas-prices=0.0001azeta --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable --home /root/.zetacored