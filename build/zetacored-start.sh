
DIR="/root/.zetacore/config/gentx"
export MYIP=$(hostname -i)
if [[ -d "$DIR" ]]; then
    echo "$DIR already exists."
    echo "Skipping ZetaCore Init"
else
    # rm -rf ~/.zetacore/* # zetacored stores all states in this directory
    zetacored init mocknet
    cd ~/.zetacore/config
    zetacored config keyring-backend test
    zetacored keys add val
    MY_VALIDATOR_ADDRESS=$(zetacored keys show val -a)
    zetacored add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake #--node "tcp://0.0.0.0:26657"
    zetacored gentx val 100000000stake --chain-id zetacore #--node "tcp://0.0.0.0:26657"
    zetacored collect-gentxs &> gentxs
    export NODE_ID=$(cat gentxs | jq -r .node_id)
    echo $NODE_ID > NODE_ID
fi

zetacored start --rpc.laddr "tcp://0.0.0.0:26657" --proxy_app "tcp://0.0.0.0:26658" --rpc.pprof_laddr "0.0.0.0:6060" 2>&1 | tee /root/.zetacore/zetacored.log
