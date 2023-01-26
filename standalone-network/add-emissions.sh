zetacored tx emissions add-token-emission ValidatorEmission 100000000000 --from=zeta --keyring-backend=test --chain-id=localnet_101-1 --fees=200000azeta --yes
zetacored q emissions list-balances


#Zeta rewards self Delegation
zetacored q distribution rewards zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq


#Mario delegate
zetacored tx staking delegate zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq 1000000000000000000000000azeta --from=mario --keyring-backend=test --chain-id=localnet_101-1 --fees=200000azeta
zetacored q distribution rewards zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50 zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq

#Zeta outstanding rewards (Total)
zetacored q distribution validator-outstanding-rewards zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq

#Zeta commission rewards
zetacored q distribution commission zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq
