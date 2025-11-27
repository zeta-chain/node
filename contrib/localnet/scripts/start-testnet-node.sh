#!/bin/bash
set -e

MONIKER=${MONIKER:-"testNode"}
FORCE_DOWNLOAD=${FORCE_DOWNLOAD:-false}
SNAPSHOT_CACHE="/root/zetacored_snapshot_testnet"
ZETACORED_HOME="/root/.zetacored"
ZETACORED_CONFIG="${ZETACORED_HOME}/config"

ATHENS3_CONFIG_BASE="https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3"
ATHENS3_GENESIS_URL="${ATHENS3_CONFIG_BASE}/genesis.json"
ATHENS3_CLIENT_URL="${ATHENS3_CONFIG_BASE}/client.toml"
ATHENS3_CONFIG_URL="${ATHENS3_CONFIG_BASE}/config.toml"
ATHENS3_APP_URL="${ATHENS3_CONFIG_BASE}/app.toml"

download_file() {
    local url=$1
    local dest=$2
    wget -q "${url}" -O "${dest}" || {
        echo "Failed to download ${url}"
        exit 1
    }
}

if [ ! -f "${ZETACORED_CONFIG}/genesis.json" ]; then
    zetacored init "${MONIKER}" --chain-id=athens_7001-1
    download_file "${ATHENS3_GENESIS_URL}" "${ZETACORED_CONFIG}/genesis.json"
    download_file "${ATHENS3_CLIENT_URL}" "${ZETACORED_CONFIG}/client.toml"
    download_file "${ATHENS3_CONFIG_URL}" "${ZETACORED_CONFIG}/config.toml"
    download_file "${ATHENS3_APP_URL}" "${ZETACORED_CONFIG}/app.toml"
fi

SNAPSHOT_DATA_DIR="${SNAPSHOT_CACHE}/data"

if [ "${FORCE_DOWNLOAD}" = "true" ]; then
    rm -rf "${SNAPSHOT_DATA_DIR}"
fi

if [ -d "${SNAPSHOT_DATA_DIR}" ] && [ "$(ls -A ${SNAPSHOT_DATA_DIR})" ]; then
    if [ ! -d "${ZETACORED_HOME}/data" ] || [ -z "$(ls -A ${ZETACORED_HOME}/data)" ]; then
        cp -r "${SNAPSHOT_DATA_DIR}" "${ZETACORED_HOME}/"
    fi
else
    python3 /root/download_snapshot.py
    if [ -d "${SNAPSHOT_DATA_DIR}" ]; then
        if [ -d "${ZETACORED_HOME}/data" ]; then
            rm -rf "${ZETACORED_HOME}/data"
        fi
        cp -r "${SNAPSHOT_DATA_DIR}" "${ZETACORED_HOME}/"
    fi
fi

EXTERNAL_IP=$(curl -4 -s icanhazip.com || echo "127.0.0.1")

#Update config.toml
sed -i "s/{YOUR_EXTERNAL_IP_ADDRESS_HERE}/${EXTERNAL_IP}/g" "${ZETACORED_CONFIG}/config.toml"
sed -i "s/{MONIKER}/${MONIKER}/g" "${ZETACORED_CONFIG}/config.toml"
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' "${ZETACORED_CONFIG}/config.toml"
sed -i 's/enable = false/enable = true/' "${ZETACORED_CONFIG}/app.toml"
sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/' "${ZETACORED_CONFIG}/app.toml"
sed -i 's/address = "localhost:9090"/address = "0.0.0.0:9090"/' "${ZETACORED_CONFIG}/app.toml"
sed -i 's/prometheus = false/prometheus = true/' "${ZETACORED_CONFIG}/config.toml"

exec zetacored start