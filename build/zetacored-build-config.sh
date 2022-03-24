#!/bin/bash

echo "Starting Zetacore"
echo $1 $2 $3

NODE_NUMBER=$1
MAX_NODE_NUMBER=$2 #Whats the highest node number? If you have nodes 0,1,2,3 MAX_NODE_NUMBER=3
echo "MAX_NODE_NUMBER: $MAX_NODE_NUMBER"

export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin
export MYIP=$(hostname -i)

rm -rf ~/.zetacore/

mkdir -p ~/.zetacore/data/ ~/.zetacore/config/gentx/ ~/.zetacore/keyring-test/

if (( $NODE_NUMBER == 0 )); then
    echo "This is Node $NODE_NUMBER"

    zetacored init testnet --chain-id zetacore
    zetacored config keyring-backend test
    zetacored keys add val
    cd ~/.zetacore/config

    NODE_0_VALIDATOR=$(zetacored keys show val -a)
    echo $NODE_0_VALIDATOR > NODE_VALIDATOR_ID
    zetacored add-genesis-account $NODE_0_VALIDATOR 100000000000stake

    i=1
    while [ $i -le $MAX_NODE_NUMBER ]
    do
        echo "i = $i"
        until [ -f /zetashared/node$i/config/NODE_VALIDATOR_ID ]
            echo "Waiting for Node $i to generate new keys"
            do
                sleep 3
            done
        echo "VALIDATOR_ID for node$i found"
        VALIDATOR_ID=$(cat /zetashared/node$i/config/NODE_VALIDATOR_ID)
        echo "Node $i VALIDATOR_ID: $VALIDATOR_ID"
        zetacored add-genesis-account $VALIDATOR_ID 100000000000stake
        i=$[$i+1]
    done

    cp ~/.zetacore/config/genesis.json /zetashared/genesis/init-genesis.json

    i=1
    while [ $i -le $MAX_NODE_NUMBER ]
    do
        echo "i = $i"
        until [ -f /zetashared/node$i/config/gentx/gentx-*.json ]
            do
                echo "Waiting for Node $i to generate gentx files"
                sleep 3
            done
        cp /zetashared/node$i/config/gentx/gentx-*.json ~/.zetacore/config/gentx/
        i=$[$i+1]
    done

    zetacored gentx val 100000000stake --chain-id zetacore --ip $MYIP 
    zetacored collect-gentxs &> gentxs

    cp /root/.zetacore/config/genesis.json /zetashared/genesis/genesis.json
    cp -r /root/.zetacore/config/* /zetashared/node$NODE_NUMBER/config/
    cp -r /root/.zetacore/data/* /zetashared/node$NODE_NUMBER/data/
    cp -r /root/.zetacore/keyring-test/* /zetashared/node$NODE_NUMBER/keyring-test/

   echo "Config Built -- Node $NODE_NUMBER"

fi

if (( $NODE_NUMBER > 0 )); then
    echo "This is Node $NODE_NUMBER"
    echo "Generating new keys"
    zetacored config keyring-backend test
    zetacored keys add val
    NODE_VALIDATOR=$(zetacored keys show val -a)
    echo $NODE_VALIDATOR
    echo $NODE_VALIDATOR > NODE_VALIDATOR_ID
    cp NODE_VALIDATOR_ID /zetashared/node$NODE_NUMBER/config/

    echo "Waiting for Node 0 to Create Genesis..."

    until [ -f /zetashared/genesis/init-genesis.json ]
        do
            sleep 3
        done
    echo "init-genesis.json found"

    sleep 5 # Can probably be removed

    # Happens after Node 0 creates the init-genesis file but before it runs collect-gentxs
    cp /zetashared/genesis/init-genesis.json  ~/.zetacore/config/genesis.json 
    zetacored gentx val 100000000stake --chain-id zetacore --ip $MYIP 
    cp -r /root/.zetacore/config/* /zetashared/node$NODE_NUMBER/config/
    cp -r /root/.zetacore/keyring-test/* /zetashared/node$NODE_NUMBER/keyring-test/
    cp -r /root/.zetacore/data/* /zetashared/node$NODE_NUMBER/data/

    until [ -f /zetashared/genesis/genesis.json ]
        do
            echo "Waiting for updated genesis file..."
            sleep 3
        done
    echo "Final genesis.json found"
    cp /zetashared/genesis/genesis.json  ~/.zetacore/config/genesis.json 
    cp -r /root/.zetacore/config/* /zetashared/node$NODE_NUMBER/config/

    echo "Config Built -- Node $NODE_NUMBER"

fi



## Manual Steps for setting up Node 1+ the first time 
    # export MYIP=$(hostname -i)
   
    # mkdir -p ~/.zetacore/config/ ~/.zetacore/config/gentx/ ~/.zetacore/keyring-test/
    #  cp /zetashared/genesis/genesis.json  ~/.zetacore/config/genesis.json 
    #  zetacored config keyring-backend test
    #  zetacored keys add val # or  cp -r /zetashared/node$NODE_NUMBER/keyring-test/* ~/.zetacore/keyring-test/
    # NODE_VALIDATOR=$(zetacored keys show val -a)
    # zetacored add-genesis-account $NODE_VALIDATOR 100000000000stake
    # MYIP=$(hostname -i)
    # zetacored gentx val 1000000000stake --node "tcp://$MYIP:26657" --ip $MYIP --chain-id zetacore
    # cp -r /root/.zetacore/config/* /zetashared/node$NODE_NUMBER/config/
    # cp -r /root/.zetacore/keyring-test/* /zetashared/node$NODE_NUMBER/keyring-test/
    # cp -r /root/.zetacore/data/* /zetashared/node$NODE_NUMBER/data/

    # cp -r /root/.zetacore/config/gentx/gentx-*.json /zetashared/node$NODE_NUMBER/config/gentx/
    # cp /root/.zetacore/config/node_key.json /zetashared/node$NODE_NUMBER/config/node_key.json
    # cp /root/.zetacore/config/priv_validator_key.json /zetashared/node$NODE_NUMBER/config/priv_validator_key.json
