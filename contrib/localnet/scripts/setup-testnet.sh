#!/bin/bash
set -x
ZETACORED=/usr/local/bin/zetacored
NODES="zetacore1"
HOSTNAME=$(hostname)
if [ $HOSTNAME != "zetacore0" ]
then
  echo "You should run this only on zetacore0."
  exit 1
fi

if $ZETACORED validate-genesis; then
  echo "Genesis file is valid"
else
  echo "Genesis file is invalid"
  exit 1
fi

NODES="zetacore0 zetacore1"
for NODE in $NODES; do
  ssh  $NODE $ZETACORED validate-genesis
  scp
done

