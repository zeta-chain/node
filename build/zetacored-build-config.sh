echo "Starting Zetacore"
echo $1 $2 $3

NODE_NUMBER=$1
NODE_0_DNS=$2
 
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin
export MYIP=$(hostname -i)

rm -rf ~/.zetacore/
rm -rf /zetashared/*/data/* /zetashared/*/config/* /zetashared/*/keyring-test/*
rm -rf /zetashared/genesis/*

mkdir -p ~/.zetacore/config/ ~/.zetacore/config/gentx/ ~/.zetacore/keyring-test/

if (( $NODE_NUMBER == 0 )); then
    echo "This is Node $NODE_NUMBER"

    zetacored init testnet --chain-id zetacore
    zetacored config keyring-backend test
    zetacored keys add val
    cd ~/.zetacore/config

    NODE_0_VALIDATOR=$(zetacored keys show val -a)
    echo $NODE_0_VALIDATOR > NODE_VALIDATOR_ID
    zetacored add-genesis-account $NODE_0_VALIDATOR 100000000000stake

    echo "Waiting for other Nodes to Generate Keys"
    until [ -f /zetashared/node1/config/NODE_VALIDATOR_ID ]
    do
        sleep 5
    done
    echo "NODE_1_VALIDATOR_ID found"
    NODE_1_VALIDATOR_ID=$(cat /zetashared/node1/config/NODE_VALIDATOR_ID)
    zetacored add-genesis-account $NODE_1_VALIDATOR_ID 100000000000stake

    # until [ -f /zetashared/node2/config/NODE_VALIDATOR_ID ]
    # do
    #     sleep 5
    # done
    # NODE_2_VALIDATOR_ID=$(cat /zetashared/node2/config/NODE_VALIDATOR_ID)
    # zetacored add-genesis-account $NODE_2_VALIDATOR_ID 100000000000stake

    cp ~/.zetacore/config/genesis.json /zetashared/genesis/init-genesis.json

    echo "Waiting for other Nodes to generate gentx files"
    until [ -f /zetashared/node1/config/gentx/gentx-*.json ]
    do
        sleep 5
    done
    cp /zetashared/node1/config/gentx/gentx-*.json ~/.zetacore/config/gentx/
    zetacored gentx val 100000000stake --chain-id zetacore --ip $MYIP 

    # Node 2 
    # cp /zetashared/node2/config/gentx/gentx-*.json ~/.zetacore/config/gentx/
    # NODE_2_VALIDATOR_ID=$(cat /zetashared/node2/config/NODE_VALIDATOR_ID)

    zetacored collect-gentxs &> gentxs

    # echo $(cat gentxs | jq -r .node_id) > NODE_0_ID
    # cp ~/.zetacore/config/NODE_0_ID /zetashared/

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
    echo $NODE_VALIDATOR > NODE_VALIDATOR_ID
    cp NODE_VALIDATOR_ID /zetashared/node$NODE_NUMBER/config/

    echo "Waiting for Node 0 to Create Genesis..."

    until [ -f /zetashared/genesis/init-genesis.json ]
    do
        sleep 5
    done
    echo "init-genesis.json found"
    # This needs to happen after Node 0 creates the init-genesis file but before it runs collect-gentxs

    cp /zetashared/genesis/init-genesis.json  ~/.zetacore/config/genesis.json 
    zetacored gentx val 1000000000stake --chain-id zetacore --ip $MYIP

    cp -r /root/.zetacore/config/* /zetashared/node$NODE_NUMBER/config/
    cp -r /root/.zetacore/keyring-test/* /zetashared/node$NODE_NUMBER/keyring-test/
    cp -r /root/.zetacore/data/* /zetashared/node$NODE_NUMBER/data/

    echo "Waiting for updated genesis file"  # Should be loop, waiting for the file to show up
    until [ -f /zetashared/genesis/genesis.json ]
    do
        sleep 5
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
