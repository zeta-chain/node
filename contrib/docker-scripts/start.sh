#!/bin/bash

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}


function load_defaults {
  #DEFAULT: Mainnet Statesync.
  export DAEMON_HOME=${DAEMON_HOME:=/root/.zetacored}
  export NETWORK=${NETWORK:=mainnet}
  export RESTORE_TYPE=${RESTORE_TYPE:=statesync}
  export SNAPSHOT_API=${SNAPSHOT_API:=https://snapshots.rpc.zetachain.com}
  export TRUST_HEIGHT_DIFFERENCE_STATE_SYNC=${TRUST_HEIGHT_DIFFERENCE_STATE_SYNC:=40000}
  export COSMOVISOR_VERSION=${COSMOVISOR_VERSION:=v1.5.0}
  export CHAIN_ID=${CHAIN_ID:=zetachain_7000-1}
  export COSMOVISOR_CHECKSUM=${COSMOVISOR_CHECKSUM:=626dfc58c266b85f84a7ed8e2fe0e2346c15be98cfb9f9b88576ba899ed78cdc}
  export VISOR_NAME=${VISOR_NAME:=cosmovisor}
  export DAEMON_NAME=${DAEMON_NAME:=zetacored}
  export DAEMON_ALLOW_DOWNLOAD_BINARIES=${DAEMON_ALLOW_DOWNLOAD_BINARIES:=false}
  export DAEMON_RESTART_AFTER_UPGRADE=${DAEMON_RESTART_AFTER_UPGRADE:=true}
  export UNSAFE_SKIP_BACKUP=${UNSAFE_SKIP_BACKUP:=true}
  export CLIENT_DAEMON_NAME=${CLIENT_DAEMON_NAME:=zetaclientd}
  export CLIENT_DAEMON_ARGS=${CLIENT_DAEMON_ARGS:""}
  export CLIENT_SKIP_UPGRADE=${CLIENT_SKIP_UPGRADE:=true}
  export CLIENT_START_PROCESS=${CLIENT_START_PROCESS:=false}
  export MONIKER=${MONIKER:=local-test}
  export RE_DO_START_SEQUENCE=${RE_DO_START_SEQUENCE:=false}

  #ATHENS3
  export BINARY_LIST_ATHENS3=${BINARY_LIST_ATHENS3:=https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/binary_list.json}
  export STATE_SYNC_RPC_NODE_FILE_ATHENS3=${STATE_SYNC_RPC_NODE_FILE_ATHENS3:=https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/state_sync_node}
  export RPC_STATE_SYNC_RPC_LIST_FILE_ATHENS3=${RPC_STATE_SYNC_RPC_LIST_FILE_ATHENS3:=https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/rpc_state_sync_nodes}
  export APP_TOML_FILE_ATHENS3=${APP_TOML_FILE_ATHENS3:=https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/app.toml}
  export CONFIG_TOML_FILE_ATHENS3=${CONFIG_TOML_FILE_ATHENS3:=https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/config.toml}
  export CLIENT_TOML_FILE_ATHENS3=${CLIENT_TOML_FILE_ATHENS3:=https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/client.toml}
  export GENESIS_FILE_ATHENS3=${GENESIS_FILE_ATHENS3:=https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3/genesis.json}

  #MAINNET
  export BINARY_LIST_MAINNET=${BINARY_LIST_MAINNET:=https://raw.githubusercontent.com/zeta-chain/network-config/main/mainnet/binary_list.json}
  export STATE_SYNC_RPC_NODE_FILE_MAINNET=${STATE_SYNC_RPC_NODE_FILE_MAINNET:=https://raw.githubusercontent.com/zeta-chain/network-config/main/mainnet/state_sync_node}
  export RPC_STATE_SYNC_RPC_LIST_FILE_MAINNET=${RPC_STATE_SYNC_RPC_LIST_FILE_MAINNET:=https://raw.githubusercontent.com/zeta-chain/network-config/main/mainnet/rpc_state_sync_nodes}
  export APP_TOML_FILE_MAINNET=${APP_TOML_FILE_MAINNET:=https://raw.githubusercontent.com/zeta-chain/network-config/main/mainnet/app.toml}
  export CONFIG_TOML_FILE_MAINNET=${CONFIG_TOML_FILE_MAINNET:=https://raw.githubusercontent.com/zeta-chain/network-config/main/mainnet/config.toml}
  export CLIENT_TOML_FILE_MAINNET=${CLIENT_TOML_FILE_MAINNET:=https://raw.githubusercontent.com/zeta-chain/network-config/main/mainnet/client.toml}
  export GENESIS_FILE_MAINNET=${GENESIS_FILE_MAINNET:=https://raw.githubusercontent.com/zeta-chain/network-config/main/mainnet/genesis.json}

}

function init_chain {
  if [ -d "${DAEMON_HOME}/config" ]; then
      logt "${DAEMON_NAME} home directory already initialized."
  else
      logt "${DAEMON_NAME} home directory not initialized."
      logt "MONIKER: ${MONIKER}"
      logt "DAEMON_HOME: ${DAEMON_HOME}"
      ${DAEMON_NAME} init ${MONIKER} --home ${DAEMON_HOME} --chain-id ${CHAIN_ID}
  fi
}

function download_configs {
  if [ "${NETWORK}" == "mainnet" ]; then
    wget -q ${APP_TOML_FILE_MAINNET} -O ${DAEMON_HOME}/config/app.toml
    wget -q ${CONFIG_TOML_FILE_MAINNET} -O ${DAEMON_HOME}/config/config.toml
    wget -q ${CLIENT_TOML_FILE_MAINNET} -O ${DAEMON_HOME}/config/client.toml
    wget -q ${GENESIS_FILE_MAINNET} -O ${DAEMON_HOME}/config/genesis.json
    wget -q ${BINARY_LIST_MAINNET}
    export DOWNLOAD_BINARIES=$(cat binary_list.json | tr -d '\n')
    rm -rf binary_list.json
    logt "BINARY_LIST: ${DOWNLOAD_BINARIES}"
  elif [ "${NETWORK}" == "athens3" ]; then
    wget -q ${APP_TOML_FILE_ATHENS3} -O ${DAEMON_HOME}/config/app.toml
    wget -q ${CONFIG_TOML_FILE_ATHENS3} -O ${DAEMON_HOME}/config/config.toml
    wget -q ${CLIENT_TOML_FILE_ATHENS3} -O ${DAEMON_HOME}/config/client.toml
    wget -q ${GENESIS_FILE_ATHENS3} -O ${DAEMON_HOME}/config/genesis.json
    wget -q ${BINARY_LIST_ATHENS3}
    export DOWNLOAD_BINARIES=$(cat binary_list.json | tr -d '\n')
    rm -rf binary_list.json
    logt "BINARY_LIST: ${DOWNLOAD_BINARIES}"
  else
    logt "Initialize for Localnet."
  fi
}

function setup_restore_type {
  if [ "${RESTORE_TYPE}" == "statesync" ]; then
    logt "Statesync restore. Download state sync rpc address from network-config"
    if [ "${NETWORK}" == "mainnet" ]; then
      logt "MAINNET STATE SYNC"
      logt "STATE_SYNC_RPC_NODE_FILE_MAINNET: ${STATE_SYNC_RPC_NODE_FILE_MAINNET}"
      logt "RPC_STATE_SYNC_RPC_LIST_FILE_MAINNET: ${RPC_STATE_SYNC_RPC_LIST_FILE_MAINNET}"
      wget -q ${STATE_SYNC_RPC_NODE_FILE_MAINNET}
      wget -q ${RPC_STATE_SYNC_RPC_LIST_FILE_MAINNET}
      export STATE_SYNC_SERVER=$(cat state_sync_node)
      export RPC_STATE_SYNC_SERVERS=$(cat rpc_state_sync_nodes)
      rm -rf state_sync_node
      rm -rf rpc_state_sync_nodes
    elif [ "${NETWORK}" == "athens3" ]; then
      logt "ATHENS STATE SYNC"
      logt "STATE_SYNC_RPC_NODE_FILE_MAINNET: ${STATE_SYNC_RPC_NODE_FILE_ATHENS3}"
      logt "RPC_STATE_SYNC_RPC_LIST_FILE_MAINNET: ${RPC_STATE_SYNC_RPC_LIST_FILE_ATHENS3}"
      wget -q ${STATE_SYNC_RPC_NODE_FILE_ATHENS3}
      wget -q ${RPC_STATE_SYNC_RPC_LIST_FILE_ATHENS3}
      export STATE_SYNC_SERVER=$(cat state_sync_node)
      export RPC_STATE_SYNC_SERVERS=$(cat rpc_state_sync_nodes)
      rm -rf state_sync_node
      rm -rf rpc_state_sync_nodes
    fi
  elif [ "${RESTORE_TYPE}" == "snapshot"  ]; then
    if [ "${NETWORK}" == "mainnet" ]; then
      logt "Get Latest Snapshot URL"
      SNAPSHOT_URL=$(curl -s ${SNAPSHOT_API}/mainnet/fullnode/latest.json | jq -r '.snapshots[0].link')
      SNAPSHOT_FILENAME=$(curl -s ${SNAPSHOT_API}/mainnet/fullnode/latest.json | jq -r '.snapshots[0].filename')
      SNAPSHOT_DIR=$(pwd)
      logt "Download Snapshot from url: ${SNAPSHOT_URL}"
      curl -o "${SNAPSHOT_FILENAME}" "${SNAPSHOT_URL}"
      logt "Change to: ${DAEMON_HOME} and extract snapshot."
      cd ${DAEMON_HOME}
      tar xvf ${SNAPSHOT_DIR}/${SNAPSHOT_FILENAME}
      logt " Cleanup Snapshot"
      rm -rf ${SNAPSHOT_DIR}/${SNAPSHOT_FILENAME}
    elif [ "${NETWORK}" == "athens3" ]; then
      SNAPSHOT_URL=$(curl -s ${SNAPSHOT_API}/testnet/fullnode/latest.json | jq -r '.snapshots[0].link')
      SNAPSHOT_FILENAME=$(curl -s ${SNAPSHOT_API}/testnet/fullnode/latest.json | jq -r '.snapshots[0].filename')
      SNAPSHOT_DIR=$(pwd)
      logt "Download Snapshot from url: ${SNAPSHOT_URL}"
      curl -o "${SNAPSHOT_FILENAME}" "${SNAPSHOT_URL}"
      logt "Change to: ${DAEMON_HOME} and extract snapshot."
      cd ${DAEMON_HOME}
      tar xvf ${SNAPSHOT_DIR}/${SNAPSHOT_FILENAME}
      logt " Cleanup Snapshot"
      rm -rf ${SNAPSHOT_DIR}/${SNAPSHOT_FILENAME}
    fi
  elif [ "${RESTORE_TYPE}" == "snapshot-archive"  ]; then
    if [ "${NETWORK}" == "mainnet" ]; then
      logt "Get Latest Snapshot URL"
      SNAPSHOT_URL=$(curl -s ${SNAPSHOT_API}/mainnet/archive/latest.json | jq -r '.snapshots[0].link')
      SNAPSHOT_FILENAME=$(curl -s ${SNAPSHOT_API}/mainnet/archive/latest.json | jq -r '.snapshots[0].filename')
      SNAPSHOT_DIR=$(pwd)
      logt "Download Snapshot from url: ${SNAPSHOT_URL}"
      curl -o "${SNAPSHOT_FILENAME}" "${SNAPSHOT_URL}"
      logt "Change to: ${DAEMON_HOME} and extract snapshot."
      cd ${DAEMON_HOME}
      tar xvf ${SNAPSHOT_DIR}/${SNAPSHOT_FILENAME}
      logt " Cleanup Snapshot"
      rm -rf ${SNAPSHOT_DIR}/${SNAPSHOT_FILENAME}
    elif [ "${NETWORK}" == "athens3" ]; then
      SNAPSHOT_URL=$(curl -s ${SNAPSHOT_API}/testnet/archive/latest.json | jq -r '.snapshots[0].link')
      SNAPSHOT_FILENAME=$(curl -s ${SNAPSHOT_API}/testnet/archive/latest.json | jq -r '.snapshots[0].filename')
      SNAPSHOT_DIR=$(pwd)
      logt "Download Snapshot from url: ${SNAPSHOT_URL}"
      curl -o "${SNAPSHOT_FILENAME}" "${SNAPSHOT_URL}"
      logt "Change to: ${DAEMON_HOME} and extract snapshot."
      cd ${DAEMON_HOME}
      tar xvf ${SNAPSHOT_DIR}/${SNAPSHOT_FILENAME}
      logt " Cleanup Snapshot"
      rm -rf ${SNAPSHOT_DIR}/${SNAPSHOT_FILENAME}
    fi
  else
    logt "Initialize for Localnet."
  fi
}

function change_config_values {
  if [ "${RESTORE_TYPE}" == "statesync" ]; then
    export STATE_SYNC_SERVER="${STATE_SYNC_SERVER}"
    export TRUST_HEIGHT=$(curl -s ${STATE_SYNC_SERVER}/block | jq -r '.result.block.header.height')
    export HEIGHT=$((TRUST_HEIGHT-${TRUST_HEIGHT_DIFFERENCE_STATE_SYNC}))
    export TRUST_HASH=$(curl -s "${STATE_SYNC_SERVER}/block?height=${HEIGHT}" | jq -r '.result.block_id.hash')
    export RPC_STATE_SYNC_SERVERS="${RPC_STATE_SYNC_SERVERS}"
    export EXTERNAL_IP=$(curl -4 icanhazip.com)

    logt "******* DEBUG STATE SYNC VALUES *******"
    logt "STATE_SYNC_SERVER: ${STATE_SYNC_SERVER}"
    logt "RPC_STATE_SYNC_SERVERS: ${RPC_STATE_SYNC_SERVERS}"
    logt "TRUST_HEIGHT: ${TRUST_HEIGHT}"
    logt "TRUST_HASH: ${TRUST_HASH}"
    logt "HEIGHT: ${HEIGHT}"
    logt "EXTERNAL_IP: ${EXTERNAL_IP}"

    logt "SED Change Config Files."
    sed -i -e "s/^enable = .*/enable = \"true\"/" ${DAEMON_HOME}/config/config.toml
    sed -i -e "s/^rpc_servers = .*/rpc_servers = \"${RPC_STATE_SYNC_SERVERS}\"/" ${DAEMON_HOME}/config/config.toml
    sed -i -e "s/^trust_height = .*/trust_height = \"${HEIGHT}\"/" ${DAEMON_HOME}/config/config.toml
    sed -i -e "s/^trust_hash = .*/trust_hash = \"${TRUST_HASH}\"/" ${DAEMON_HOME}/config/config.toml
    sed -i -e "s/^moniker = .*/moniker = \"${MONIKER}\"/" ${DAEMON_HOME}/config/config.toml
    sed -i -e "s/^external_address = .*/external_address = \"${EXTERNAL_IP}:26656\"/" ${DAEMON_HOME}/config/config.toml
  else
    export EXTERNAL_IP=$(curl -4 icanhazip.com)
    logt "******* DEBUG STATE SYNC VALUES *******"
    logt "EXTERNAL_IP: ${EXTERNAL_IP}"
    logt "SED Change Config Files."
    sed -i -e "s/^enable = .*/enable = \"true\"/" ${DAEMON_HOME}/config/config.toml
    sed '/^\[statesync\]/,/^\[/ s/enable = "true"/enable = "false"/' ${DAEMON_HOME}/config/config.toml
    sed -i -e "s/^moniker = .*/moniker = \"${MONIKER}\"/" ${DAEMON_HOME}/config/config.toml
    sed -i -e "s/^external_address = .*/external_address = \"${EXTERNAL_IP}:26656\"/" ${DAEMON_HOME}/config/config.toml
  fi
}

function setup_basic_keyring {
  if ${DAEMON_NAME} keys show "$MONIKER" --keyring-backend test > /dev/null 2>&1; then
    echo "Key $MONIKER already exists."
  else
    ${DAEMON_NAME} keys add "$MONIKER" --keyring-backend test
    echo "Key $MONIKER created."
  fi
}

function download_binary_version {
  if [ "${NETWORK}" == "mainnet" ]; then
    wget -q ${BINARY_LIST_MAINNET}
    export DOWNLOAD_BINARIES=$(cat binary_list.json | tr -d '\n')
    rm -rf binary_list.json
    logt "BINARY_LIST: ${DOWNLOAD_BINARIES}"
  elif [ "${NETWORK}" == "athens3" ]; then
    wget -q ${BINARY_LIST_ATHENS3}
    export DOWNLOAD_BINARIES=$(cat binary_list.json | tr -d '\n')
    rm -rf binary_list.json
    logt "BINARY_LIST: ${DOWNLOAD_BINARIES}"
  fi
  python3 /scripts/download_binaries.py
}

function move_zetacored_binaries {
  mkdir -p ${DAEMON_HOME}/cosmovisor || logt "Directory already exists ${DAEMON_HOME}/cosmovisor"
  mkdir -p ${DAEMON_HOME}/cosmovisor/genesis || logt "Directory already exists ${DAEMON_HOME}/cosmovisor/genesis"
  mkdir -p ${DAEMON_HOME}/cosmovisor/genesis/bin || logt "Directory already exists ${DAEMON_HOME}/cosmovisor/genesis/bin"
  cp /usr/local/bin/zetacored ${DAEMON_HOME}/cosmovisor/genesis/bin/zetacored

  if [ "${RESTORE_TYPE}" == "statesync" ]; then
      logt "Its statesync so cosmosvisor won't know which binary to start from so make sure it starts from the latest version reported in ABCI_INFO from statesync server rpc."
      export VERSION_CHECK=$(curl -s ${STATE_SYNC_SERVER}/abci_info | jq -r '.result.response.version')
      logt "CURRENT VERSION_CHECK: ${VERSION_CHECK}"
      cp ${DAEMON_HOME}/cosmovisor/upgrades/v${VERSION_CHECK}/bin/zetacored ${DAEMON_HOME}/cosmovisor/genesis/bin/zetacored
  fi
}

function start_network {
  if [ "${IS_LOCAL_DEVELOPMENT}" == "true" ]; then
    cp /usr/local/bin/zetacored ${DAEMON_HOME}/cosmovisor/genesis/bin/zetacored
    find /root/.zetacored/cosmovisor/upgrades/ -type f -path "*/bin/zetacored" -exec cp /usr/local/bin/zetacored {} \;
  fi
  expected_major_version=$(cat /scripts/expected_major_version | cut -d '-' -f 1)
  VISOR_VERSION=v$(${VISOR_NAME} version | tail -n 1 | tr -d '(devel)' | tr -d '\n')
  DAEMON_VERSION=$(${DAEMON_NAME} version)
  VISOR_MAJOR_VERSION=$(echo $VISOR_VERSION | grep -o '^v[0-9]*')
  DAEMON_MAJOR_VERSION=$(echo $DAEMON_VERSION | grep -o '^v[0-9]*')

  logt "EXPECTED_VERSION_WITH_HASH: $(cat /scripts/expected_major_version | cut -d '-' -f 1)"
  logt "EXPECTED_MAJOR_VERSION: ${expected_major_version}"
  logt "VISOR_VERSION: ${VISOR_VERSION}"
  logt "DAEMON_VERSION: ${DAEMON_VERSION}"
  logt "VISOR_MAJOR_VERSION: ${VISOR_MAJOR_VERSION}"
  logt "DAEMON_MAJOR_VERSION: ${DAEMON_MAJOR_VERSION}"

  if [ "$VISOR_MAJOR_VERSION" != "$expected_major_version" ] || [ "$DAEMON_MAJOR_VERSION" != "$expected_major_version" ]; then
      logt "One or both versions don't match the expected major release version: $expected_major_version"
  else
      logt "Both versions match the expected major release version: $expected_major_version"
  fi

  if [ "$VISOR_VERSION" != "$DAEMON_VERSION" ]; then
      logt "cosmovisor version doesn't appear to match your daemon version. Start ${DAEMON_NAME}"
  else
    logt "cosmovisor version match your daemon version. Start ${VISOR_NAME}"
  fi

  ${VISOR_NAME} run start --home ${DAEMON_HOME} \
    --log_level info \
    --moniker ${MONIKER} \
    --rpc.laddr tcp://0.0.0.0:26657 \
    --minimum-gas-prices 1.0azeta "--grpc.enable=true"
}

logt "Load Default Values for ENV Vars if not set."
load_defaults

if [[ -f "${DAEMON_HOME}/start_sequence_status" ]] && grep -q "START_SEQUENCE_COMPLETE" "${DAEMON_HOME}/start_sequence_status" && [[ "$RE_DO_START_SEQUENCE" != "true" ]]; then
    logt "The start sequence is complete and no redo is required."

    logt "Download Configs"
    download_configs

    logt "Download Historical Binaries"
    download_binary_version

    if [ "${RESTORE_TYPE}" == "statesync" ]; then
      logt "Setup Restore Type: ${RESTORE_TYPE}"
      logt "During restarts, we re-do this to ensure to update the configs with valid values. When you call change config the stuff that gets set in this function for statesync needs to be set. Doesn't effect to re-set this."
      setup_restore_type
    fi

    logt "Modify Chain Configs"
    change_config_values

    logt "Move Zetacored Binaries."
    move_zetacored_binaries

    logt "Start sequence has completed, echo into file so on restart it doesn't download snapshots again."
    echo "START_SEQUENCE_COMPLETE" >> ${DAEMON_HOME}/start_sequence_status

    logt "Start Network"
    start_network
else
  logt "START_SEQUENCE_COMPLETE is not true, or RE_DO_START_SEQUENCE is set to true."

  if [[ "$RE_DO_START_SEQUENCE" == "true" ]]; then
    logt "Clean any files that may exist in: ${DAEMON_HOME}"
    rm -rf ${DAEMON_HOME}/* || logt "directory doesn't exist."
  fi

  logt "Init Chain"
  init_chain

  logt "Download Configs"
  download_configs

  logt "Download Historical Binaries"
  download_binary_version

  logt "Setup Restore Type: ${RESTORE_TYPE}"
  setup_restore_type

  logt "Modify Chain Configs"
  change_config_values

  logt "Move root binaries to current"
  move_zetacored_binaries

  logt "Start sequence has completed, echo into file so on restart it doesn't download snapshots again."
  echo "START_SEQUENCE_COMPLETE" >> ${DAEMON_HOME}/start_sequence_status

  logt "Start Network"
  start_network
fi
