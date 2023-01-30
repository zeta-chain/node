#!/bin/bash
HOSTNAME=$(hostname)

if [ $HOSTNAME == "node0" ]
then
  TSSPATH=~/.tss2 zetaclientd -val val -log-console -enable-chains GOERLI,BSCTESTNET -pre-params ~/preParams.json
else
  SEED=$(curl --retry 10 --retry-delay 5 --retry-connrefused  -s node0:8123/p2p)
  echo "SEED:" $SEED
  echo $SEED > ~/.seed
fi
