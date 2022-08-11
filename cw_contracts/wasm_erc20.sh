#!/usr/bin/env bash

zetacored tx wasm store cw1_subkeys.wasm \
--from zeta --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet  \
-y

zetacored tx wasm instantiate 1 '{"admins":["zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk"],"mutable":false}' \
--amount 5000000stake \
--label "CW1 Subkey" \
--no-admin \
--from zeta --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet \
-y

#Receiptient is contract : <zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql>

zetacored q bank balances zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql
#Query mario's allowance
zetacored query wasm contract-state smart zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql  \
 '{"allowance":{"spender":"zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50"}}' \
 --chain-id localnet

zetacored query wasm contract-state smart zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql \
 '{"all_allowances":{}}' \
 --chain-id localnet

zetacored query wasm contract-state smart zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql \
 '{"admin_list":{}}' \
 --chain-id localnet


# Add allowance for Mario
zetacored tx wasm execute zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql \
  '{"increase_allowance":{"spender":"zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50","amount":{"denom":"stake","amount":"200000"}}}' \
  --from zeta \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y

zetacored q bank balances $(zetacored keys show -a zeta --keyring-backend test) | grep -B1 stake
zetacored q bank balances $(zetacored keys show -a mario --keyring-backend test) | grep -B1 stake

#Proxy send from mario
zetacored tx wasm execute zeta14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62znql \
  '{"execute":{"msgs":[{"bank":{"send":{"to_address":"zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk","amount":[{"denom":"stake","amount":"999"}]}}}]}}' \
  --from mario \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y