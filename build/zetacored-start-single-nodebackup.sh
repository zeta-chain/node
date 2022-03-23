echo "Starting Zetacore"
echo $1 $2 $3

NODE_NUMBER=$1
NODE_0_DNS=$2
NODE_0_ID=$3
 
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin
export MYIP=$(hostname -i)


DIR="/root/.zetacore/config/gentx"
if  (( $NODE_NUMBER == 0 )) && [[ -d "$DIR" ]]; then
    echo "This is Node $NODE_NUMBER"
    echo "$DIR already exists."
    echo "Skipping ZetaCore Init"
    zetacored start --rpc.laddr "tcp://0.0.0.0:26657" --proxy_app "tcp://0.0.0.0:26658" --rpc.pprof_laddr "0.0.0.0:6060" 2>&1 | tee /root/.zetacore/zetacored.log

elif (( $NODE_NUMBER == 0 )); then
    echo "This is Node $NODE_NUMBER"
    # rm -rf ~/.zetacore/* # zetacored stores all states in this directory
    zetacored init testnet
    cd ~/.zetacore/config
    zetacored config keyring-backend test
    zetacored keys add val
    MY_VALIDATOR_ADDRESS=$(zetacored keys show val -a)
    zetacored add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake #--node "tcp://0.0.0.0:26657"
    zetacored gentx val 100000000stake --chain-id zetacore #--node "tcp://0.0.0.0:26657"
    zetacored collect-gentxs &> gentxs
    export NODE_ID=$(cat gentxs | jq -r .node_id)
    echo $NODE_ID > NODE_0_ID
    cp ~/.zetacore/config/genesis.json /zetashared/
    cp ~/.zetacore/config/NODE_0_ID /zetashared/
    cp ~/.zetacore/config/gentx /zetashared/
    zetacored start --rpc.laddr "tcp://0.0.0.0:26657" --proxy_app "tcp://0.0.0.0:26658" --rpc.pprof_laddr "0.0.0.0:6060" 2>&1 | tee /root/.zetacore/zetacored.log
fi

DIR="/root/.zetacore/config/gentx"
if  (( $NODE_NUMBER > 0 )) && [[ -d "$DIR" ]]; then
    echo "This is Node $NODE_NUMBER"
    echo "$DIR already exists."
    echo "Skipping ZetaCore Init"
    # zetacored start --rpc.laddr "tcp://0.0.0.0:26657" --proxy_app "tcp://0.0.0.0:26658" --rpc.pprof_laddr "0.0.0.0:6060" 2>&1 | tee /root/.zetacore/zetacored.log

elif (( $NODE_NUMBER > 0 )); then
    echo "This is Node $NODE_NUMBER"
    # rm -rf ~/.zetacore/* # zetacored stores all states in this directory
    # zetacored init testnet
    # cd ~/.zetacore/config
    
    zetacored config keyring-backend test
    # zetacored keys add val
    # MY_VALIDATOR_ADDRESS=$(zetacored keys show val -a)
    # zetacored add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake #--node "tcp://0.0.0.0:26657"
    # zetacored gentx val 100000000stake --chain-id zetacore #--node "tcp://0.0.0.0:26657"
    # zetacored collect-gentxs &> gentxs
    # export NODE_ID=$(cat gentxs | jq -r .node_id)
    # echo $NODE_ID > NODE_0_ID
    cp /zetashared/genesis.json  ~/.zetacore/config/genesis.json 
    NODE_0_ID=$(cat /zetashared/NODE_0_ID)
    cp /zetashared/gentx ~/.zetacore/config/gentx 

    zetacored start --rpc.laddr "tcp://0.0.0.0:26657" --proxy_app "tcp://0.0.0.0:26658" --rpc.pprof_laddr "0.0.0.0:6060" 2>&1 | tee /root/.zetacore/zetacored.log

fi


