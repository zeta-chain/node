set -e
set -x
METACORED=/home/ubuntu/go/bin/metacored
NODES="10.0.0.22"
rm -rf ~/.metacore
for NODE in $NODES; do
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.metacore
done

$METACORED init --chain-id testing zetachain
$METACORED keys add val --keyring-backend=test
$METACORED add-genesis-account $($METACORED keys show val -a --keyring-backend=test) 1000000000stake


for NODE in $NODES; do
	ssh -i ~/.ssh/meta.pem $NODE $METACORED keys add val --keyring-backend=test
	ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $METACORED keys show val -a --keyring-backend=test)
	$METACORED add-genesis-account $ADDR 1000000000stake --keyring-backend=test
done


for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.metacore/config/genesis.json $NODE:~/.metacore/config/
done


$METACORED gentx val 1000000000stake --keyring-backend=test --chain-id=testing
for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE $METACORED gentx val 1000000000stake --keyring-backend=test --chain-id=testing
    scp -i ~/.ssh/meta.pem $NODE:~/.metacore/config/gentx/*.json ~/.metacore/config/gentx/
done

$METACORED collect-gentxs


for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.metacore/config/genesis.json $NODE:~/.metacore/config/
done


jq '.chain_id = "testing"' ~/.metacore/config/genesis.json > temp.json && mv temp.json ~/.metacore/config/genesis.json
sed -i '/\[api\]/,+3 s/enable = false/enable = true/' ~/.metacore/config/app.toml
sed -i '/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' ~/.metacore/config/app.toml
for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE jq \'.chain_id = \"testing\"\' ~/.metacore/config/genesis.json > temp.json && mv temp.json ~/.metacore/config/genesis.json
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+3 s/enable = false/enable = true/\' ~/.metacore/config/app.toml
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/\' ~/.metacore/config/app.toml
done
