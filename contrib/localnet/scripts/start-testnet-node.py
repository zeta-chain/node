#!/usr/bin/env python3
"""
Start script for ZetaCore testnet node.
This script initializes and configures zetacored for connecting to Athens3 testnet.
"""

import os
import subprocess
import sys

# Configuration
MONIKER = os.environ.get('MONIKER', 'testNode')
FORCE_DOWNLOAD = os.environ.get('FORCE_DOWNLOAD', 'false').lower() == 'true'
SNAPSHOT_CACHE = "/root/zetacored_snapshot_testnet"
ZETACORED_HOME = "/root/.zetacored"
ZETACORED_CONFIG = f"{ZETACORED_HOME}/config"

ATHENS3_CONFIG_BASE = "https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3"
ATHENS3_GENESIS_URL = f"{ATHENS3_CONFIG_BASE}/genesis.json"
ATHENS3_CLIENT_URL = f"{ATHENS3_CONFIG_BASE}/client.toml"
ATHENS3_CONFIG_URL = f"{ATHENS3_CONFIG_BASE}/config.toml"
ATHENS3_APP_URL = f"{ATHENS3_CONFIG_BASE}/app.toml"


def run_command(cmd, check=True):
    """Run a shell command."""
    result = subprocess.run(cmd, shell=True, check=check)
    return result.returncode == 0


def download_file(url, dest):
    """Download a file from URL to destination."""
    print(f"Downloading {url}...")
    result = subprocess.run(
        f'wget -q "{url}" -O "{dest}"',
        shell=True
    )
    if result.returncode != 0:
        print(f"Failed to download {url}")
        sys.exit(1)


def dir_exists_and_not_empty(path):
    """Check if directory exists and is not empty."""
    return os.path.isdir(path) and os.listdir(path)


def sed_replace(file_path, pattern, replacement):
    """Replace pattern in file using sed."""
    run_command(f"sed -i 's/{pattern}/{replacement}/' \"{file_path}\"")


def initialize_node():
    """Initialize zetacored and download config files."""
    genesis_path = f"{ZETACORED_CONFIG}/genesis.json"

    if not os.path.isfile(genesis_path):
        print(f"Initializing node with moniker: {MONIKER}")
        run_command(f'zetacored init "{MONIKER}" --chain-id=athens_7001-1')
        download_file(ATHENS3_GENESIS_URL, f"{ZETACORED_CONFIG}/genesis.json")
        download_file(ATHENS3_CLIENT_URL, f"{ZETACORED_CONFIG}/client.toml")
        download_file(ATHENS3_CONFIG_URL, f"{ZETACORED_CONFIG}/config.toml")
        download_file(ATHENS3_APP_URL, f"{ZETACORED_CONFIG}/app.toml")


def setup_snapshot():
    """Setup snapshot data."""
    snapshot_data_dir = f"{SNAPSHOT_CACHE}/data"
    zetacored_data_dir = f"{ZETACORED_HOME}/data"

    if FORCE_DOWNLOAD:
        print("Force download enabled, removing cached snapshot...")
        # Remove contents of cache directory (can't remove the dir itself as it's a mount point)
        run_command(f'rm -rf "{SNAPSHOT_CACHE}"/*', check=False)

    if dir_exists_and_not_empty(snapshot_data_dir):
        # Snapshot cache exists, copy if needed
        if not dir_exists_and_not_empty(zetacored_data_dir):
            print("Copying snapshot from cache...")
            run_command(f'cp -r "{snapshot_data_dir}" "{ZETACORED_HOME}/"')
    else:
        # Download snapshot
        print("Downloading snapshot...")
        run_command("python3 -u /root/download_snapshot.py")

        if os.path.isdir(snapshot_data_dir):
            if os.path.isdir(zetacored_data_dir):
                run_command(f'rm -rf "{zetacored_data_dir}"')
            run_command(f'cp -r "{snapshot_data_dir}" "{ZETACORED_HOME}/"')


def get_external_ip():
    """Get external IP address."""
    result = subprocess.run(
        "curl -4 -s icanhazip.com",
        shell=True,
        capture_output=True,
        text=True
    )
    if result.returncode == 0 and result.stdout.strip():
        return result.stdout.strip()
    return "127.0.0.1"


def update_configs():
    """Update configuration files."""
    print("Updating configuration files...")

    external_ip = get_external_ip()
    print(f"External IP: {external_ip}")

    config_toml = f"{ZETACORED_CONFIG}/config.toml"
    app_toml = f"{ZETACORED_CONFIG}/app.toml"

    # Update config.toml
    sed_replace(config_toml, "{YOUR_EXTERNAL_IP_ADDRESS_HERE}", external_ip)
    sed_replace(config_toml, "{MONIKER}", MONIKER)
    sed_replace(config_toml, 'laddr = "tcp:\\/\\/127.0.0.1:26657"', 'laddr = "tcp:\\/\\/0.0.0.0:26657"')
    sed_replace(config_toml, "prometheus = false", "prometheus = true")

    # Update app.toml
    sed_replace(app_toml, "enable = false", "enable = true")
    sed_replace(app_toml, 'address = "tcp:\\/\\/localhost:1317"', 'address = "tcp:\\/\\/0.0.0.0:1317"')
    sed_replace(app_toml, 'address = "localhost:9090"', 'address = "0.0.0.0:9090"')


def start_zetacored():
    """Start zetacored."""
    print("Starting zetacored...")
    # Use absolute path to avoid PATH conflicts with zetaclientd symlinks
    os.execv("/usr/local/bin/zetacored", ["/usr/local/bin/zetacored", "start"])


def main():
    """Main execution."""
    print("========================================")
    print("  Starting ZetaCore Testnet Node       ")
    print("========================================")

    initialize_node()
    setup_snapshot()
    update_configs()
    start_zetacored()


if __name__ == "__main__":
    main()
