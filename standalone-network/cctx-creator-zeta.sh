#!/usr/bin/env bash
set -x

zetacored tx observer add-observer 5 InBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
#zetacored tx observer add-observer 5 OutBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx observer add-observer 2374 InBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx observer add-observer 2374 OutBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx crosschain gas-price-voter 2374 10000000000 100 100 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored tx crosschain nonce-voter Goerli 2374  --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta --broadcast-mode=block

zetacored tx crosschain inbound-voter 0x96B05C238b99768F349135de0653b687f9c13fEE 5 0x3b9Fe88DE29efD13240829A0c18E9EC7A44C3CA7 2374 10000000000000000000 0 "" "0x19398991572a825894b34b904ac1e3692720895351466b5c9e6bb7ae1e21d680" 100 Zeta --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored q crosschain list-cctx
exit 0
zetacored tx crosschain outbound-voter 0x752139735699b0ffa87571bf519867ab2aaf355733316842f9941e0efe5d05c9 hashout 1 7997428181981842964 0 2374 1 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta

#zetacored tx crosschain zeta-conversion-rate-voter Goerli 1 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
