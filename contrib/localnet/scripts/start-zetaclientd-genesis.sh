#!/bin/bash

/usr/sbin/sshd

#DISCOVERED_HOSTNAME=$(hostname)
DISCOVERED_HOSTNAME=$(nslookup $(hostname -i) | grep '=' | awk -F'= ' '{split($2, a, "."); print a[1]}')
DISCOVERED_NETWORK=$(echo $DISCOVERED_HOSTNAME |  awk -F'-' '{split($1, a, "-"); print a[1]}')
echo "DEBUG: DISCOVERED_HOSTNAME: $DISCOVERED_HOSTNAME"
echo "DEBUG: DISCOVERED_NETWORK: $DISCOVERED_NETWORK"

cp  /root/preparams/PreParams_$DISCOVERED_HOSTNAME.json /root/preParams.json
num=$(echo $DISCOVERED_HOSTNAME | tr -dc '0-9')
#node="zetacore_node-$num" #unused
#mv  /root/zetacored/zetacored_$node /root/.zetacored
#mv /root/tss/$DISCOVERED_HOSTNAME /root/.tss

echo "Wait for zetacore to exchange genesis file" #TODO: loop instead and actually watch for the file
sleep 30
operator=$(cat $HOME/.zetacored/os.json | jq '.ObserverAddress' )
operatorAddress=$(echo "$operator" | tr -d '"')
echo "operatorAddress: $operatorAddress"
echo "Start zetaclientd"
if [ $DISCOVERED_HOSTNAME == "$DISCOVERED_NETWORK-zetaclient-1" ]
then
    rm ~/.tss/*
    MYIP=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1)
    zetaclientd init  --zetacore-url $DISCOVERED_NETWORK-zetacore_node-1 --chain-id athens_101-1 --operator "$operatorAddress"  --log-format=text --public-ip "$MYIP"
    zetaclientd start
else
  num=$(echo $DISCOVERED_HOSTNAME | tr -dc '0-9')
  node="zetacore_node-$num"
  MYIP=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1)
  #ZETACLIENT_IP=$(dig +short $DISCOVERED_NETWORK-zetaclient-1 | awk '{ print; exit }') #no dig, ya dig? Not in image.
  ZETACLIENT_IP=$(nslookup $DISCOVERED_NETWORK-zetaclient-1 | grep Address: | tail -n1 | awk '{ print $2 }')
  echo "DEBUG: ZETACLIENT_IP is $ZETACLIENT_IP"
  SEED=$(curl --retry 10 --retry-delay 5 --retry-connrefused  -s $DISCOVERED_NETWORK-zetaclient-1:8123/p2p)
  rm ~/.tss/*
  echo "DEBUG: SEED is: $SEED"
  echo "DEBUG: peer is: /ip4/$ZETACLIENT_IP/tcp/6668/p2p/$SEED"
  zetaclientd init --peer /ip4/$ZETACLIENT_IP/tcp/6668/p2p/$SEED --zetacore-url "$DISCOVERED_NETWORK-$node" --chain-id athens_101-1 --operator "$operatorAddress" --log-format=text --public-ip "$MYIP" --log-level 0
  #zetaclientd init --peer /ip4/$ZETACLIENT_IP/tcp/6668/p2p/"$SEED" --zetacore-url "$DISCOVERED_NETWORK-$node" --chain-id athens_101-1 --operator "$operatorAddress" --log-format=text --public-ip "$MYIP" --log-level 0
  #zetaclientd init --peer /dns/$DISCOVERED_NETWORK-zetaclient-1/tcp/6668/p2p/"$SEED" --zetacore-url "$DISCOVERED_NETWORK-$node" --chain-id athens_101-1 --operator "$operatorAddress" --log-format=text --public-ip "$MYIP" --log-level 0
  #echo "DEBUG: peer is: /dns/$DISCOVERED_NETWORK-zetaclient-1/tcp/6668/p2p/$SEED"
  #zetaclientd init --peer /dns/$DISCOVERED_NETWORK-zetaclient-1/tcp/6668/p2p/$SEED --zetacore-url "$DISCOVERED_NETWORK-$node" --chain-id athens_101-1 --operator "$operatorAddress" --log-format=text --public-ip "$MYIP" --log-level 0
  zetaclientd start
fi
