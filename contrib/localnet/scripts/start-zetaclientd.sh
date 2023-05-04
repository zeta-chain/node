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
    zetaclientd init \
      --pre-params ~/preParams.json  --zetacore-url zetacore0 \
      --chain-id athens_101-1 --operator zeta1z46tdw75jvh4h39y3vu758ctv34rw5z9kmyhgz --log-level 1 --hotkey=val_grantee_observer
    zetaclientd start
else
  num=$(echo $HOSTNAME | tr -dc '0-9')
  node="zetacore$num"
  SEED=$(curl --retry 10 --retry-delay 5 --retry-connrefused  -s zetaclient0:8123/p2p)
  zetaclientd init \
    --peer /ip4/172.20.0.21/tcp/6668/p2p/$SEED \
    --pre-params ~/preParams.json --zetacore-url $node \
    --chain-id athens_101-1 --operator zeta1lz2fqwzjnk6qy48fgj753h48444fxtt7hekp52 --log-level 0 --hotkey=val_grantee_observer
  zetaclientd start
fi
