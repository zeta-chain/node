#!/usr/bin/env python3
"""
Devnet Fork Script
==================
This script downloads testnet state, syncs for a period, then converts it to a local
single-validator devnet for testing purposes.

Run from the root zeta-node directory:
    python3 contrib/devnet/devnet_fork.py [--node-version VERSION]
Or:
    make devnet-fork
"""

import argparse
import os
import subprocess
import sys
import time
import json
import requests
from pathlib import Path

# ============================================================================
# Configuration Constants
# ============================================================================

# Devnet configuration
DEVNET_CHAIN_ID = "devnet_70000-1"
OPERATOR_ADDRESS = "zeta13l7ladn2crrdcl9nupqn5kzyajcn03lkzgnrze"

# Athens3 testnet config URLs
ATHENS3_CONFIG_BASE = "https://raw.githubusercontent.com/zeta-chain/network-config/main/athens3"
ATHENS3_GENESIS_URL = f"{ATHENS3_CONFIG_BASE}/genesis.json"
ATHENS3_CLIENT_URL = f"{ATHENS3_CONFIG_BASE}/client.toml"
ATHENS3_CONFIG_URL = f"{ATHENS3_CONFIG_BASE}/config.toml"
ATHENS3_APP_URL = f"{ATHENS3_CONFIG_BASE}/app.toml"

# Snapshot configuration
SNAPSHOT_JSON_URL = "https://snapshots.rpc.zetachain.com/testnet/fullnode/latest.json"

# Timing configuration
SYNC_DURATION_SECONDS = 300  # 5 minutes
SHUTDOWN_WAIT_SECONDS = 5

# Paths
HOME_DIR = Path.home()
ZETACORED_DIR = HOME_DIR / ".zetacored"
ZETACORED_CONFIG_DIR = ZETACORED_DIR / "config"
ZETACORED_LOG_FILE = HOME_DIR / "zetacored_devnet_fork.log"

# ============================================================================
# Helper Functions
# ============================================================================

def run_command(cmd, shell=True, check=True, capture_output=False):
    """Run a shell command and handle errors."""
    print(f"Running: {cmd}")
    try:
        result = subprocess.run(
            cmd,
            shell=shell,
            check=check,
            capture_output=capture_output,
            text=True
        )
        if capture_output:
            return result.stdout.strip()
        return result
    except subprocess.CalledProcessError as e:
        print(f"Error running command: {cmd}")
        print(f"Error: {e}")
        sys.exit(1)

def download_file(url, dest_path):
    """Download a file from URL to destination path."""
    print(f"Downloading {url} to {dest_path}")
    try:
        response = requests.get(url, timeout=300)
        response.raise_for_status()
        with open(dest_path, 'wb') as f:
            f.write(response.content)
        print(f"Successfully downloaded {dest_path}")
    except Exception as e:
        print(f"Error downloading {url}: {e}")
        sys.exit(1)

def get_external_ip():
    """Get external IP address."""
    print("Getting external IP address...")
    try:
        # Force IPv4 to match the shell script behavior (curl -4)
        response = requests.get("https://ipv4.icanhazip.com", timeout=10)
        response.raise_for_status()
        ip = response.text.strip()
        print(f"External IP: {ip}")
        return ip
    except Exception as e:
        print(f"Error getting external IP: {e}")
        sys.exit(1)

def replace_in_file(file_path, search, replace):
    """Replace text in a file."""
    print(f"Updating {file_path}: replacing '{search}' with '{replace}'")
    try:
        with open(file_path, 'r') as f:
            content = f.read()
        content = content.replace(search, replace)
        with open(file_path, 'w') as f:
            f.write(content)
    except Exception as e:
        print(f"Error updating {file_path}: {e}")
        sys.exit(1)

def extract_major_version(version):
    """Extract major version from a version string (e.g., 'v36.0.4' -> 'v36')."""
    # Remove 'v' prefix if present
    if version.startswith('v'):
        version_without_v = version[1:]
    else:
        version_without_v = version

    # Split by '.' and take first part
    major = version_without_v.split('.')[0]

    # Return with 'v' prefix
    return f"v{major}"

