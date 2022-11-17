#!/usr/bin/env bash
zetacored tx crosschain gas-price-voter GOERLI 1000000000000 2100000000000000 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=1.5 --gas-prices=0.1azeta
zetacored tx crosschain zeta-conversion-rate-voter GOERLI 1 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=1.5 --gas-prices=0.1azeta
zetacored tx crosschain nonce-voter GOERLI 1  --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=1.5 --gas-prices=0.1azeta


zetacored tx crosschain inbound-voter 0x96B05C238b99768F349135de0653b687f9c13fEE GOERLI 0x96B05C238b99768F349135de0653b687f9c13fEE GOERLI 1000000000000000000 0 message hash 100 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 1 --broadcast-mode=block --gas=auto --gas-adjustment=1.5 --gas-prices=0.1azeta
zetacored tx crosschain inbound-voter 0x96B05C238b99768F349135de0653b687f9c13fEE GOERLI 0x96B05C238b99768F349135de0653b687f9c13fEE GOERLI 1000000000000000000 0 message hash 100 --from=mario --keyring-backend=test --yes --chain-id=localnet_101-1 1 --broadcast-mode=block --gas=auto --gas-adjustment=1.5 --gas-prices=0.1azeta

zetacored q crosschain list-cctx
zetacored tx crosschain outbound-voter 0x953a57379e407310f2ea77c441f7749ebd1247bc8838657cd89cbd8f5c29c6f4 hashout 1 0 0 ETH 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --fees=40azeta
zetacored tx crosschain outbound-voter 0x953a57379e407310f2ea77c441f7749ebd1247bc8838657cd89cbd8f5c29c6f4 hashout 1 0 0 ETH 1 --from=mario --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --fees=40azeta

