#!/usr/bin/env bash

CHAINID="localnet_101-1"
KEYRING="test"
export DAEMON_HOME=$HOME/.zetacored
export DAEMON_NAME=zetacored

### chain init script for development purposes only ###
rm -rf ~/.zetacored
kill -9 $(lsof -ti:26657)
zetacored config keyring-backend $KEYRING --home ~/.zetacored
zetacored config chain-id $CHAINID --home ~/.zetacored
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | zetacored keys add zeta --algo=secp256k1 --recover --keyring-backend=test
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | zetacored keys add mario --algo secp256k1 --recover --keyring-backend=test
echo "lounge supply patch festival retire duck foster decline theme horror decline poverty behind clever harsh layer primary syrup depart fantasy session fossil dismiss east" | zetacored keys add zetaeth --recover --keyring-backend=test

zetacored init test --chain-id=$CHAINID

#Set config to use azeta
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json






zetacored add-genesis-account $(zetacored keys show zeta -a --keyring-backend=test) 500000000000000000000000000000000000000azeta --keyring-backend=test
zetacored add-genesis-account $(zetacored keys show mario -a --keyring-backend=test) 50000000000000000000000000000000000000azeta --keyring-backend=test
zetacored add-genesis-account $(zetacored keys show zetaeth -a --keyring-backend=test) 500000000000000000000000000000000azeta --keyring-backend=test


ADDR1=$(zetacored keys show zeta -a --keyring-backend=test)
observer+=$ADDR1
observer+=","
ADDR2=$(zetacored keys show mario -a --keyring-backend=test)
observer+=$ADDR2
observer+=","


observer_list=$(echo $observer | rev | cut -c2- | rev)

echo $observer_list



zetacored add-observer 1337 "$observer_list"
zetacored add-observer 101 "$observer_list"




zetacored gentx zeta 50000000000000000000000000azeta --chain-id=localnet_101-1 --keyring-backend=test

contents="$(jq '.app_state.gov.voting_params.voting_period = "10s"' $DAEMON_HOME/config/genesis.json)" && \
echo "${contents}" > $DAEMON_HOME/config/genesis.json

echo "Collecting genesis txs..."
zetacored collect-gentxs

echo "Validating genesis file..."
zetacored validate-genesis
#
#export DUMMY_PRICE=yes
#export DISABLE_TSS_KEYGEN=yes
#export GOERLI_ENDPOINT=https://goerli.infura.io/v3/faf5188f178a4a86b3a63ce9f624eb1b
