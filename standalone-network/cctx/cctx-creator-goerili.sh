set -x


zetacored tx crosschain gas-price-voter 1337 10000000000 100 100 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored tx crosschain create-tss-voter tsspubkey 5 0 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored tx crosschain gas-price-voter 1337 10000000000 100 100 --from=mario --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta
zetacored tx crosschain create-tss-voter tsspubkey 5 0 --from=mario --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta

zetacored tx crosschain inbound-voter \
0x96B05C238b99768F349135de0653b687f9c13fEE \
1337 \
0x3b9Fe88DE29efD13240829A0c18E9EC7A44C3CA7 \
0x96B05C238b99768F349135de0653b687f9c13fEE \
1337 \
10000000000000000000 \
"" \
"0x19398991572a825894b34b904ac1e3692720895351466b5c9e6bb7ae1e21d680" \
100 \
Gas \
"" \
--from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta


zetacored tx crosschain inbound-voter \
0x96B05C238b99768F349135de0653b687f9c13fEE \
1337 \
0x3b9Fe88DE29efD13240829A0c18E9EC7A44C3CA7 \
0x96B05C238b99768F349135de0653b687f9c13fEE \
1337 \
10000000000000000000 \
"" \
"0x19398991572a825894b34b904ac1e3692720895351466b5c9e6bb7ae1e21d680" \
100 \
Gas \
"" \
--from=mario --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta --output=json

zetacored q crosschain list-cctx

zetacored tx crosschain outbound-voter \
0xead687de84b3969b4c18480f197d2812e0acb83f851acc2830f70e94c85cef55 \
hashout \
1 \
7994721005120625032 \
0 \
1337 \
1 \
Gas \
--from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta

zetacored tx crosschain outbound-voter \
0xead687de84b3969b4c18480f197d2812e0acb83f851acc2830f70e94c85cef55 \
hashout \
1 \
7994721005120625032 \
0 \
1337 \
1 \
Gas \
--from=mario --keyring-backend=test --yes --chain-id=localnet_101-1 --broadcast-mode=block --gas=auto --gas-adjustment=2 --gas-prices=0.1azeta




