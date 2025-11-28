#!/usr/bin/env python3
"""
Download Snapshot Script
========================
This script downloads the latest snapshot for a ZetaChain network and caches it locally
for reuse in multiple devnet fork runs.

Usage:
    python3 contrib/localnet/scripts_python/download_snapshot.py [--chain-id CHAIN_ID]

Examples:
    python3 contrib/localnet/scripts_python/download_snapshot.py --chain-id athens_7001-1
    python3 contrib/localnet/scripts_python/download_snapshot.py --chain-id zetachain_7000-1
"""

import argparse
import hashlib
import os
import subprocess
import sys
import requests
from pathlib import Path

# ============================================================================
# Configuration Constants
# ============================================================================

HOME_DIR = Path.home()
TEMP_EXTRACT_DIR = HOME_DIR / "zetacored_snapshot_temp"

NETWORK_CONFIGS = {
    "athens_7001-1": {
        "snapshot_url": "https://snapshots.rpc.zetachain.com/testnet/fullnode/latest.json",
        "cache_dir": HOME_DIR / "zetacored_snapshot_testnet",
        "name": "Testnet",
    },
    "zetachain_7000-1": {
        "snapshot_url": "https://snapshots.rpc.zetachain.com/mainnet/fullnode/latest.json",
        "cache_dir": HOME_DIR / "zetacored_snapshot_mainnet",
        "name": "Mainnet",
    },
}

# Default chain ID
DEFAULT_CHAIN_ID = "athens_7001-1"

# ============================================================================
# Helper Functions
# ============================================================================

def run_command(cmd, shell=True, check=True, silent=False):
    """Run a shell command and handle errors."""
    if not silent:
        print(f"Running: {cmd}")
    try:
        subprocess.run(cmd, shell=shell, check=check)
    except subprocess.CalledProcessError as e:
        print(f"Error running command: {cmd}")
        print(f"Error: {e}")
        sys.exit(1)


def download_file(url, dest_path):
    """Download file."""
    response = requests.get(url, stream=True, timeout=30)
    response.raise_for_status()

    with open(dest_path, 'wb') as f:
        for chunk in response.iter_content(chunk_size=8192 * 128):
            if chunk:
                f.write(chunk)


def extract_archive(archive_path, dest_dir):
    """Extract lz4 archive."""
    print(f"  Extracting archive...")
    subprocess.run(
        f'lz4 -dc "{archive_path}" | tar -C "{dest_dir}/" -xf -',
        shell=True,
        check=True
    )
    print(f"  Extraction complete.")


def compute_md5(file_path):
    """Compute MD5 checksum."""
    md5_hash = hashlib.md5()
    with open(file_path, "rb") as f:
        for chunk in iter(lambda: f.read(8192 * 128), b""):
            md5_hash.update(chunk)
    return md5_hash.hexdigest()


def verify_checksum(file_path, expected_md5):
    """Verify MD5 checksum of downloaded file."""
    print("  Computing MD5 checksum...")
    computed_md5 = compute_md5(file_path)

    if computed_md5 == expected_md5:
        print("  ✓ Checksum verification passed!")
        return True
    else:
        print("  ✗ Checksum verification FAILED!")
        print(f"    Expected: {expected_md5}")
        print(f"    Got:      {computed_md5}")
        return False


def get_network_config(chain_id):
    """Get network configuration for the given chain ID."""
    if chain_id not in NETWORK_CONFIGS:
        print(f"Error: Unknown chain ID: {chain_id}")
        print(f"Supported chain IDs: {', '.join(NETWORK_CONFIGS.keys())}")
        sys.exit(1)
    return NETWORK_CONFIGS[chain_id]


# ============================================================================
# Main Script
# ============================================================================

