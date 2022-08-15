#!/bin/bash

zetacored tx wasm execute zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql \
  '{"add_to_watch_list":{"chain":"ETH","nonce":1,"tx_hash":"123"}}' \
  --from zeta \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y

 zetacored query wasm contract-state smart zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql \
   '{"watchlist":{}}' \
   --chain-id localnet