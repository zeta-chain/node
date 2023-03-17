zetacored tx authz grant zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50 delegate --allowed-validators zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --spend-limit=100000000000azeta

zetacored q authz grants zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50
zetacored tx authz exec tx.json --from mario --gas=auto --gas-prices=0.1azeta --gas-adjustment=20 --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored q staking delegations-to zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq



zetacored tx authz grant zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50 generic --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --msg-type=/zetachain.zetacore.crosschain.MsgGasPriceVoter

