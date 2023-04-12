#!/bin/bash

# create keys
CHAINID="athens_101-1"
KEYRING="test"
HOSTNAME=$(hostname)
NODES="zetacore1"
echo "HOSTNAME: $HOSTNAME"
zetacored init Zetanode-Localnet --chain-id=$CHAINID

zetacored config keyring-backend $KEYRING --home ~/.zetacored
zetacored config chain-id $CHAINID --home ~/.zetacored


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




if [ $HOSTNAME == "zetacore0" ]
then
  for NODE in $NODES; do
    scp $NODE:~/.zetacored/os_info/os.json ~/.zetacored/os_info/os_z1.json
  done
  zetacored collect-observer-info
  zetacored add-observer-list
fi

sleep infinity
#  for i in {1..4} ; do
#  scp zetacored$i:~/.zetacored/os.json os$i.json
#  done
#
## concatenate OS JSON files
#
#
## create genesis file
#
#for i in {1..4} ; do
#  scp genesis.json zetacored$i:~/genesis.json
#done
#
#sleep 10
#
#zetacored gentx ...


# start the network

# set peer addresses in config.toml
#jq ".fdlasjkf=sd;lfja" config.toml > config.toml
#
#zetcored start ...