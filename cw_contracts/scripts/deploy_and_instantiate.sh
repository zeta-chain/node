#!/bin/bash

zetacored tx wasm store ../contract/watcher.wasm --from=zeta --keyring-backend=test --broadcast-mode block --chain-id=localnet -y --gas=1000000000000000000

sleep 7
zetacored tx wasm instantiate 1 '{}' \
--amount 50000stake \
--label "watcher" \
--no-admin \
--from zeta --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet \
-y

# Contract : zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql