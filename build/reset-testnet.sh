#set -e
set -x
ZETACORED=/home/ubuntu/go/bin/zetacored
# AWS EC2: testnet1, testnet2, testnet3, testnet5
NODES="3.137.46.147 3.20.194.40 3.19.64.252"

rm -rf ~/.zetacore/data
rm -rf ~/.zetacore/config
for NODE in $NODES; do
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.zetacore/data
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.zetacore/config
done

$ZETACORED init --chain-id athens-1 zetachain --home ~/.zetacore
ADDR=$($ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacore)
if [ -z "$ADDR" ]; then
  echo "No val key found; generate new val key"
  $ZETACORED keys add val --keyring-backend=test --home ~/.zetacore
fi
$ZETACORED add-genesis-account $($ZETACORED keys show val -a --keyring-backend=test) 1000000000stake --home ~/.zetacore


for NODE in $NODES; do
  ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacore)
  if [ -z "$ADDR" ]; then
    echo "No val key found; generate new val key"
	  ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys add val --keyring-backend=test --home ~/.zetacore
  fi
	ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacore)
	$ZETACORED add-genesis-account $ADDR 1000000000stake --keyring-backend=test --home ~/.zetacore
done

 
for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacore/config/genesis.json $NODE:~/.zetacore/config/
done


$ZETACORED gentx val 1000000000stake --keyring-backend=test --chain-id=athens-1 --home ~/.zetacore

for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE $ZETACORED gentx val 1000000000stake --keyring-backend=test --chain-id=athens-1 --ip $NODE --home ~/.zetacore
    scp -i ~/.ssh/meta.pem $NODE:~/.zetacore/config/gentx/*.json ~/.zetacore/config/gentx/
done


$ZETACORED collect-gentxs --home ~/.zetacore


for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacore/config/genesis.json $NODE:~/.zetacore/config/
done

#

jq '.chain_id = "athens-1"' ~/.zetacore/config/genesis.json > temp.json && mv temp.json ~/.zetacore/config/genesis.json
sed -i '/\[api\]/,+3 s/enable = false/enable = true/' ~/.zetacore/config/app.toml
sed -i '/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' ~/.zetacore/config/app.toml

for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE jq \'.chain_id = \"athens-1\"\' ~/.zetacore/config/genesis.json > temp.json && mv temp.json ~/.zetacore/config/genesis.json
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+3 s/enable = false/enable = true/\' ~/.zetacore/config/app.toml
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/\' ~/.zetacore/config/app.toml
done

