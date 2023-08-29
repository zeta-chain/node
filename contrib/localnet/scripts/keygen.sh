#!/bin/bash
if [ $# -ne 1 ]
then
  echo "Usage: keygen.sh <blocknumber>"
  exit 1
fi 
BLOCKNUM=$1
#DISCOVERED_HOSTNAME=$(hostname)
DISCOVERED_HOSTNAME=$(nslookup $(hostname -i) | grep '=' | awk -F'= ' '{split($2, a, "."); print a[1]}')
DISCOVERED_NETWORK=$(echo $DISCOVERED_HOSTNAME |  awk -F'-' '{split($1, a, "-"); print a[1]}')

cp  /root/preparams/PreParams_$DISCOVERED_HOSTNAME.json /root/preParams.json
num=$(echo $DISCOVERED_HOSTNAME | tr -dc '0-9')
node="zetacore_node-$num"
mv  /root/zetacored/zetacored_$node /root/.zetacored

mv /root/tss/$DISCOVERED_HOSTNAME /root/.tss


if [ $DISCOVERED_HOSTNAME == "$DISCOVERED_NETWORK-zetaclient-1" ]
then
  rm ~/.tss/address_book.seed
  export TSSPATH=~/.tss2
  zetaclientd init --val val --log-console --enable-chains "GOERLI,BSCTESTNET" \
    --pre-params ~/preParams.json --keygen-block $BLOCKNUM --zetacore-url $DISCOVERED_NETWORK-zetacore_node-1 \
    --chain-id athens_101-1
  zetaclientd start
else
  num=$(echo $DISCOVERED_HOSTNAME | tr -dc '0-9')
  node="zetacore_node-$num"
  SEED=$(curl --retry 10 --retry-delay 5 --retry-connrefused  -s $DISCOVERED_NETWORK-zetaclient-1:8123/p2p)

  export TSSPATH=~/.tss2
  zetaclientd init --val val --log-console --enable-chains "GOERLI,BSCTESTNET"  \
    --peer /dns/$DISCOVERED_NETWORK-zetaclient-1/tcp/6668/p2p/$SEED \
    --pre-params ~/preParams.json --keygen-block $BLOCKNUM --zetacore-url $DISCOVERED_NETWORK-$node \
    --chain-id athens_101-1
  zetaclientd start
fi


