zetacored tx  group create-group-with-policy zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax group-metadata group-policy-metadata members.json policy.json --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx group submit-proposal proposal_keygen.json --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block --exec=1
zetacored tx group vote 2 zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax VOTE_OPTION_YES metadata --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored tx group vote 2 zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2 VOTE_OPTION_YES metadata --from mario --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored q group proposal 2
zetacored tx group exec 2 --from zeta --fees=40azeta --chain-id=localnet_101-1 --keyring-backend=test -y --broadcast-mode=block
zetacored q group group-policies-by-admin zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax


#zetacored q group group-policy-info zeta1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsxn0x73
#zetacored q group group-policies-by-group 1