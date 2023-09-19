#!/bin/bash

clibuilder()
{
   echo ""
   echo "Usage: $0 -u UPGRADE_HEIGHT"
   echo -e "\t-u Height of upgrade, should match governance proposal"
   echo -e "\t-n Number of clients in the network"
   exit 1 # Exit script after printing help
}

while getopts "u:n" opt
do
   case "$opt" in
      u ) UPGRADE_HEIGHT="$OPTARG" ;;
      n ) NUM_OF_NODES="$OPTARG" ;;
      ? ) clibuilder ;; # Print cliBuilder in case parameter is non-existent
   esac
done

# generate node list
START=1
END=$((NUM_OF_NODES - 1))

CLIENT_LIST=()
for i in $(eval echo "{$START..$END}")
do
  CLIENT_LIST+=("zetaclient$i")
done

echo "$UPGRADE_HEIGHT"

CURRENT_HEIGHT=0

while [[ $CURRENT_HEIGHT -lt $UPGRADE_HEIGHT ]]
do
    CURRENT_HEIGHT=$(curl zetacore0:26657/status | jq '.result.sync_info.latest_block_height' | tr -d '"')
    sleep 5
done

echo current height is "$CURRENT_HEIGHT", restarting zetaclients
for NODE in "${NODELIST[@]}"; do
    ssh $NODE "killall zetaclientd; $GOPATH/bin/new/zetaclientd start > $HOME/zetaclient.log 2>&1 &"
done