def setup_cosmovisor(node_version, upgrade_version=None):
    """Set up Cosmovisor directory structure and binaries."""
    print("\n[Setting up Cosmovisor]")

    # Create cosmovisor directory structure
    cosmovisor_dir = ZETACORED_DIR / "cosmovisor"
    genesis_bin_dir = cosmovisor_dir / "genesis" / "bin"

    print(f"Creating genesis bin directory: {genesis_bin_dir}")
    genesis_bin_dir.mkdir(parents=True, exist_ok=True)

    # Install genesis version (node_version) and copy to genesis bin
    print(f"\nInstalling genesis version: {node_version}")
    run_command(f"NODE_VERSION={node_version} make install")

    # Find the zetacored binary path
    zetacored_path = run_command("which zetacored", capture_output=True)
    print(f"Found zetacored at: {zetacored_path}")

    # Copy to genesis bin
    genesis_binary = genesis_bin_dir / "zetacored"
    print(f"Copying binary to: {genesis_binary}")
    run_command(f"cp {zetacored_path} {genesis_binary}")
    run_command(f"chmod +x {genesis_binary}")

    print(f"\nCosmovisor genesis setup complete!")
    print(f"  Genesis binary: {genesis_binary}")

    # Only set up upgrade directory if upgrade_version is provided
    if upgrade_version:
        # Create upgrade directory
        upgrade_handler_version = extract_major_version(upgrade_version)
        upgrade_bin_dir = cosmovisor_dir / "upgrades" / upgrade_handler_version / "bin"
        print(f"\nCreating upgrade bin directory: {upgrade_bin_dir}")
        upgrade_bin_dir.mkdir(parents=True, exist_ok=True)

        # Install upgrade version and copy to upgrade bin
        print(f"Installing upgrade version: {upgrade_version}")
        run_command(f"NODE_VERSION={upgrade_version} make install")

        # Find the upgraded zetacored binary path
        zetacored_path = run_command("which zetacored", capture_output=True)
        print(f"Found upgraded zetacored at: {zetacored_path}")

        # Copy to upgrade bin
        upgrade_binary = upgrade_bin_dir / "zetacored"
        print(f"Copying binary to: {upgrade_binary}")
        run_command(f"cp {zetacored_path} {upgrade_binary}")
        run_command(f"chmod +x {upgrade_binary}")

        print(f"\nCosmovisor upgrade setup complete!")
        print(f"  Upgrade binary: {upgrade_binary}")
        print(f"  Upgrade handler: {upgrade_handler_version}")

    # Set environment variables for Cosmovisor
    print("\nSetting up Cosmovisor environment variables...")
    env_vars = {
        "DAEMON_HOME": str(ZETACORED_DIR),
        "DAEMON_NAME": "zetacored",
        "DAEMON_ALLOW_DOWNLOAD_BINARIES": "true",
        "DAEMON_RESTART_AFTER_UPGRADE": "true",
        "CLIENT_DAEMON_NAME": "zetaclientd",
        "CLIENT_DAEMON_ARGS": "-enable-chains,GOERLI,-val,operator",
        "DAEMON_DATA_BACKUP_DIR": str(ZETACORED_DIR),
        "CLIENT_SKIP_UPGRADE": "true",
        "CLIENT_START_PROCESS": "false",
        "UNSAFE_SKIP_BACKUP": "true"
    }

    for key, value in env_vars.items():
        os.environ[key] = value
        print(f"  export {key}={value}")

# ============================================================================
# Main Script
# ============================================================================

