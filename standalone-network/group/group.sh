zetacored tx  group create-group-with-policy zeta19wzjdtah4kl2vh77jks68cyy5gpjyurqltys99 group-metadata group-policy-metadata members.json policy.json --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx group submit-proposal proposal_keygen.json --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --exec=1
zetacored tx group vote 3 zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk VOTE_OPTION_YES metadata --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx group vote 3 zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50 VOTE_OPTION_YES metadata --from mario --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored q group proposal 3
zetacored tx group exec 3 --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored q group group-policies-by-admin zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk


#zetacored q group group-policy-info zeta1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsxn0x73
#zetacored q group group-policies-by-group 1