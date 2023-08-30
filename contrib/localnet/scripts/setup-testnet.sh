#!/bin/bash
set -x
ZETACORED=/usr/local/bin/zetacored
NODES="zetacore-2"
#DISCOVERED_HOSTNAME=$(hostname)
DISCOVERED_HOSTNAME=$(nslookup $(hostname -i) | grep '=' | awk -F'= ' '{split($2, a, "."); print a[1]}')
DISCOVERED_NETWORK=$(echo $DISCOVERED_HOSTNAME |  awk -F'-' '{split($1, a, "-"); print a[1]}')
if [ $DISCOVERED_HOSTNAME != "$DISCOVERED_NETWORK-zetacore-1" ]
then
  echo "You should run this only on $DISCOVERED_NETWORK-zetacore-1."
  exit 1
fi

if $ZETACORED validate-genesis; then
  echo "Genesis file is valid"
else
  echo "Genesis file is invalid"
  exit 1
fi

NODES="zetacore-1 zetacore-2"
for NODE in $NODES; do
  ssh  $DISCOVERED_NETWORK-$NODE $ZETACORED validate-genesis
  scp
done

