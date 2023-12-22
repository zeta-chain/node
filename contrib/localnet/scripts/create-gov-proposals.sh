#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
cd "$SCRIPT_DIR" || exit

## Create localnet_gov_admin on host system

LOCALNET_GOV_ADMIN_MNEMONIC="mercy oblige six giant absorb crunch derive tornado sleep friend blame border avocado fine script dilemma vacant dad buddy occur trigger energy today minimum"
WALLET_NAME=localnet_gov_admin

if ! zetacored keys show localnet_gov_admin -a --keyring-backend test > /dev/null 2>&1; then
    echo  "Creating localnet_gov_admin key"
    echo "$LOCALNET_GOV_ADMIN_MNEMONIC" | zetacored keys add $WALLET_NAME --keyring-backend test --recover
fi

# Create a few short lived proposals for variety of testing
zetacored tx gov submit-proposal proposals/proposal_for_failure.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 10
zetacored tx gov vote 1 no --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --yes --fees 2000000000000000azeta --yes

zetacored tx gov submit-proposal proposals/proposal_for_success.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 10
zetacored tx gov vote 2 yes --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --yes --fees 2000000000000000azeta --yes

zetacored tx gov submit-proposal proposals/v100.0.0_proposal.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 10
zetacored tx gov vote 3 yes --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --yes --fees 2000000000000000azeta --yes

# Increase the length of the voting period to 1 week
zetacored tx gov submit-proposal proposals/proposal_voting_period.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 10
zetacored tx gov vote 4 yes --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --yes --fees 2000000000000000azeta --yes

# Create a few long lived proposals for variety of testing

zetacored tx gov submit-proposal proposals/proposal_voting_period.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 10
zetacored tx gov vote 5 yes --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --yes --fees 2000000000000000azeta --yes


zetacored tx gov submit-proposal proposals/v100.0.0_proposal.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 10
zetacored tx gov vote 6 yes --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --yes --fees 2000000000000000azeta --yes

zetacored tx gov submit-proposal proposals/proposal_for_deposit.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 10
zetacored tx gov vote 7 yes --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --yes --fees 2000000000000000azeta --yes




