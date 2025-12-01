#!/usr/bin/env python3
"""
Download Testnet Snapshot Script
==================================
This script downloads the latest testnet snapshot and caches it locally
for reuse in multiple devnet fork runs.

The cached snapshot is stored in ~/zetacored_snapshot_testnet/

Run from the root zeta-node directory:
    python3 contrib/devnet/download_snapshot.py
"""

import os
import subprocess
import sys
import hashlib
import requests
from pathlib import Path

# ============================================================================
# Configuration Constants
# ============================================================================

# Snapshot configuration
SNAPSHOT_JSON_URL = "https://snapshots.rpc.zetachain.com/testnet/fullnode/latest.json"

# Paths
HOME_DIR = Path.home()
SNAPSHOT_CACHE_DIR = HOME_DIR / "zetacored_snapshot_testnet"
TEMP_EXTRACT_DIR = HOME_DIR / "zetacored_snapshot_temp"

# ============================================================================
# Helper Functions
# ============================================================================

def run_command(cmd, shell=True, check=True):
    """Run a shell command and handle errors."""
    print(f"Running: {cmd}")
    try:
        subprocess.run(cmd, shell=shell, check=check)
    except subprocess.CalledProcessError as e:
        print(f"Error running command: {cmd}")
        print(f"Error: {e}")
        sys.exit(1)

def compute_md5(file_path, chunk_size=8192):
    """Compute MD5 checksum of a file."""
    md5_hash = hashlib.md5()
    with open(file_path, "rb") as f:
        for chunk in iter(lambda: f.read(chunk_size), b""):
            md5_hash.update(chunk)
    return md5_hash.hexdigest()

def verify_checksum(file_path, expected_md5):
    """Verify MD5 checksum of downloaded file."""
    print(f"Computing MD5 checksum of {file_path.name}...")
    print("This may take a few minutes for large files...")
    computed_md5 = compute_md5(file_path)
    print(f"Expected MD5:  {expected_md5}")
    print(f"Computed MD5:  {computed_md5}")

    if computed_md5 == expected_md5:
        print("✓ Checksum verification passed!")
        return True
    else:
        print("✗ Checksum verification FAILED!")
        print("The downloaded file may be corrupted or tampered with.")
        return False

# ============================================================================
# Main Script
# ============================================================================

def main():
    print("=" * 80)
    print("ZetaChain Testnet Snapshot Download Script")
    print("=" * 80)
    print(f"Cache location: {SNAPSHOT_CACHE_DIR}")
    print("=" * 80)

    # Check if cache already exists
    if SNAPSHOT_CACHE_DIR.exists() and any(SNAPSHOT_CACHE_DIR.iterdir()):
        print(f"\nWarning: Cached snapshot already exists at {SNAPSHOT_CACHE_DIR}")
        response = input("Do you want to re-download and overwrite? (yes/no): ")
        if response.lower() not in ['yes', 'y']:
            print("Aborted. Using existing cache.")
            sys.exit(0)
        print("Removing existing cache...")
        run_command(f'rm -rf "{SNAPSHOT_CACHE_DIR}"')

    # Fetch snapshot info
    print("\n[1/5] Fetching snapshot information...")
    try:
        snapshot_json = requests.get(SNAPSHOT_JSON_URL, timeout=30).json()
        snapshot_data = snapshot_json['snapshots'][0]
        snapshot_link = snapshot_data['link']
        snapshot_filename = snapshot_data['filename']
        expected_md5 = snapshot_data.get('checksums', {}).get('md5')
        print(f"Snapshot: {snapshot_filename}")
        print(f"Link: {snapshot_link}")
        if expected_md5:
            print(f"Expected MD5: {expected_md5}")
        else:
            print("Warning: No MD5 checksum available for verification")
    except Exception as e:
        print(f"Error fetching snapshot info: {e}")
        sys.exit(1)

    snapshot_path = HOME_DIR / snapshot_filename

    # Download snapshot
    print(f"\n[2/5] Downloading snapshot to {snapshot_path}...")
    print("This may take a while depending on your internet connection...")
    run_command(f'curl "{snapshot_link}" -o "{snapshot_path}"')
    print("Download complete!")

    # Verify checksum
    if expected_md5:
        print(f"\n[3/5] Verifying snapshot integrity...")
        if not verify_checksum(snapshot_path, expected_md5):
            print("\nChecksum verification failed. Aborting...")
            print(f"Removing corrupted file: {snapshot_path}")
            snapshot_path.unlink()
            sys.exit(1)
    else:
        print(f"\n[3/5] Skipping checksum verification (no checksum available)")

    # Create temp directory and extract
    print(f"\n[4/5] Extracting snapshot to temporary location...")
    if TEMP_EXTRACT_DIR.exists():
        run_command(f'rm -rf "{TEMP_EXTRACT_DIR}"')
    TEMP_EXTRACT_DIR.mkdir(parents=True, exist_ok=True)

    print("Extracting... (this may take several minutes)")
    run_command(f'lz4 -dc "{snapshot_path}" | tar -C "{TEMP_EXTRACT_DIR}/" -xvf -')

    # Move data directory to cache location
    print(f"\n[5/5] Moving snapshot data to cache location...")
    temp_data_dir = TEMP_EXTRACT_DIR / "data"
    if not temp_data_dir.exists():
        print(f"Error: Expected data directory not found at {temp_data_dir}")
        sys.exit(1)

    # Move the data directory to the cache location
    run_command(f'mv "{temp_data_dir}" "{SNAPSHOT_CACHE_DIR}"')

    # Cleanup
    print("\nCleaning up temporary files...")
    run_command(f'rm -rf "{TEMP_EXTRACT_DIR}"')
    print(f"Removing snapshot archive: {snapshot_path}")
    snapshot_path.unlink()

    print("\n" + "=" * 80)
    print("Snapshot cached successfully!")
    print(f"Cache location: {SNAPSHOT_CACHE_DIR}")
    print("=" * 80)

if __name__ == "__main__":
    main()
