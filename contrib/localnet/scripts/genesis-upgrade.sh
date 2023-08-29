#!/bin/bash

/usr/sbin/sshd

if [ $# -ne 1 ]
then
  echo "Usage: genesis.sh <num of nodes>"
  exit 1
fi
NUMOFNODES=$1

# create keys
CHAINID="athens_101-1"
KEYRING="test"
#DISCOVERED_HOSTNAME=$(hostname)
DISCOVERED_HOSTNAME=$(nslookup $(hostname -i) | grep '=' | awk -F'= ' '{split($2, a, "."); print a[1]}')
DISCOVERED_NETWORK=$(echo $DISCOVERED_HOSTNAME |  awk -F'-' '{split($1, a, "-"); print a[1]}')
INDEX=${DISCOVERED_HOSTNAME:0-1}

# generate node list
START=2
# shellcheck disable=SC2100
END=$((NUMOFNODES))
NODELIST=()
for i in $(eval echo "{$START..$END}")
do
  NODELIST+=("zetacore_node-$i")
done

echo "DISCOVERED_HOSTNAME: $DISCOVERED_HOSTNAME"

# Init a new node to generate genesis file .
# Copy config files from existing folders which get copied via Docker Copy when building images
mkdir -p ~/.backup/config
zetacored init Zetanode-Localnet --chain-id=$CHAINID
rm -rf ~/.zetacored/config/app.toml
rm -rf ~/.zetacored/config/client.toml
rm -rf ~/.zetacored/config/config.toml
cp -r ~/zetacored/common/app.toml ~/.zetacored/config/
cp -r ~/zetacored/common/client.toml ~/.zetacored/config/
cp -r ~/zetacored/common/config.toml ~/.zetacored/config/
sed -i -e "/moniker =/s/=.*/= \"$DISCOVERED_HOSTNAME\"/" "$HOME"/.zetacored/config/config.toml

# Add two new keys for operator and hotkey and create the required json structure for os_info
source ~/os-info.sh

# Pause other nodes so that the primary can node can do the genesis creation
if [ $DISCOVERED_HOSTNAME != "$DISCOVERED_NETWORK-zetacore_node-1" ]
then
  echo "Waiting for $DISCOVERED_NETWORK-zetacore_node-1 to create genesis.json"
  sleep $((7*NUMOFNODES))
  echo "genesis.json created"
fi

# Genesis creation following steps
# 1. Accumulate all the os_info files from other nodes on zetcacore0 and create a genesis.json
# 2. Add the observers , authorizations and required params to the genesis.json
# 3. Copy the genesis.json to all the nodes .And use it to create a gentx for every node
# 4. Collect all the gentx files in zetacore_node-1 and create the final genesis.json
# 5. Copy the final genesis.json to all the nodes and start the nodes
# 6. Update Config in zetacore_node-1 so that it has the correct persistent peer list
# 7. Start the nodes

# Start of genesis creation . This is done only on $DISCOVERED_NETWORK-zetacore_node-1
if [ $DISCOVERED_HOSTNAME == "$DISCOVERED_NETWORK-zetacore_node-1" ]
then
  # Misc : Copying the keyring to the client nodes so that they can sign the transactions
  ssh $DISCOVERED_NETWORK-zetaclient-1 mkdir -p ~/.zetacored/keyring-test/
  scp ~/.zetacored/keyring-test/* $DISCOVERED_NETWORK-zetaclient-1:~/.zetacored/keyring-test/

# 1. Accumulate all the os_info files from other nodes on $DISCOVERED_NETWORK-zetacore_node-1 and create a genesis.json
  for NODE in "${NODELIST[@]}"; do
    INDEX=${NODE:0-1}
    ssh $DISCOVERED_NETWORK-zetaclient-"$INDEX" mkdir -p ~/.zetacored/
    scp "$DISCOVERED_NETWORK-$NODE":~/.zetacored/os_info/os.json ~/.zetacored/os_info/os_z"$INDEX".json
    scp ~/.zetacored/os_info/os_z"$INDEX".json $DISCOVERED_NETWORK-zetaclient-"$INDEX":~/.zetacored/os.json
  done

  ssh $DISCOVERED_NETWORK-zetaclient-1 mkdir -p ~/.zetacored/
  scp ~/.zetacored/os_info/os.json $DISCOVERED_NETWORK-zetaclient-1:/root/.zetacored/os.json

# 2. Add the observers , authorizations and required params to the genesis.json
  zetacored collect-observer-info
  zetacored add-observer-list
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="500000000"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="100s"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json


# 3. Copy the genesis.json to all the nodes .And use it to create a gentx for every node
  zetacored gentx operator 1000000000000000000000azeta --chain-id=$CHAINID --keyring-backend=$KEYRING
  # Copy host gentx to other nodes
  for NODE in "${NODELIST[@]}"; do
    ssh $DISCOVERED_NETWORK-$NODE mkdir -p ~/.zetacored/config/gentx/peer/
    scp ~/.zetacored/config/gentx/* $DISCOVERED_NETWORK-$NODE:~/.zetacored/config/gentx/peer/
  done
  # Create gentx files on other nodes and copy them to host node
  mkdir ~/.zetacored/config/gentx/z2gentx
  for NODE in "${NODELIST[@]}"; do
      ssh $DISCOVERED_NETWORK-$NODE rm -rf ~/.zetacored/genesis.json
      scp ~/.zetacored/config/genesis.json $DISCOVERED_NETWORK-$NODE:~/.zetacored/config/genesis.json
      ssh $DISCOVERED_NETWORK-$NODE zetacored gentx operator 1000000000000000000000azeta --chain-id=$CHAINID --keyring-backend=$KEYRING
      scp $DISCOVERED_NETWORK-$NODE:~/.zetacored/config/gentx/* ~/.zetacored/config/gentx/
      scp $DISCOVERED_NETWORK-$NODE:~/.zetacored/config/gentx/* ~/.zetacored/config/gentx/z2gentx/
  done

# 4. Collect all the gentx files in zetacore_node-1 and create the final genesis.json
  zetacored collect-gentxs
  zetacored validate-genesis
# 5. Copy the final genesis.json to all the nodes
  for NODE in "${NODELIST[@]}"; do
      ssh $DISCOVERED_NETWORK-$NODE rm -rf ~/.zetacored/genesis.json
      scp ~/.zetacored/config/genesis.json $DISCOVERED_NETWORK-$NODE:~/.zetacored/config/genesis.json
  done
# 6. Update Config in zetacore_node-1 so that it has the correct persistent peer list
   sleep 2
   pp=$(cat $HOME/.zetacored/config/gentx/z2gentx/*.json | jq '.body.memo' )
   pps=${pp:1:58}
   sed -i -e "/persistent_peers =/s/=.*/= \"$pps\"/" "$HOME"/.zetacored/config/config.toml
fi
# End of genesis creation steps . The steps below are common to all the nodes

# Update persistent peers
if [ $DISCOVERED_HOSTNAME != "$DISCOVERED_NETWORK-zetacore_node-1" ]
then
  # Misc : Copying the keyring to the client nodes so that they can sign the transactions
  ssh $DISCOVERED_NETWORK-zetaclient-"$INDEX" mkdir -p ~/.zetacored/keyring-test/
  scp ~/.zetacored/keyring-test/* "$DISCOVERED_NETWORK-zetaclient-$INDEX":~/.zetacored/keyring-test/

  pp=$(cat $HOME/.zetacored/config/gentx/peer/*.json | jq '.body.memo' )
  pps=${pp:1:58}
  sed -i -e "/persistent_peers =/s/=.*/= \"$pps\"/" "$HOME"/.zetacored/config/config.toml
fi

/root/zetacored-proposal.sh &

# 7 Start the nodes
exec zetacored start --pruning=nothing --minimum-gas-prices=0.0001azeta --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable --home /root/.zetacored