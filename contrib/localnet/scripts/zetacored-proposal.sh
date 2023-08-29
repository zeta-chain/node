#!/bin/bash
echo "proposal script running..."
sleep 100

#DISCOVERED_HOSTNAME=$(hostname)
DISCOVERED_HOSTNAME=$(nslookup $(hostname -i) | grep '=' | awk -F'= ' '{split($2, a, "."); print a[1]}')
DISCOVERED_NETWORK=$(echo $DISCOVERED_HOSTNAME |  awk -F'-' '{split($1, a, "-"); print a[1]}')
num=$(echo $DISCOVERED_HOSTNAME | tr -dc '0-9')

cd /root

if [ "$num" == "0" ]
then
  zetacored tx gov submit-proposal draft_proposal.json --from operator --chain-id athens_101-1 --fees 20azeta --yes
fi

sleep 10
zetacored tx gov vote 1 yes --from operator --keyring-backend test --chain-id athens_101-1 --yes --fees 20azeta --yes


