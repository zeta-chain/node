zetacored tx authz grant zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50 delegate --allowed-validators zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block

zetacored q authz grants zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50
zetacored tx authz exec tx.json --from mario --gas=auto --gas-prices=0.1azeta --gas-adjustment=1.5 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored q staking delegations-to zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq
