#!/bin/bash

# This script is used to restart zetaclientd after an upgrade
# It waits for the upgrade height to be reached and then restarts the zetaclientd on all nodes in the network
# It interacts with the network using the zetaclientd binary

clibuilder()
{
   echo ""
   echo "Usage: $0 -u UPGRADE_HEIGHT"
   echo -e "\t-u Height of upgrade, should match governance proposal"
   echo -e "\t-n Number of clients in the network"
   exit 1 # Exit script after printing help
}

while getopts "u:n:" opt
do
   case "$opt" in
      u ) UPGRADE_HEIGHT="$OPTARG" ;;
      n ) NUM_OF_NODES="$OPTARG" ;;
      ? ) clibuilder ;; # Print cliBuilder in case parameter is non-existent
   esac
done

# generate client list
START=0
END=$((NUM_OF_NODES-1))
CLIENT_LIST=()
for i in $(eval echo "{$START..$END}")
do
  CLIENT_LIST+=("zetaclient$i")
done

echo "$UPGRADE_HEIGHT"

CURRENT_HEIGHT=0

while [[ $CURRENT_HEIGHT -lt $UPGRADE_HEIGHT ]]
do
    CURRENT_HEIGHT=$(curl -s zetacore0:26657/status | jq '.result.sync_info.latest_block_height' | tr -d '"')
    echo current height is "$CURRENT_HEIGHT", waiting for "$UPGRADE_HEIGHT"
    sleep 5
done

echo upgrade height reached, restarting zetaclients

for NODE in "${CLIENT_LIST[@]}"; do
    ssh -o "StrictHostKeyChecking no" "$NODE" -i ~/.ssh/localtest.pem killall zetaclientd
    ssh -o "StrictHostKeyChecking no" "$NODE" -i ~/.ssh/localtest.pem "$GOPATH/bin/new/zetaclientd start < /root/password.file > $HOME/zetaclient.log 2>&1 &"
done
