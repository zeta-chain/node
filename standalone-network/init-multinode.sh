#!/usr/bin/env bash

### chain init script for development purposes only ###
rm -rf ~/.zetacored
kill -9 $(lsof -ti:26657)

NODES="node0 node1 node2 node3"

genesis_addresses=()
for NODE in $NODES; do
  mkdir -p $HOME/.zetacored/$NODE
  zetacored init test --chain-id=localnet_101-1 -o --home=$HOME/.zetacored/$NODE
  zetacored keys add zeta --algo secp256k1 --keyring-backend=test --home=$HOME/.zetacored/$NODE
  zetacored keys add mario --algo secp256k1 --keyring-backend=test --home=$HOME/.zetacored/$NODE
  genesis_addresses+=$(zetacored keys show zeta -a --keyring-backend=test --home=$HOME/.zetacored/$NODE)
  genesis_addresses+=" "
  genesis_addresses+=$(zetacored keys show mario -a --keyring-backend=test --home=$HOME/.zetacored/$NODE)
  genesis_addresses+=" "
done

for address in $genesis_addresses; do
   zetacored add-genesis-account $address 500000000000000000000000000000000azeta --keyring-backend=test --home=$HOME/.zetacored/node0
done

zetacored gentx zeta 10000000000000000azeta --chain-id=localnet_101-1 --keyring-backend=test --home=$HOME/.zetacored/node0
for NODE in $NODES; do
  if [ $NODE != "node0" ]
  then
    rm -rf $HOME/.zetacored/$NODE/config/genesis.json
    cp $HOME/.zetacored/node0/config/genesis.json $HOME/.zetacored/$NODE/config/
    zetacored gentx zeta 10000000000000000azeta --chain-id=localnet_101-1 --keyring-backend=test --home=$HOME/.zetacored/$NODE
    cp $HOME/.zetacored/$NODE/config/gentx/*.json $HOME/.zetacored/node0/config/gentx/
  fi
done


cat $HOME/.zetacored/node0/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacored/node0/config/tmp_genesis.json && mv $HOME/.zetacored/node0/config/tmp_genesis.json $HOME/.zetacored/node0/config/genesis.json
cat $HOME/.zetacored/node0/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacored/node0/config/tmp_genesis.json && mv $HOME/.zetacored/node0/config/tmp_genesis.json $HOME/.zetacored/node0/config/genesis.json
cat $HOME/.zetacored/node0/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacored/node0/config/tmp_genesis.json && mv $HOME/.zetacored/node0/config/tmp_genesis.json $HOME/.zetacored/node0/config/genesis.json
cat $HOME/.zetacored/node0/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacored/node0/config/tmp_genesis.json && mv $HOME/.zetacored/node0/config/tmp_genesis.json $HOME/.zetacored/node0/config/genesis.json
cat $HOME/.zetacored/node0/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacored/node0/config/tmp_genesis.json && mv $HOME/.zetacored/node0/config/tmp_genesis.json $HOME/.zetacored/node0/config/genesis.json
cat $HOME/.zetacored/node0/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.zetacored/node0/config/tmp_genesis.json && mv $HOME/.zetacored/node0/config/tmp_genesis.json $HOME/.zetacored/node0/config/genesis.json

zetacored collect-gentxs --home=$HOME/.zetacored/node0
zetacored validate-genesis --home=$HOME/.zetacored/node0
for NODE in $NODES; do
  if [ $NODE != "node0" ]
  then
    rm -rf $HOME/.zetacored/$NODE/config/genesis.json
    cp $HOME/.zetacored/node0/config/genesis.json $HOME/.zetacored/$NODE/config/
  fi
done

killall zetacored


p2p=27655
grpc=9085
grpcweb=9093
tcp=27659
rpcladdr=2665
jsonrpc=8545
ws=9545

for NODE in $NODES; do
  echo "Starting $NODE"
  zetacored start \
  --minimum-gas-prices=0.0001azeta \
  --home $HOME/.zetacored/$NODE \
  --p2p.laddr 0.0.0.0:$p2p  \
  --grpc.address 0.0.0.0:$grpc \
  --grpc-web.address 0.0.0.0:$grpcweb \
  --address tcp://0.0.0.0:$tcp \
  --json-rpc.address 0.0.0.0:$jsonrpc \
  --json-rpc.ws-address 0.0.0.0:$ws \
  --rpc.laddr tcp://127.0.0.1:$rpcladdr  >> $HOME/.zetacored/$NODE/abci.log 2>&1 &
  ((p2p=p2p+1))
  ((grpc=grpc+1))
  ((grpcweb=grpcweb+1))
  ((tcp=tcp+1))
  ((rpcladdr=rpcladdr+1))
  ((jsonrpc=jsonrpc+1))
  ((ws=ws+1))
done

# TODO ADD peers to config