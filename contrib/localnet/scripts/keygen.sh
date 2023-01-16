#!/bin/bash
if [ $# -ne 1 ]
then
  echo "Usage: keygen.sh <blocknumber>"
  exit 1
fi 
BLOCKNUM=$1
HOSTNAME=$(hostname)

if [ $HOSTNAME == "node0" ]
then
  TSSPATH=~/.tss2 zetaclientd -val val -log-console -enable-chains GOERLI,BSCTESTNET -pre-params ~/preParams.json -keygen-block $BLOCKNUM
else
  SEED=$(curl --retry 10 --retry-delay 5 --retry-connrefused  -s node0:8123/p2p)
  TSSPATH=~/.tss2 zetaclientd -val val -log-console -enable-chains GOERLI,BSCTESTNET  -peer /dns/node0/tcp/6668/p2p/$SEED -pre-params ~/preParams.json -keygen-block $BLOCKNUM
fi