def main(chain_id, force=False):
    config = get_network_config(chain_id)
    snapshot_url = config["snapshot_url"]
    snapshot_cache_dir = config["cache_dir"]
    network_name = config["name"]

    force = force or os.environ.get('FORCE_DOWNLOAD', '').lower() == 'true'

    print("=" * 60)
    print(f"  ZetaChain {network_name} Snapshot Download")
    print(f"  Chain ID: {chain_id}")
    print("=" * 60)

    snapshot_data_dir = snapshot_cache_dir / "data"
    if snapshot_data_dir.exists() and any(snapshot_data_dir.iterdir()):
        if not force:
            print(f"\nWarning: Cached snapshot already exists at {snapshot_cache_dir}")
            response = input("Do you want to re-download and overwrite? (yes/no): ")
            if response.lower() not in ['yes', 'y']:
                print("Aborted. Using existing cache.")
                sys.exit(0)
        print("Removing existing cache...")
        run_command(f'rm -rf "{snapshot_cache_dir}"/*', silent=True)

    # Create cache directory if it doesn't exist
    snapshot_cache_dir.mkdir(parents=True, exist_ok=True)

    # Fetch snapshot info
    print("\n[1/5] Fetching snapshot information...")
    try:
        snapshot_json = requests.get(snapshot_url, timeout=30).json()
        snapshot_data = snapshot_json['snapshots'][0]
        snapshot_link = snapshot_data['link']
        snapshot_filename = snapshot_data['filename']
        expected_md5 = snapshot_data.get('checksums', {}).get('md5')
        print(f"  Snapshot: {snapshot_filename}")
    except Exception as e:
        print(f"  Error fetching snapshot info: {e}")
        sys.exit(1)

    snapshot_path = HOME_DIR / snapshot_filename

    # Download snapshot
    print("\n[2/5] Downloading snapshot...")
    try:
        download_file(snapshot_link, snapshot_path)
        print("  ✓ Download complete!")
    except Exception as e:
        print(f"\n  Error downloading snapshot: {e}")
        sys.exit(1)

    # Verify checksum
    print("\n[3/5] Verifying snapshot integrity...")
    if expected_md5:
        if not verify_checksum(snapshot_path, expected_md5):
            print("\n  Checksum verification failed. Aborting...")
            snapshot_path.unlink()
            sys.exit(1)
    else:
        print("  Skipping checksum verification (no checksum available)")

    # Create temp directory and extract
    print("\n[4/5] Extracting snapshot...")
    if TEMP_EXTRACT_DIR.exists():
        run_command(f'rm -rf "{TEMP_EXTRACT_DIR}"', silent=True)
    TEMP_EXTRACT_DIR.mkdir(parents=True, exist_ok=True)

    extract_archive(snapshot_path, TEMP_EXTRACT_DIR)

    # Move data directory to cache location
    print("\n[5/5] Moving snapshot to cache...")
    temp_data_dir = TEMP_EXTRACT_DIR / "data"
    if not temp_data_dir.exists():
        print(f"  Error: Expected data directory not found at {temp_data_dir}")
        sys.exit(1)

    # Move the data directory INTO the cache location (cache_dir/data/)
    target_data_dir = snapshot_cache_dir / "data"
    if target_data_dir.exists():
        run_command(f'rm -rf "{target_data_dir}"', silent=True)
    run_command(f'mv "{temp_data_dir}" "{target_data_dir}"', silent=True)

    # Cleanup
    print("  Cleaning up temporary files...")
    run_command(f'rm -rf "{TEMP_EXTRACT_DIR}"', silent=True)
    snapshot_path.unlink()

    print("\n" + "=" * 60)
    print("  ✓ Snapshot cached successfully!")
    print(f"  Location: {snapshot_cache_dir}")
    print("=" * 60)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Download and cache ZetaChain snapshot",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    parser.add_argument(
        "--chain-id",
        default=DEFAULT_CHAIN_ID,
        help=f"Chain ID to download snapshot for (default: {DEFAULT_CHAIN_ID})"
    )
    parser.add_argument(
        "--force",
        action="store_true",
        help="Force re-download even if cache exists"
    )
    args = parser.parse_args()

    main(args.chain_id, args.force)
