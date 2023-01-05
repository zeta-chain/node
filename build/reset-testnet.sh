#set -e
set -x
ZETACORED=/home/ubuntu/go/bin/zetacored
# AWS EC2: testnet1, testnet2, testnet3, testnet5
NODES="3.137.46.147 3.20.194.40 3.19.64.252"

observer=()

rm -rf ~/.zetacored/data
rm -rf ~/.zetacored/config
for NODE in $NODES; do
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.zetacored/data
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.zetacored/config
done

$ZETACORED init --chain-id athens_101-1 zetachain
ADDR=$($ZETACORED keys show val -a --keyring-backend=test)
observer+=$ADDR
observer+=","

if [ -z "$ADDR" ]; then
  echo "No val key found; generate new val key"
  $ZETACORED keys add val --keyring-backend=test --algo=secp256k1
fi
$ZETACORED add-genesis-account $($ZETACORED keys show val -a --keyring-backend=test) 100000000000000000000000000azeta

# give test address bob 10000 ZETA
# cosmos:zeta1h4m2lf04kpzn4c6fj7tfcnr29fmgk2vfpkjma6
# hex: 0xBD76aFa5f5b0453aE34997969C4c6a2A768b2989
$ZETACORED add-genesis-account zeta1h4m2lf04kpzn4c6fj7tfcnr29fmgk2vfpkjma6 10000000000000000000000azeta


for NODE in $NODES; do
  ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys show val -a --keyring-backend=test)
  if [ -z "$ADDR" ]; then
    echo "No val key found; generate new val key"
	  ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys add val --keyring-backend=test --algo=secp256k1
  fi
	ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys show val -a --keyring-backend=test)
	observer+=$ADDR
	observer+=","
	$ZETACORED add-genesis-account $ADDR 100000000000000000000000000azeta --keyring-backend=test
done

observer_list=$(echo $observer | rev | cut -c2- | rev)
zetacored add-observer Goerli InBoundTx "$observer_list" #goerli
zetacored add-observer Goerli OutBoundTx "$observer_list"
zetacored add-observer BscTestnet InBoundTx "$observer_list" #bsctestnet
zetacored add-observer BscTestnet OutBoundTx "$observer_list"
zetacored add-observer Mumbai InBoundTx "$observer_list" #mumbai
zetacored add-observer Mumbai OutBoundTx "$observer_list"
zetacored add-observer BTCTestnet InBoundTx "$observer_list" #btctestnet
zetacored add-observer BTCTestnet OutBoundTx "$observer_list"
zetacored add-observer Baobab InBoundTx "$observer_list" #baobab klaytn
zetacored add-observer Baobab OutBoundTx "$observer_list"
 
for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacored/config/genesis.json $NODE:~/.zetacored/config/
done


$ZETACORED gentx val 10000000000000000000000000azeta --keyring-backend=test --chain-id=athens_101-1

for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE $ZETACORED gentx val 10000000000000000000000000azeta --keyring-backend=test --chain-id=athens_101-1 --ip $NODE
    scp -i ~/.ssh/meta.pem $NODE:~/.zetacored/config/gentx/*.json ~/.zetacored/config/gentx/
done


$ZETACORED collect-gentxs

# Change parameter token denominations to aphoton
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json

# set block gas limit to 100 million
cat $HOME/.zetacored/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="100000000"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json

zetacored validate-genesis --home ~/.zetacored


for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacored/config/genesis.json $NODE:~/.zetacored/config/
done

#

jq '.chain_id = "athens_101-1"' ~/.zetacored/config/genesis.json > temp.json && mv temp.json ~/.zetacored/config/genesis.json
sed -i '/\[api\]/,+3 s/enable = false/enable = true/' ~/.zetacored/config/app.toml
sed -i '/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' ~/.zetacored/config/app.toml

for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE jq \'.chain_id = \"athens_101-1\"\' ~/.zetacored/config/genesis.json > temp.json && mv temp.json ~/.zetacored/config/genesis.json
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+3 s/enable = false/enable = true/\' ~/.zetacored/config/app.toml
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/\' ~/.zetacored/config/app.toml
done

