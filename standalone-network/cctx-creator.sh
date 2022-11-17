#!/usr/bin/env bash
zetacored tx crosschain gas-price-voter Goerli 1000000000000 2100000000000000 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored tx crosschain zeta-conversion-rate-voter Goerli 1 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored tx crosschain nonce-voter Goerli 1  --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta



zetacored tx crosschain inbound-voter 0x96B05C238b99768F349135de0653b687f9c13fEE 5 0x96B05C238b99768F349135de0653b687f9c13fEE 5 1000000000000000000 0 message hash 100 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored tx crosschain inbound-voter 0x96B05C238b99768F349135de0653b687f9c13fEE 5 0x96B05C238b99768F349135de0653b687f9c13fEE 5 1000000000000000000 0 message hash 100 1 --from=mario --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta

zetacored q crosschain list-cctx
zetacored tx crosschain outbound-voter 0xa968ea9d648d5759ec66ed6dc55790ac6465f167009d180b3e7c68a5cc5e06a9 hashout 1 0 0 5 1 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored tx crosschain outbound-voter 0xa968ea9d648d5759ec66ed6dc55790ac6465f167009d180b3e7c68a5cc5e06a9 hashout 1 0 0 5 1 1 --from=mario --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta


