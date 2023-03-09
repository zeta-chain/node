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
echo "elevator gorilla destroy silk any hip wrap bench diary adult powder blame coast future gas affair angry write pig country crew beef rack laundry" | zetacored keys add val_tss_signer --algo=secp256k1 --recover --keyring-backend=test
echo "gloom skull world long ready panel special glance method maximum evil book learn salt level logic acoustic basket survey usual crawl olive enemy gown" | zetacored keys add val_grantee_validator --algo secp256k1 --recover --keyring-backend=test
echo "banner powder select list size resource drum cloth clap strike fade duty supply wheel device asthma impact whisper off settle gorilla pizza cage modify" | zetacored keys add val_grantee_observer --recover --keyring-backend=test --algo secp256k1
echo "category lucky claw exhaust hobby hero garden hero tooth involve play decorate general organ shock cloud inquiry provide inspire depend spread simple width orient" | zetacored keys add val --recover --keyring-backend=test --algo secp256k1
zetacored init test --chain-id=$CHAINID

#Set config to use azeta
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json






zetacored add-genesis-account $(zetacored keys show zeta -a --keyring-backend=test) 500000000000000000000000000000000azeta --keyring-backend=test
zetacored add-genesis-account $(zetacored keys show mario -a --keyring-backend=test) 500000000000000000000000000000000azeta --keyring-backend=test
zetacored add-genesis-account $(zetacored keys show executer -a --keyring-backend=test) 500000000000000000000000000000000azeta --keyring-backend=test


#ADDR1=$(zetacored keys show zeta -a --keyring-backend=test)
#observer+=$ADDR1
#observer+=","
#ADDR2=$(zetacored keys show mario -a --keyring-backend=test)
#observer+=$ADDR2
#observer+=","
#
#
#observer_list=$(echo $observer | rev | cut -c2- | rev)
#
#echo $observer_list
#
#
#
#zetacored add-observer 1337 "$observer_list" #goerli
#zetacored add-observer 101  "$observer_list" #goerli
zetacored add-observer-list standalone-network/observers.json





zetacored gentx zeta 1000000000000000000000azeta --chain-id=localnet_101-1 --keyring-backend=test

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
