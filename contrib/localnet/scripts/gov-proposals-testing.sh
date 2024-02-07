#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
cd "$SCRIPT_DIR" || exit

WALLET_NAME=operator

# Create a few short lived proposals for variety of testing
zetacored tx gov submit-proposal proposals/proposal_for_failure.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12
zetacored tx gov vote 1 VOTE_OPTION_NO --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12

zetacored tx gov submit-proposal proposals/proposal_for_success.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12
zetacored tx gov vote 2 VOTE_OPTION_YES --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12

zetacored tx gov submit-proposal proposals/v100.0.0_proposal.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12
zetacored tx gov vote 3 VOTE_OPTION_YES --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12

# Increase the length of the voting period to 1 week
zetacored tx gov submit-proposal proposals/proposal_voting_period.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12
zetacored tx gov vote 4 VOTE_OPTION_YES --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12

# Create a few long lived proposals for variety of testing

echo "Sleeping for 3 minutes to allow the voting period to end and the voting period will be increased to 1 week on future proposals"
sleep 180

zetacored tx gov submit-proposal proposals/proposal_voting_period.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12
zetacored tx gov vote 5 VOTE_OPTION_YES --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12

zetacored tx gov submit-proposal proposals/v100.0.0_proposal.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12
zetacored tx gov vote 6 VOTE_OPTION_YES --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12

zetacored tx gov submit-proposal proposals/proposal_for_deposit.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12
zetacored tx gov vote 7 VOTE_OPTION_YES --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12

# Consensus param test
#zetacored tx gov submit-legacy-proposal param-change proposals/proposal_for_consensus_params.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12
#zetacored tx gov vote 1 VOTE_OPTION_YES --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12

#zetacored tx gov submit-legacy-proposal param-change proposals/emissions_change.json --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --gas 1000000 --yes && sleep 12
#zetacored tx gov vote 1 VOTE_OPTION_YES --from $WALLET_NAME --keyring-backend test --chain-id athens_101-1 --fees 2000000000000000azeta --yes && sleep 12
