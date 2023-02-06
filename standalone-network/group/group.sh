zetacored tx  group create-group-with-policy zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk group-metadata group-policy-metadata members.json policy.json --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx group submit-proposal proposal_group.json --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --exec=1
zetacored tx group vote 1 zeta1syavy2npfyt9tcncdtsdzf7kny9lh777heefxk VOTE_OPTION_YES metadata --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx group vote 1 zeta1l7hypmqk2yc334vc6vmdwzp5sdefygj2w5yj50 VOTE_OPTION_YES metadata --from mario --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
exit 0
zetacored q group proposal 1
zetacored tx group exec 1 --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block

