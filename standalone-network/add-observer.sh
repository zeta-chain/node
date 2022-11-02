#!/usr/bin/env bash
zetacored tx observer add-observer Eth InBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx observer add-observer Eth OutBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y
zetacored q observer list-observer

zetacored tx staking unbond zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq 10000000000azeta --from zeta --gas=auto --gas-prices=1azeta --gas-adjustment=3.0 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
