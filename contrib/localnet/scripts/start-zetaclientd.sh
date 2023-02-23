#!/bin/bash
HOSTNAME=$(hostname)

cp  /root/preparams/PreParams_$HOSTNAME.json /root/preParams.json
num=$(echo $HOSTNAME | tr -dc '0-9')
node="zetacore$num"
mv  /root/zetacored/zetacored_$node /root/.zetacored

mv /root/tss/$HOSTNAME /root/.tss


if [ $HOSTNAME == "zetaclient0" ]
then
    rm ~/.tss/address_book.seed
    export TSSPATH=~/.tss
    exec zetaclientd -val val -log-console -enable-chains goerli_localnet,btc_regtest \
      -pre-params ~/preParams.json  -zetacore-url zetacore0 \
      -chain-id athens_101-1 -dev
else
  num=$(echo $HOSTNAME | tr -dc '0-9')
  node="zetacore$num"
  SEED=$(curl --retry 10 --retry-delay 5 --retry-connrefused  -s zetaclient0:8123/p2p)

  export TSSPATH=~/.tss
  exec zetaclientd -val val -log-console -enable-chains goerli_localnet,btc_regtest  \
    -peer /ip4/172.20.0.21/tcp/6668/p2p/$SEED \
    -pre-params ~/preParams.json -zetacore-url $node \
    -chain-id athens_101-1 -dev
fi
