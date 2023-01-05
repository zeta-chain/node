#!/usr/bin/env bash
zetacored tx observer add-observer Eth InBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx observer add-observer Eth OutBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored q observer list-observer

zetacored tx staking unbond zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq 99999900000000000000azeta --from zeta --gas=auto --gas-prices=0.0001azeta --gas-adjustment=1.5 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block


#Add via genesis
#zetacored add-observer-list sample-observer-list.json
#zetacored add-observer Eth InBoundTx $(zetacored keys show zeta -a --keyring-backend=test)
#zetacored add-observer Eth OutBoundTx $(zetacored keys show zeta -a --keyring-backend=test)
#zetacored add-genesis-account $(zetacored keys show zetaeth -a --keyring-backend=test) 50000000000000000000000000000000azeta,500000000000000000000000000000000stake --keyring-backend=test

