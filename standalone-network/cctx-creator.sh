#!/usr/bin/env bash

zetacored tx zetacore inbound-voter 0x96B05C238b99768F349135de0653b687f9c13fEE ETH 0x96B05C238b99768F349135de0653b687f9c13fEE ETH 1000000000000000000 0 message hash 100 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1
zetacored tx zetacore inbound-voter 0x96B05C238b99768F349135de0653b687f9c13fEE ETH 0x96B05C238b99768F349135de0653b687f9c13fEE ETH 1000000000000000000 0 message hash 100 --from=mario --keyring-backend=test --yes --chain-id=localnet_101-1
zetacored q zetacore list-cctx

zetacored tx zetacore outbound-voter 0x9ea007f0f60e32d58577a8cf25678942d2b10791c2a34f48e237b76a7e998e4d hashout 1 0 0 ETH 1 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1
zetacored tx zetacore outbound-voter 0x9ea007f0f60e32d58577a8cf25678942d2b10791c2a34f48e237b76a7e998e4d hashout 1 0 0 ETH 1 --from=mario --keyring-backend=test --yes --chain-id=localnet_101-1
