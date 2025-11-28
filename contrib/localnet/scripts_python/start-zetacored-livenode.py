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
    if subprocess.run(f'wget -q "{url}" -O "{dest}"', shell=True).returncode != 0:
        sys.exit(1)


def initialize_node(config):
    genesis = f"{ZETACORED_CONFIG}/genesis.json"
    if not os.path.isfile(genesis):
        run(f'zetacored init "{MONIKER}" --chain-id={CHAIN_ID}')
        for f in ["genesis.json", "client.toml", "config.toml", "app.toml"]:
            download(f"{config['config_base']}/{f}", f"{ZETACORED_CONFIG}/{f}")


def setup_snapshot(config):
    snapshot_data = f"{config['snapshot_cache']}/data"
    data_dir = f"{ZETACORED_HOME}/data"

    if FORCE_DOWNLOAD:
        run(f'rm -rf "{config["snapshot_cache"]}"/*', check=False)

    # Check if snapshot data exists (application.db indicates real snapshot, not just init)
    has_snapshot = os.path.isdir(f"{data_dir}/application.db")

    if os.path.isdir(snapshot_data) and os.listdir(snapshot_data):
        if not has_snapshot:
            if os.path.isdir(data_dir):
                run(f'rm -rf "{data_dir}"')
            run(f'cp -r "{snapshot_data}" "{ZETACORED_HOME}/"')
    else:
        run(f"python3 -u /root/download_snapshot.py --chain-id {CHAIN_ID}")
        if os.path.isdir(snapshot_data) and os.listdir(snapshot_data):
            if os.path.isdir(data_dir):
                run(f'rm -rf "{data_dir}"')
            run(f'cp -r "{snapshot_data}" "{ZETACORED_HOME}/"')

    # Create priv_validator_state.json if missing
    pvs = f"{data_dir}/priv_validator_state.json"
    if not os.path.isfile(pvs):
        with open(pvs, 'w') as f:
            json.dump({"height": "0", "round": 0, "step": 0}, f)


def get_external_ip():
    result = subprocess.run(["curl", "-4", "-s", "icanhazip.com"], capture_output=True, text=True)
    return result.stdout.strip() if result.returncode == 0 else "127.0.0.1"


def update_configs():
    p2p_port = 26656
    rpc_port = 26657
    api_port = 1317
    grpc_port = 9090

    external_ip = get_external_ip()

    # Update config.toml
    config_path = f"{ZETACORED_CONFIG}/config.toml"
    config = toml.load(config_path)
    config["moniker"] = MONIKER
    config["p2p"]["external_address"] = f"{external_ip}:{p2p_port}"
    config["rpc"]["laddr"] = f"tcp://0.0.0.0:{rpc_port}"
    config["instrumentation"]["prometheus"] = True
    with open(config_path, "w") as f:
        toml.dump(config, f)

    # Update app.toml
    app_path = f"{ZETACORED_CONFIG}/app.toml"
    app = toml.load(app_path)
    app["api"]["enable"] = True
    app["api"]["address"] = f"tcp://0.0.0.0:{api_port}"
    app["grpc"]["address"] = f"0.0.0.0:{grpc_port}"
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
