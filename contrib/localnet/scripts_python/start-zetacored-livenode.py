#!/usr/bin/env python3
"""Start script for ZetaCore node connecting to testnet or mainnet."""

import json
import os
import subprocess
import sys
import toml

MONIKER = os.environ.get('MONIKER', 'testNode')
FORCE_DOWNLOAD = os.environ.get('FORCE_DOWNLOAD', 'false').lower() == 'true'
CHAIN_ID = os.environ.get('CHAIN_ID', 'athens_7001-1')

ZETACORED_HOME = "/root/.zetacored"
ZETACORED_CONFIG = f"{ZETACORED_HOME}/config"

NETWORK_CONFIGS = {
    "athens_7001-1": {
        "name": "Testnet",
        "config_base": "https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3",
        "snapshot_cache": "/root/zetacored_snapshot_testnet",
    },
    "zetachain_7000-1": {
        "name": "Mainnet",
        "config_base": "https://raw.githubusercontent.com/zeta-chain/network-config/main/mainnet",
        "snapshot_cache": "/root/zetacored_snapshot_mainnet",
    },
}


def run(cmd, check=True):
    return subprocess.run(cmd, shell=True, check=check).returncode == 0


def download(url, dest):
    subprocess.run(f'wget -q "{url}" -O "{dest}"', shell=True, check=True)


def initialize_node(config):
    genesis = f"{ZETACORED_CONFIG}/genesis.json"
    if not os.path.isfile(genesis):
        run(f'zetacored init "{MONIKER}" --chain-id={CHAIN_ID}')
        for f in ["genesis.json", "client.toml", "config.toml", "app.toml"]:
            download(f"{config['config_base']}/{f}", f"{ZETACORED_CONFIG}/{f}")


def copy_snapshot_data(snapshot_data, data_dir):
    """Copy snapshot data to zetacored home if not already present."""
    has_snapshot = os.path.isdir(f"{data_dir}/application.db")
    if not has_snapshot and os.path.isdir(snapshot_data) and os.listdir(snapshot_data):
        if os.path.isdir(data_dir):
            run(f'rm -rf "{data_dir}"')
        run(f'cp -r "{snapshot_data}" "{ZETACORED_HOME}/"')


def setup_snapshot(config):
    snapshot_data = f"{config['snapshot_cache']}/data"
    data_dir = f"{ZETACORED_HOME}/data"

    if FORCE_DOWNLOAD:
        run(f'rm -rf "{config["snapshot_cache"]}"/*', check=False)

    # Try cached snapshot first, otherwise download
    if not (os.path.isdir(snapshot_data) and os.listdir(snapshot_data)):
        run(f"python3 -u /root/download_snapshot.py --chain-id {CHAIN_ID}")

    copy_snapshot_data(snapshot_data, data_dir)

    # Create priv_validator_state.json if missing
    pvs = f"{data_dir}/priv_validator_state.json"
    if not os.path.isfile(pvs):
        with open(pvs, 'w') as f:
            json.dump({"height": "0", "round": 0, "step": 0}, f)


def get_external_ip():
    result = subprocess.run(["curl", "-4", "-s", "icanhazip.com"], capture_output=True, text=True)
    return result.stdout.strip() if result.returncode == 0 else "127.0.0.1"


def update_configs():
    external_ip = get_external_ip()

    # Update config.toml
    config_path = f"{ZETACORED_CONFIG}/config.toml"
    config = toml.load(config_path)
    config["moniker"] = MONIKER
    config["p2p"]["external_address"] = f"{external_ip}:26656"
    config["rpc"]["laddr"] = "tcp://0.0.0.0:26657"
    config["instrumentation"]["prometheus"] = True
    with open(config_path, "w") as f:
        toml.dump(config, f)

    # Update app.toml
    app_path = f"{ZETACORED_CONFIG}/app.toml"
    app = toml.load(app_path)
    app["api"]["enable"] = True
    app["api"]["address"] = "tcp://0.0.0.0:1317"
    app["grpc"]["address"] = "0.0.0.0:9090"
    with open(app_path, "w") as f:
        toml.dump(app, f)


def main():
    if CHAIN_ID not in NETWORK_CONFIGS:
        print(f"Unknown chain ID: {CHAIN_ID}")
        sys.exit(1)

    config = NETWORK_CONFIGS[CHAIN_ID]
    print(f"Starting {config['name']} node (chain: {CHAIN_ID})")

    initialize_node(config)
    setup_snapshot(config)
    update_configs()

    os.execv("/usr/local/bin/zetacored", ["/usr/local/bin/zetacored", "start"])


if __name__ == "__main__":
    main()
