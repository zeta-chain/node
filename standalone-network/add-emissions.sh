zetacored tx bank send zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk zeta1w43fn2ze2wyhu5hfmegr6vp52c3dgn0srdgymy 100000000000azeta --keyring-backend=test --chain-id=localnet_101-1 --fees=200000azeta --yes

#total left in reserves
zetacored q bank balances zeta1w43fn2ze2wyhu5hfmegr6vp52c3dgn0srdgymy
#observer undistributed
zetacored q bank balances zeta1pyks89mqljlpgzenwa0g8zch0hptk6usd9vcuh
#tss undistributed
zetacored q bank balances zeta1v8v7zkyt7j3dc526k4alsu8vspvqqg342t27vu


#Zeta rewards self Delegation
zetacored q distribution rewards zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq


#Mario delegate
zetacored tx staking delegate zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq 1000000000000000000000000azeta --from=mario --keyring-backend=test --chain-id=localnet_101-1 --fees=200000azeta
zetacored q distribution rewards zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50 zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq

#Zeta outstanding rewards (Total)
zetacored q distribution validator-outstanding-rewards zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq

#Zeta commission rewards
zetacored q distribution commission zetavaloper1syavy2npfyt9tcncdtsdzf7kny9lh777nep4tq
