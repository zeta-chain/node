#!/usr/bin/env bash
zetacored tx observer add-observer 5 InBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx observer add-observer 5 OutBoundTx --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx crosschain gas-price-voter 5 1000000000000 2100000000000000 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
#zetacored tx crosschain zeta-conversion-rate-voter Goerli 1 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored tx crosschain nonce-voter Goerli 1  --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta --broadcast-mode=block



zetacored tx crosschain inbound-voter 0x96B05C238b99768F349135de0653b687f9c13fEE 5 0x96B05C238b99768F349135de0653b687f9c13fEE 5 1000000000000000000 0 message hash 100 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta

zetacored q crosschain list-cctx
zetacored tx crosschain outbound-voter 0x18b11d2eccd0601e3c0ffaef981892d579e74ed7eb36e7ab159f5c410c829a57 hashout 1 0 0 5 1 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta


