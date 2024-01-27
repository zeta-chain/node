#!/usr/bin/env bash

CHAINID="localnet_101-1"
KEYRING="test"
HOSTNAME=$(hostname)
signer="zeta"


killall zetacored
zetacored start --trace \
--minimum-gas-prices=0.0001azeta \
--json-rpc.api eth,txpool,personal,net,debug,web3,miner \
--api.enable >> ~/.zetacored/zetacored.log 2>&1  & \
#>> "$HOME"/.zetacored/zetanode.log 2>&1  & \

#--home ~/.zetacored \
#--p2p.laddr 0.0.0.0:27655  \
#--grpc.address 0.0.0.0:9096 \
#--grpc-web.address 0.0.0.0:9093 \
#--address tcp://0.0.0.0:27659 \
#--rpc.laddr tcp://127.0.0.1:26657 \
#--pruning custom \
#--pruning-keep-recent 54000 \
#--pruning-interval 10 \
#--min-retain-blocks 54000 \
#--state-sync.snapshot-interval 14400 \
#--state-sync.snapshot-keep-recent 3

#echo "--> Submitting proposal to update admin policies "
#sleep 7
#zetacored tx gov submit-legacy-proposal param-change standalone-network/proposal.json --from $signer --gas=auto --gas-adjustment=1.5 --gas-prices=0.001azeta --chain-id=$CHAINID --keyring-backend=$KEYRING -y --broadcast-mode=block
#echo "--> Submitting vote for proposal"
#sleep 7
#zetacored tx gov vote 1 yes --from $signer --keyring-backend $KEYRING --chain-id $CHAINID --yes --fees=40azeta --broadcast-mode=block
sleep 7
zetacored tx fungible deploy-system-contracts --from zeta --gas=auto --gas-prices=10000000000azeta --gas-adjustment=1.5 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block

tail -f ~/.zetacored/zetacored.log

