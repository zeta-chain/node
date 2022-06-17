#!/bin/bash

NODE_NUMBER=$1
MAX_NODE_NUMBER=$2 #Whats the highest node number? If you have nodes 0,1,2,3 MAX_NODE_NUMBER=3
REUSE_EXISTING_KEYS=$3
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin

if [ -z "${MYIP}" ]; then 
    echo "MYIP ENV Variable Not Set -- Setting it automatically using host IP"
    export MYIP=$(hostname -i)
fi
echo "MYIP: $MYIP"
echo "MyLocalIP: $(hostname -i)"

REUSE_EXISTING_KEYS="true"
if [ -z "${REUSE_EXISTING_KEYS}" ]; then 
    echo "Generating new keys"
    rm -rf ~/.zetacore/
    rm -rf ~/.tssnew/
    rm -rf ~/.tss/
    rm -rf ~/.zetaclient/
elif [ "${REUSE_EXISTING_KEYS}" == "true" ]; then 
    echo "Reusing existing keys"
    rm -rf ~/.zetaclient/
    rm -rf ~/.zetacore/data
    rm -rf ~/.zetacore/config
    cp ~/backup-2022-06-17/zetacore/keyring-test/* ~/.zetacore/keyring-test/ ## Temporary to restore exisiting keys from athens
else
    echo "Unknown Input -- REUSE_EXISTING_KEYS=$REUSE_EXISTING_KEYS"
    exit 1
fi



echo "testing"
echo "ls ~/"
ls ~/
# echo "ls ~/.zetacore/keyring-test/"
# ls ~/.zetacore/keyring-test/
# echo "ls ~/.zetacore/*"
# ls ~/.zetacore/*
# echo "ls /zetashared/"
# ls /zetashared/

if (( $NODE_NUMBER == 0 )); then
    echo "This is Node $NODE_NUMBER"
    rm -r /zetashared/*
    mkdir -p /zetashared/genesis/ /zetashared/node"${NODE_NUMBER}"/config/gentx/ /zetashared/node"${NODE_NUMBER}"/data/ /zetashared/node"${NODE_NUMBER}"/keyring-test/

    zetacored init --chain-id athens-1 zetachain
    zetacored config keyring-backend test
    # zetacored keys add val
    cd ~/.zetacore/config || exit
    NODE_0_VALIDATOR=$(zetacored keys show val -a)
    echo "NODE_0_VALIDATOR: $NODE_0_VALIDATOR"
    echo "$NODE_0_VALIDATOR" > NODE_VALIDATOR_ID
    zetacored add-genesis-account "$NODE_0_VALIDATOR" 100000000000stake

    if [ "$STAKER_ACCOUNT_MEMONIC" != "" ]; then
        echo "$STAKER_ACCOUNT_MEMONIC"
        echo "CREATING STAKE ACCOUNT WITH 1000000000000000000000000stake"
        # echo "hip stick bless tank flame raw basket solution deposit share must rookie harbor warfare method joke cram umbrella clump they wasp notice blind empower" | zetacored keys add staker --recover
        STAKER_ADDR=$(zetacored keys show staker -a)
        echo "STAKER ADDR: $STAKER_ADDR"
        zetacored add-genesis-account "$STAKER_ADDR" 1000000000000000000000000stake
    fi

    i=1
    while [ $i -le "$MAX_NODE_NUMBER" ]
    do
        until [ -f /zetashared/node$i/config/NODE_VALIDATOR_ID ]
            echo "Waiting for Node $i to generate new keys"
            do
                sleep 3
            done
        echo "VALIDATOR_ID for node$i found"
        VALIDATOR_ID=$(cat /zetashared/node$i/config/NODE_VALIDATOR_ID)
        echo "Node $i VALIDATOR_ID: $VALIDATOR_ID"
        zetacored add-genesis-account "$VALIDATOR_ID" 100000000000stake
        i=$[$i+1]
    done

    cp ~/.zetacore/config/genesis.json /zetashared/genesis/init-genesis.json
    
    i=1
    while [ $i -le "$MAX_NODE_NUMBER" ]
    do
        # echo "i = $i"
        until [ -f /zetashared/node$i/config/gentx/gentx-*.json ]
            do
                echo "Waiting for Node $i to generate gentx files"
                sleep 3
            done
        cp /zetashared/node$i/config/gentx/gentx-*.json ~/.zetacore/config/gentx/
        i=$[$i+1]
    done
    zetacored gentx val 100000000stake --chain-id athens-1 --ip "$MYIP" --moniker "node$NODE_NUMBER" 
    zetacored collect-gentxs &> gentxs

    sed -i '/\[instrumentation\]/,+3 s/prometheus = false/prometheus = true/' /root/.zetacore/config/config.toml
    sed -i '/\[instrumentation\]/,+3 s/namespace = "tendermint"/namespace = "zetachain-athens"/' /root/.zetacore/config/config.toml
    sed -i '/\[telemetry\]/,+6 s/enabled = false/enabled = false/' /root/.zetacore/config/app.toml
    sed -i '/\[api\]/,+3 s/enable = false/enable = true/' /root/.zetacore/config/app.toml
    sed -i 's/enable-hostname-label = false/enable-hostname-label = true/' /root/.zetacore/config/app.toml
    sed -i 's/prometheus-retention-time = 5/prometheus-retention-time = 5/' /root/.zetacore/config/app.toml

    cp /root/.zetacore/config/genesis.json /zetashared/genesis/genesis.json
    cp -r /root/.zetacore/config/* /zetashared/node"$NODE_NUMBER"/config/
    cp -r /root/.zetacore/data/* /zetashared/node"$NODE_NUMBER"/data/
    cp -r /root/.zetacore/keyring-test/* /zetashared/node"$NODE_NUMBER"/keyring-test/

   echo "Config Built -- Node $NODE_NUMBER"

fi

if (( $NODE_NUMBER > 0 )); then
    echo "This is Node $NODE_NUMBER"
    sleep 20
    mkdir -p /zetashared/node"${NODE_NUMBER}"/config/gentx/ /zetashared/node"${NODE_NUMBER}"/data/ /zetashared/node"${NODE_NUMBER}"/keyring-test/

    echo "Generating new keys"
    zetacored config keyring-backend test
    # zetacored keys add val
    NODE_VALIDATOR=$(zetacored keys show val -a)
    echo "NODE_VALIDATOR: $NODE_VALIDATOR"
    echo "$NODE_VALIDATOR" > NODE_VALIDATOR_ID
    cp NODE_VALIDATOR_ID /zetashared/node"$NODE_NUMBER"/config/

    echo "Waiting for Node 0 to Create Genesis..."

    until [ -f /zetashared/genesis/init-genesis.json ]
        do
            sleep 3
        done
    echo "init-genesis.json found"

    sleep 5 # Wait to make sure node0 has finished configuring the genesis file

    # Happens after Node 0 creates the init-genesis file but before it runs collect-gentxs
    cp /zetashared/genesis/init-genesis.json  ~/.zetacore/config/genesis.json 
    zetacored gentx val 100000000stake --chain-id athens-1 --ip "$MYIP" --moniker "node$NODE_NUMBER" 

    sed -i '/\[instrumentation\]/,+3 s/prometheus = false/prometheus = true/' /root/.zetacore/config/config.toml
    sed -i '/\[instrumentation\]/,+3 s/namespace = "tendermint"/namespace = "zetachain-athens"/' /root/.zetacore/config/config.toml
    sed -i '/\[telemetry\]/,+6 s/enabled = false/enabled = false/' /root/.zetacore/config/app.toml
    sed -i '/\[api\]/,+3 s/enable = false/enable = true/' /root/.zetacore/config/app.toml
    sed -i 's/enable-hostname-label = false/enable-hostname-label = true/' /root/.zetacore/config/app.toml
    sed -i 's/prometheus-retention-time = 5/prometheus-retention-time = 5/' /root/.zetacore/config/app.toml

    cp -r /root/.zetacore/config/* /zetashared/node"$NODE_NUMBER"/config/
    cp -r /root/.zetacore/keyring-test/* /zetashared/node"$NODE_NUMBER"/keyring-test/
    cp -r /root/.zetacore/data/* /zetashared/node"$NODE_NUMBER"/data/


    until [ -f /zetashared/genesis/genesis.json ]
        do
            echo "Waiting for updated genesis file..."
            sleep 3
        done
    # echo "Final genesis.json found"
    sleep 5 
    cp /zetashared/genesis/genesis.json  ~/.zetacore/config/genesis.json 
    cp -r /root/.zetacore/config/* /zetashared/node"$NODE_NUMBER"/config/

    echo "Config Built -- Node $NODE_NUMBER"

fi
