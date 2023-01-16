#!/bin/bash
HOSTNAME=$(hostname)

if [ $HOSTNAME == "node0" ]
then
  TSSPATH=~/.tss2 zetaclientd -val val -log-console -enable-chains GOERLI,BSCTESTNET -pre-params ~/preParams.json 
else
  SEED=$(cat ~/.seed)
  TSSPATH=~/.tss2 zetaclientd -val val -log-console -enable-chains GOERLI,BSCTESTNET -peer /dns/node0/tcp/6668/p2p/$SEED -pre-params ~/preParams.json 
fi
