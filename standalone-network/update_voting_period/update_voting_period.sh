#!/bin/bash
CHAINID="atheens_7001-1"
KEYRING="test"
HOSTNAME=$(hostname)
signer="operator"
proposal_count=10

#PID=1

signerAddress=$(zetacored keys show $signer -a --keyring-backend=test)
echo "signerAddress: $signerAddress"
for (( i = 0; i < proposal_count; i++ )); do
  zetacored tx gov submit-legacy-proposal param-change proposal.json --from $signer --gas=auto --gas-adjustment=1.5 --gas-prices=0.001azeta --chain-id=$CHAINID --keyring-backend=$KEYRING -y --broadcast-mode=block
done


#zetacored tx gov vote "$PID" yes --from $signer --keyring-backend $KEYRING --chain-id $CHAINID --yes --fees=40azeta --broadcast-mode=block