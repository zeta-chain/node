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

echo "$UPGRADE_HEIGHT"

CURRENT_HEIGHT=0

while [ $CURRENT_HEIGHT -lt "$UPGRADE_HEIGHT" ]
do
    CURRENT_HEIGHT=$(curl localhost:26657/status | jq '.result.sync_info.latest_block_height' )
done

for i in {$NUM_OF_NODES}
do
    ssh "zetaclient$i" "killall zetaclientd; zetaclientd start > $HOME/zetaclient.log 2>&1 &"
done



