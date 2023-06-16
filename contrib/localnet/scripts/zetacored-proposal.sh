#!/bin/bash
echo "proposal script running..."
sleep 100

HOSTNAME=$(hostname)
num=$(echo $HOSTNAME | tr -dc '0-9')

cd /root

if [ "$num" == "0" ]
then
  zetacored tx gov submit-proposal draft_proposal.json --from operator --chain-id athens_101-1 --fees 20azeta --yes
fi

sleep 10
zetacored tx gov vote 1 yes --from operator --keyring-backend test --chain-id athens_101-1 --yes --fees 20azeta --yes


