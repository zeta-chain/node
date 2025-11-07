make clean
make install
zetacored init test --chain-id=athens_7001-1

wget https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/genesis.json -O ~/.zetacored/config/genesis.json &&\
wget https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/client.toml -O ~/.zetacored/config/client.toml &&\
wget https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/config.toml -O ~/.zetacored/config/config.toml &&\
wget https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/app.toml -O ~/.zetacored/config/app.toml

# Get external IP address
EXTERNAL_IP=$(curl -4 icanhazip.com)

# Update config.toml with external IP and moniker
sed -i.bak "s/{YOUR_EXTERNAL_IP_ADDRESS_HERE}/$EXTERNAL_IP/g" ~/.zetacored/config/config.toml
sed -i.bak "s/{MONIKER}/testNode/g" ~/.zetacored/config/config.toml


SNAPSHOT_JSON=$(curl -s https://snapshots.rpc.zetachain.com/testnet/fullnode/latest.json)
SNAPSHOT_LINK=$(echo $SNAPSHOT_JSON | jq -r '.snapshots[].link')
SNAPSHOT_FILENAME=$(echo $SNAPSHOT_JSON | jq -r '.snapshots[].filename')


curl "$SNAPSHOT_LINK" -o "$HOME/$SNAPSHOT_FILENAME" && \
lz4 -dc "$HOME/$SNAPSHOT_FILENAME" | tar -C "$HOME/.zetacored/" -xvf - && \
rm "$HOME/$SNAPSHOT_FILENAME"


sleep 10
echo "Starting zetacored..."
export NODE_VERSION=36.0.1
NODE_VERSION=36.0.1 zetacored start