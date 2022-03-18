
FILE="/root/.zetacore/config"
export MYIP=$(hostname -i)
if [[ -d "$FILE" ]]; then
    echo "$FILE already exists."
    echo "Skipping ZetaCore Init"
else
    rm -rf ~/.zetacore/* # zetacored stores all states in this directory
    zetacored init mocknet
    cd ~/.zetacore/config
    zetacored config keyring-backend test
    zetacored keys add val
    MY_VALIDATOR_ADDRESS=$(zetacored keys show val -a)
    zetacored add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake
    zetacored gentx val 100000000stake --chain-id zetacore
    zetacored collect-gentxs &> gentxs
    export NODE_ID=$(cat gentxs | jq -r .node_id)
    echo $NODE_ID > NODE_ID
fi
zetacored start 2>&1 | tee /root/.zetacore/zetacored.log
