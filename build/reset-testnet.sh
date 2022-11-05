#set -e
set -x
ZETACORED=/home/ubuntu/go/bin/zetacored
# AWS EC2: testnet1, testnet2, testnet3, testnet5
NODES="3.137.46.147 3.20.194.40 3.19.64.252"

rm -rf ~/.zetacored/data
rm -rf ~/.zetacored/config
for NODE in $NODES; do
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.zetacored/data
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.zetacored/config
done

$ZETACORED init --chain-id athens-1 zetachain --home ~/.zetacored
ADDR=$($ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacored)
if [ -z "$ADDR" ]; then
  echo "No val key found; generate new val key"
  $ZETACORED keys add val --keyring-backend=test --home ~/.zetacored
fi
$ZETACORED add-genesis-account $($ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacored) 1000000000stake --home ~/.zetacored


for NODE in $NODES; do
  ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacored)
  if [ -z "$ADDR" ]; then
    echo "No val key found; generate new val key"
	  ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys add val --keyring-backend=test --home ~/.zetacored
  fi
	ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys show val -a --keyring-backend=test --home ~/.zetacored)
	$ZETACORED add-genesis-account $ADDR 1000000000stake --keyring-backend=test --home ~/.zetacored
done

 
for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacored/config/genesis.json $NODE:~/.zetacored/config/
done


$ZETACORED gentx val 1000000000stake --keyring-backend=test --chain-id=athens-1 --home ~/.zetacored

for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE $ZETACORED gentx val 1000000000stake --keyring-backend=test --chain-id=athens-1 --ip $NODE --home ~/.zetacored
    scp -i ~/.ssh/meta.pem $NODE:~/.zetacored/config/gentx/*.json ~/.zetacored/config/gentx/
done


$ZETACORED collect-gentxs --home ~/.zetacored


for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacored/config/genesis.json $NODE:~/.zetacored/config/
done

#

jq '.chain_id = "athens-1"' ~/.zetacored/config/genesis.json > temp.json && mv temp.json ~/.zetacored/config/genesis.json
sed -i '/\[api\]/,+3 s/enable = false/enable = true/' ~/.zetacored/config/app.toml
sed -i '/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' ~/.zetacored/config/app.toml

for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE jq \'.chain_id = \"athens-1\"\' ~/.zetacored/config/genesis.json > temp.json && mv temp.json ~/.zetacored/config/genesis.json
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+3 s/enable = false/enable = true/\' ~/.zetacored/config/app.toml
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/\' ~/.zetacored/config/app.toml
done

