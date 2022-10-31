#!/usr/bin/env bash
zetacored tx observer add-observer Eth InBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y
zetacored tx observer add-observer Eth OutBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y