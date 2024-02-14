#!/bin/bash

function logt(){
 echo "$(date) - $1"
}

function download_binary() {
  wget https://github.com/zeta-chain/node/releases/download/${BINARY_VERSION}/zetacored-darwin-amd64 -O /usr/local/bin/zetacored
  chmod a+x /usr/local/bin/zetacored
  zetacored version || echo "BINARY NOT INSTALLED" && exit 1
}

function chain_init() {
  ZETACORED_DIR="$HOME/.zetacored"
  # Check if the .zetacored directory exists
  if [ -d "$ZETACORED_DIR" ]; then
      echo ".zetacored directory already exists at $ZETACORED_DIR."
  else
      # Directory does not exist, initialize zetacored
      zetacored init "$MONIKER" --chain-id "$CHAIN_ID"
      echo ".zetacored initialized for $MONIKER with chain ID $CHAIN_ID."
  fi
}

function modify_chain_configs() {
  sed -i -e "s/^enable = .*/enable = \"true\"/" /root/.zetacored/config/config.toml
  sed -i -e "s/^rpc_servers = .*/rpc_servers = \"${RPC_STATE_SYNC_SERVERS}\"/" /root/.zetacored/config/config.toml
  sed -i -e "s/^trust_height = .*/trust_height = \"${HEIGHT}\"/" /root/.zetacored/config/config.toml
  sed -i -e "s/^trust_hash = .*/trust_hash = \"${TRUST_HASH}\"/" /root/.zetacored/config/config.toml
  sed -i -e "s/^moniker = .*/moniker = \"${MONIKER}\"/" /root/.zetacored/config/config.toml
  sed -i -e "s/^external_address = .*/external_address = \"${EXTERNAL_IP}:26656\"/" /root/.zetacored/config/config.toml
  sed -i -e "s/^seeds = .*/seeds = \"${SEED}\"/" /root/.zetacored/config/config.toml
  sed -i -e 's/^max_num_inbound_peers = .*/max_num_inbound_peers = 120/' /root/.zetacored/config/config.toml
  sed -i -e 's/^max_num_outbound_peers = .*/max_num_outbound_peers = 60/' /root/.zetacored/config/config.toml
  sed -i -e "s/^persistent_peers = .*/persistent_peers = \"${PERSISTENT_PEERS}\"/" /root/.zetacored/config/config.toml
}

function setup_basic_keyring() {
  if zetacored keys show "$MONIKER" --keyring-backend test > /dev/null 2>&1; then
    echo "Key $MONIKER already exists."
  else
    zetacored keys add "$MONIKER" --keyring-backend test
    echo "Key $MONIKER created."
  fi
}

function start_network() {
  zetacored start --home /root/.zetacored/ \
    --log_level info \
    --moniker ${MONIKER} \
    --rpc.laddr tcp://0.0.0.0:26657 \
    --minimum-gas-prices 1.0azeta "--grpc.enable=true"
}

function install_dependencies() {
  apt-get update
  apt-get install nano jq -y
}

function check_configs_debug() {
  logt "Check home config directory ensure configs present."
  ls -lah /root/.zetacored/config

  logt "Check the zetacored binary is in /usr/local/bin"
  ls -lah /usr/local/bin/

  logt "Check zetacored root directory"
  ls -lah /root/.zetacored

  logt "Config.toml"
  cat /root/.zetacored/config/config.toml
  logt "******"

  logt "Config.toml"
  cat /root/.zetacored/config/app.toml
  logt "******"

  logt "Config.toml"
  cat /root/.zetacored/config/client.toml
  logt "******"

  logt "Config.toml"
  cat /root/.zetacored/config/genesis.json
  logt "******"
}

logt "Install Dependencies"
install_dependencies

if [ "${DEBUG}" == "true" ]; then
  check_configs_debug
fi

logt "Setup script variables."
export STATE_SYNC_SERVER="${STATE_SYNC_SERVER}"
export TRUST_HEIGHT=$(curl -s http://${STATE_SYNC_SERVER}/block | jq -r '.result.block.header.height')
export HEIGHT=$((TRUST_HEIGHT-40000))
#export HEIGHT=$((TRUST_HEIGHT-100))
export TRUST_HASH=$(curl -s "http://${STATE_SYNC_SERVER}/block?height=${HEIGHT}" | jq -r '.result.block_id.hash')
export RPC_STATE_SYNC_SERVERS="${RPC_STATE_SYNC_SERVERS}"
export SEED="${SEED_NODE}"
export PERSISTENT_PEERS="${PEERS}"
export EXTERNAL_IP=$(curl -4 icanhazip.com)

if [ "$DOWNLOAD_BINARY" = true ]; then
  logt "Download chain binary"
  download_binary
else
  logt "User built binary."
fi

logt "Init the chain directory"
chain_init

logt "Modify chain configs."
modify_chain_configs

if [ "${DEBUG}" == "true" ]; then
  check_configs_debug
fi

logt "Start network"
start_network
