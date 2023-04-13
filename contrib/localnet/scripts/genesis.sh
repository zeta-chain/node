#!/bin/bash

/usr/sbin/sshd

# create keys
CHAINID="athens_101-1"
KEYRING="test"
HOSTNAME=$(hostname)
NODES="zetacore1"
echo "HOSTNAME: $HOSTNAME"

mkdir -p ~/.backup/config
zetacored init Zetanode-Localnet --chain-id=$CHAINID
rm -rf ~/.zetacored/config/config.toml
rm -rf ~/.zetacored/config/app.toml
rm -rf ~/.zetacored/config/client.toml
cp -r ~/zetacored/zetacored_"$HOSTNAME"/config/config.toml ~/.zetacored/config/
cp -r ~/zetacored/zetacored_"$HOSTNAME"/config/app.toml ~/.zetacored/config/
cp -r ~/zetacored/zetacored_"$HOSTNAME"/config/client.toml ~/.zetacored/config/



zetacored keys add operator --algo=secp256k1 --keyring-backend=$KEYRING
zetacored keys add hotkey --algo=secp256k1 --keyring-backend=$KEYRING
operator_address=$(zetacored keys show operator -a --keyring-backend=$KEYRING)
hotkey_address=$(zetacored keys show hotkey -a --keyring-backend=$KEYRING)
pubkey=$(zetacored get-pubkey hotkey|sed -e 's/secp256k1:"\(.*\)"/\1/' | sed 's/ //g' )
echo "operator_address: $operator_address"
echo "hotkey_address: $hotkey_address"
echo "pubkey: $pubkey"
mkdir ~/.zetacored/os_info
jq -n --arg operator_address "$operator_address" --arg hotkey_address "$hotkey_address" --arg pubkey "$pubkey" '{"ObserverAddress":$operator_address,"ZetaClientGranteeAddress":$hotkey_address,"ZetaClientGranteePubKey":$pubkey}' > ~/.zetacored/os_info/os.json


if [ $HOSTNAME != "zetacore0" ]
then
  echo "Waiting for zetacore0 to create genesis.json"
  sleep 6
  echo "genesis.json created"
fi

if [ $HOSTNAME == "zetacore0" ]
then
  for NODE in $NODES; do
    scp $NODE:~/.zetacored/os_info/os.json ~/.zetacored/os_info/os_z1.json
  done
    scp ~/.zetacored/os_info/os.json zetaclient0:~/.zetacored/os.json
    scp ~/.zetacored/os_info/os_z1.json zetaclient1:~/.zetacored/os.json
  zetacored collect-observer-info
  zetacored add-observer-list
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
  cat $HOME/.zetacored/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json


  zetacored gentx operator 1000000000000000000000azeta --chain-id=$CHAINID --keyring-backend=$KEYRING
  mkdir ~/.zetacored/config/gentx/z2gentx
  for NODE in $NODES; do
      ssh $NODE rm -rf ~/.zetacored/genesis.json
      scp ~/.zetacored/config/genesis.json $NODE:~/.zetacored/config/genesis.json
      ssh $NODE zetacored gentx operator 1000000000000000000000azeta --chain-id=$CHAINID --keyring-backend=$KEYRING
      scp $NODE:~/.zetacored/config/gentx/* ~/.zetacored/config/gentx/
      scp $NODE:~/.zetacored/config/gentx/* ~/.zetacored/config/gentx/z2gentx/
  done
#  pp=$(cat $HOME/.zetacored/config/gentx/z2gentx/*.json | jq '.body.memo' )
#  pps=${pp:1:58}
#  sed -i -e "/persistent_peers =/s/=.*/= \"$pps\"/" "$HOME"/.zetacored/config/config.toml
  zetacored collect-gentxs
  zetacored validate-genesis
  for NODE in $NODES; do
      ssh $NODE rm -rf ~/.zetacored/genesis.json
      scp ~/.zetacored/config/genesis.json $NODE:~/.zetacored/config/genesis.json
  done
   sleep 2
   pp=$(cat $HOME/.zetacored/config/gentx/z2gentx/*.json | jq '.body.memo' )
   pps=${pp:1:58}
   sed -i -e "/persistent_peers =/s/=.*/= \"$pps\"/" "$HOME"/.zetacored/config/config.toml
fi

if [ $HOSTNAME == "zetacore0" ]
then
  scp ~/.zetacored/keyring-test/* zetaclient0:~/.zetacored/keyring-test/
else
  scp ~/.zetacored/keyring-test/* zetaclient1:~/.zetacored/keyring-test/
fi


exec zetacored start --pruning=nothing --minimum-gas-prices=0.0001azeta --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable --home /root/.zetacored