def main(node_version, upgrade_version=None):
    print("=" * 80)
    print("ZetaChain Devnet Fork Script")
    print("=" * 80)
    print(f"Using node version: {node_version}")
    if upgrade_version:
        print(f"Upgrade version: {upgrade_version}")
    print("=" * 80)

    # Step 1: Clean and build
    print("\n[1/10] Cleaning and building zetacored...")
    run_command("make clean")
    run_command(f"NODE_VERSION={node_version} make install")

    # Step 2: Initialize node
    print("\n[2/10] Initializing zetacored...")
    run_command("zetacored init test --chain-id=athens_7001-1")

    # Step 3: Download config files
    print("\n[3/10] Downloading Athens3 testnet configuration...")
    download_file(ATHENS3_GENESIS_URL, ZETACORED_CONFIG_DIR / "genesis.json")
    download_file(ATHENS3_CLIENT_URL, ZETACORED_CONFIG_DIR / "client.toml")
    download_file(ATHENS3_CONFIG_URL, ZETACORED_CONFIG_DIR / "config.toml")
    download_file(ATHENS3_APP_URL, ZETACORED_CONFIG_DIR / "app.toml")

    # Step 4: Update config with external IP and moniker
    print("\n[4/10] Updating configuration with external IP...")
    external_ip = get_external_ip()
    config_file = ZETACORED_CONFIG_DIR / "config.toml"
    replace_in_file(config_file, "{YOUR_EXTERNAL_IP_ADDRESS_HERE}", external_ip)
    replace_in_file(config_file, "{MONIKER}", "testNode")

    # Step 5: Download and extract snapshot
    print("\n[5/10] Loading testnet snapshot...")

    # Check for local snapshot cache
    snapshot_cache_dir = HOME_DIR / "zetacored_snapshot_testnet"
    use_cache = False

    if snapshot_cache_dir.exists() and any(snapshot_cache_dir.iterdir()):
        print(f"Found cached snapshot at: {snapshot_cache_dir}")
        print("Using cached snapshot instead of downloading...")
        use_cache = True
    else:
        print("No cached snapshot found, will download from remote...")

    if use_cache:
        # Copy from cache
        print(f"Copying cached snapshot data to {ZETACORED_DIR}/data...")
        data_dir = ZETACORED_DIR / "data"
        if data_dir.exists():
            run_command(f'rm -rf "{data_dir}"')
        run_command(f'cp -r "{snapshot_cache_dir}" "{data_dir}"')
        print("Cached snapshot copied successfully!")
    else:
        # Download snapshot
        try:
            snapshot_json = requests.get(SNAPSHOT_JSON_URL, timeout=30).json()
            snapshot_link = snapshot_json['snapshots'][0]['link']
            snapshot_filename = snapshot_json['snapshots'][0]['filename']
            print(f"Snapshot: {snapshot_filename}")
            print(f"Link: {snapshot_link}")
        except Exception as e:
            print(f"Error fetching snapshot info: {e}")
            sys.exit(1)

        snapshot_path = HOME_DIR / snapshot_filename

        # Download snapshot
        print(f"Downloading snapshot to {snapshot_path}...")
        run_command(f'curl "{snapshot_link}" -o "{snapshot_path}"')

        # Extract snapshot
        print("Extracting snapshot...")
        run_command(f'lz4 -dc "{snapshot_path}" | tar -C "{ZETACORED_DIR}/" -xvf -')

        # Remove snapshot file
        print("Cleaning up snapshot file...")
        snapshot_path.unlink()

    # Step 6: Initial sleep
    print("\n[6/10] Waiting 10 seconds before starting node...")
    time.sleep(10)

    # Step 7: Start zetacored in background
    print("\n[7/10] Starting zetacored in background...")
    # Open log file for writing
    log_file = open(ZETACORED_LOG_FILE, 'w')
    process = subprocess.Popen(
        ["zetacored", "start"],
        stdout=log_file,
        stderr=subprocess.STDOUT  # Redirect stderr to stdout
    )
    print(f"Started zetacored with PID: {process.pid}")
    print(f"Logs are being written to: {ZETACORED_LOG_FILE}")

    # Start tailing the log file in background
    print("Starting log tail in background...")
    subprocess.Popen(
        ["tail", "-f", str(ZETACORED_LOG_FILE)],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL
    )

    # Wait for RPC server to start
    print("Waiting 30 seconds for node to start up...")
    time.sleep(30)

    # Poll RPC endpoint until sync is complete
    print("Polling node sync status...")
    rpc_url = "http://localhost:26657/status"
    max_wait_time = 600  # 10 minutes timeout
    start_time = time.time()
    poll_counter = 0

    while time.time() - start_time < max_wait_time:
        try:
            response = requests.get(rpc_url, timeout=5)
            if response.status_code == 200:
                status = response.json()
                sync_info = status.get("result", {}).get("sync_info", {})
                catching_up = sync_info.get("catching_up", True)
                latest_height = sync_info.get("latest_block_height", "unknown")

                # Print status every 30 seconds
                if poll_counter % 30 == 0:
                    catching_up_str = "Yes" if catching_up else "No"
                    print(f"Sync status - height: {latest_height}, still syncing: {catching_up_str}")

                # Check if sync is complete (catching_up == False)
                if catching_up == False:
                    print("Node sync complete!")
                    break
            else:
                if poll_counter % 30 == 0:
                    print(f"RPC returned status code {response.status_code}, retrying...")
        except Exception as e:
            if poll_counter % 30 == 0:
                print(f"RPC not ready yet: {e}")

        poll_counter += 1
        time.sleep(1)
    else:
        print(f"Warning: Reached maximum wait time of {max_wait_time} seconds")

    # Wait for node to produce some blocks
    print("Waiting 30 seconds for node to produce blocks...")
    time.sleep(30)

    # Step 8: Stop zetacored
    print("\n[8/10] Stopping zetacored...")
    process.terminate()
    try:
        process.wait(timeout=SHUTDOWN_WAIT_SECONDS)
        print("Process terminated gracefully")
    except subprocess.TimeoutExpired:
        print("Process didn't stop gracefully, killing...")
        process.kill()
        process.wait()

    # Close the log file
    log_file.close()
    print(f"Logs saved to: {ZETACORED_LOG_FILE}")

    print(f"Waiting {SHUTDOWN_WAIT_SECONDS} seconds for clean shutdown...")
    time.sleep(SHUTDOWN_WAIT_SECONDS)

    # Step 9: Run devnet command
    print("\n[9/10] Running devnet command to modify state...")
    devnet_cmd = f"zetacored devnet {DEVNET_CHAIN_ID} {OPERATOR_ADDRESS} --skip-confirmation"
    if upgrade_version:
        # Extract major version (e.g., v36.0.4 -> v36) to match handler name
        upgrade_handler_version = extract_major_version(upgrade_version)
        devnet_cmd += f" --upgrade-version {upgrade_handler_version}"
        print(f"Scheduling upgrade to version: {upgrade_version} (handler: {upgrade_handler_version})")

    # Run devnet command in background (for testing)
    print(f"Running: {devnet_cmd}")
    test_log_file = open(HOME_DIR / "zetacored_devnet.log", 'w')
    test_process = subprocess.Popen(
        devnet_cmd,
        shell=True,
        stdout=test_log_file,
        stderr=subprocess.STDOUT
    )
    print(f"Started zetacored devnet with PID: {test_process.pid}")
    print("Waiting 30 seconds...")
    time.sleep(30)

    print("Killing devnet node...")
    test_process.kill()
    test_process.wait()
    test_log_file.close()
    print(f"Devnet logs: {HOME_DIR / 'zetacored_devnet.log'}")

    # Step 10: Set up Cosmovisor and start
    print("\n[10/10] Setting up Cosmovisor and starting node...")
    setup_cosmovisor(node_version, upgrade_version)

    # Run cosmovisor start in foreground
    print("\nStarting Cosmovisor...")
    cosmovisor_cmd = "cosmovisor start"
    print(f"Running: {cosmovisor_cmd}")
    print("Note: Cosmovisor is running in foreground. Press Ctrl+C to stop.")
    run_command(cosmovisor_cmd)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Fork testnet state to create a local single-validator devnet",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    parser.add_argument(
        "--node-version",
        required=True,
        help="Node version to build (e.g., 36.0.1)"
    )
    parser.add_argument(
        "--upgrade-version",
        default=None,
        help="Schedule an upgrade to this version (e.g., v36.0.4). If not provided, no upgrade is scheduled."
    )
    args = parser.parse_args()

    try:
        main(args.node_version, args.upgrade_version)
    except KeyboardInterrupt:
        print("\n\nScript interrupted by user. Exiting...")
        sys.exit(1)
    except Exception as e:
        print(f"\n\nUnexpected error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
