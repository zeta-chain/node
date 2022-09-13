#set -e
set -x

CHAINID="athens_8666-2"
KEYALGO="eth_secp256k1"
KEYRING="test"

ZETACORED=/home/ubuntu/go/bin/zetacored
# AWS EC2: testnet1, testnet2, testnet3, testnet5
NODES="3.137.46.147 3.20.194.40 3.19.64.252"

rm -rf ~/.zetacore/data
rm -rf ~/.zetacore/config
for NODE in $NODES; do
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.zetacore/data
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.zetacore/config
done

$ZETACORED init --chain-id ${CHAINID} zetachain --home ~/.zetacore
ADDR=$($ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacore)
if [ -z "$ADDR" ]; then
  echo "No val key found; generate new val key"
  $ZETACORED keys add val --keyring-backend=test --home ~/.zetacore
fi
$ZETACORED add-genesis-account $($ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacore) 100000000000000000000000000azeta --home ~/.zetacore


echo "Generating deterministic account - alice"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | $ZETACORED keys add alice --recover --keyring-backend $KEYRING --home ~/.zetacore

echo "Generating deterministic account - bob"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | $ZETACORED keys add bob --recover --keyring-backend $KEYRING --home ~/.zetacore

$ZETACORED add-genesis-account alice 1000000000000000000000azeta --keyring-backend=test --home ~/.zetacore
$ZETACORED add-genesis-account bob 1000000000000000000000azeta --keyring-backend=test --home ~/.zetacore

for NODE in $NODES; do
  ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacore)
  if [ -z "$ADDR" ]; then
    echo "No val key found; generate new val key"
	  ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys add val --keyring-backend=test --home ~/.zetacore
  fi
	ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacore)
	$ZETACORED add-genesis-account $ADDR 1000000000azeta --keyring-backend=test --home ~/.zetacore
done

 
for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacore/config/genesis.json $NODE:~/.zetacore/config/
done


$ZETACORED gentx val 1000000000stake --keyring-backend=test --chain-id=${CHAINID} --home ~/.zetacore

for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE $ZETACORED gentx val 1000000000azeta --keyring-backend=test --chain-id=${CHAINID} --ip $NODE --home ~/.zetacore
    scp -i ~/.ssh/meta.pem $NODE:~/.zetacore/config/gentx/*.json ~/.zetacore/config/gentx/
done


# Change parameter token denominations to aphoton
cat $HOME/.zetacore/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacore/config/tmp_genesis.json && mv $HOME/.zetacore/config/tmp_genesis.json $HOME/.zetacore/config/genesis.json
cat $HOME/.zetacore/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacore/config/tmp_genesis.json && mv $HOME/.zetacore/config/tmp_genesis.json $HOME/.zetacore/config/genesis.json
cat $HOME/.zetacore/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacore/config/tmp_genesis.json && mv $HOME/.zetacore/config/tmp_genesis.json $HOME/.zetacore/config/genesis.json
cat $HOME/.zetacore/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacore/config/tmp_genesis.json && mv $HOME/.zetacore/config/tmp_genesis.json $HOME/.zetacore/config/genesis.json
cat $HOME/.zetacore/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacore/config/tmp_genesis.json && mv $HOME/.zetacore/config/tmp_genesis.json $HOME/.zetacore/config/genesis.json


# Set gas limit in genesis
cat $HOME/.zetacore/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.zetacore/config/tmp_genesis.json && mv $HOME/.zetacore/config/tmp_genesis.json $HOME/.zetacore/config/genesis.json
sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.zetacore/config/config.toml


$ZETACORED collect-gentxs --home ~/.zetacore
# Run this to ensure everything worked and that the genesis file is setup correctly
$ZETACORED validate-genesis --home ~/.zetacore


for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacore/config/genesis.json $NODE:~/.zetacore/config/
done

#

jq '.chain_id = "athens_8666-2"' ~/.zetacore/config/genesis.json > temp.json && mv temp.json ~/.zetacore/config/genesis.json
sed -i '/\[api\]/,+3 s/enable = false/enable = true/' ~/.zetacore/config/app.toml
sed -i '/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' ~/.zetacore/config/app.toml

for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE jq \'.chain_id = \"athens_8666-2\"\' ~/.zetacore/config/genesis.json > temp.json && mv temp.json ~/.zetacore/config/genesis.json
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+3 s/enable = false/enable = true/\' ~/.zetacore/config/app.toml
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/\' ~/.zetacore/config/app.toml
done

