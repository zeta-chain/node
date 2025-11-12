#!/usr/bin/env python3
"""
Devnet Fork Script
==================
This script downloads testnet state, syncs for a period, then converts it to a local
single-validator devnet for testing purposes.

Run from the root zeta-node directory:
    python3 contrib/devnet/devnet_fork.py
Or:
    make devnet-fork
"""

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

# Node version to build
NODE_VERSION = "36.0.1"

# Devnet configuration (forked from Athens3 testnet)
DEVNET_CHAIN_ID = "devnet_70000-1"
OPERATOR_ADDRESS = "zeta13l7ladn2crrdcl9nupqn5kzyajcn03lkzgnrze"

# Athens3 testnet config URLs (source network for fork)
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

# ============================================================================
# Main Script
# ============================================================================

def main():
    print("=" * 80)
    print("ZetaChain Devnet Fork Script")
    print("=" * 80)

    # Step 1: Clean and build
    print("\n[1/9] Cleaning and building zetacored...")
    run_command("make clean")
    run_command(f"NODE_VERSION={NODE_VERSION} make install")

    # Step 2: Initialize node
    print("\n[2/9] Initializing zetacored...")
    run_command("zetacored init test --chain-id=athens_7001-1")

    # Step 3: Download config files
    print("\n[3/9] Downloading Athens3 configuration...")
    download_file(ATHENS3_GENESIS_URL, ZETACORED_CONFIG_DIR / "genesis.json")
    download_file(ATHENS3_CLIENT_URL, ZETACORED_CONFIG_DIR / "client.toml")
    download_file(ATHENS3_CONFIG_URL, ZETACORED_CONFIG_DIR / "config.toml")
    download_file(ATHENS3_APP_URL, ZETACORED_CONFIG_DIR / "app.toml")

    # Step 4: Update config with external IP and moniker
    print("\n[4/9] Updating configuration with external IP...")
    external_ip = get_external_ip()
    config_file = ZETACORED_CONFIG_DIR / "config.toml"
    replace_in_file(config_file, "{YOUR_EXTERNAL_IP_ADDRESS_HERE}", external_ip)
    replace_in_file(config_file, "{MONIKER}", "devNode")

    # Step 5: Download and extract snapshot
    print("\n[5/9] Downloading snapshot...")
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
    print("\n[6/9] Waiting 10 seconds before starting node...")
    time.sleep(10)

    # Step 7: Start zetacored in background
    print("\n[7/9] Starting zetacored in background...")
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
    print("\n[8/9] Stopping zetacored...")
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
    print("\n[9/9] Running devnet command to modify state and start the network")
    devnet_cmd = f"zetacored devnet {DEVNET_CHAIN_ID} {OPERATOR_ADDRESS} --skip-confirmation"
    run_command(devnet_cmd)

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\nScript interrupted by user. Exiting...")
        sys.exit(1)
    except Exception as e:
        print(f"\n\nUnexpected error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
