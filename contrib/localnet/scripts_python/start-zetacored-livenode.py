#!/usr/bin/env python3
"""Start script for ZetaCore node connecting to testnet or mainnet."""

import json
import os
import subprocess
import sys

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


def update_configs():
    ip = subprocess.run("curl -4 -s icanhazip.com", shell=True, capture_output=True, text=True)
    external_ip = ip.stdout.strip() if ip.returncode == 0 else "127.0.0.1"

    config_toml = f"{ZETACORED_CONFIG}/config.toml"
    app_toml = f"{ZETACORED_CONFIG}/app.toml"

    for pattern, repl in [
        ("{YOUR_EXTERNAL_IP_ADDRESS_HERE}", external_ip),
        ("{MONIKER}", MONIKER),
        ('laddr = "tcp:\\/\\/127.0.0.1:26657"', 'laddr = "tcp:\\/\\/0.0.0.0:26657"'),
        ("prometheus = false", "prometheus = true"),
    ]:
        run(f"sed -i 's/{pattern}/{repl}/' \"{config_toml}\"")

    for pattern, repl in [
        ("enable = false", "enable = true"),
        ('address = "tcp:\\/\\/localhost:1317"', 'address = "tcp:\\/\\/0.0.0.0:1317"'),
        ('address = "localhost:9090"', 'address = "0.0.0.0:9090"'),
    ]:
        run(f"sed -i 's/{pattern}/{repl}/' \"{app_toml}\"")


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
