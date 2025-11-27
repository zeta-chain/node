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


def format_size(size_bytes):
    """Format bytes to human readable size."""
    for unit in ['B', 'KB', 'MB', 'GB']:
        if size_bytes < 1024:
            return f"{size_bytes:.1f} {unit}"
        size_bytes /= 1024
    return f"{size_bytes:.1f} TB"


def download_with_progress(url, dest_path):
    """Download file with progress indicator."""
    response = requests.get(url, stream=True, timeout=30)
    response.raise_for_status()

    total_size = int(response.headers.get('content-length', 0))
    downloaded = 0
    chunk_size = 8192 * 128  # 1MB chunks
    last_percent_printed = -10

    with open(dest_path, 'wb') as f:
        for chunk in response.iter_content(chunk_size=chunk_size):
            if chunk:
                f.write(chunk)
                downloaded += len(chunk)

                if total_size > 0:
                    percent = int((downloaded / total_size) * 100)
                    # Print every 10%
                    if percent >= last_percent_printed + 10:
                        last_percent_printed = (percent // 10) * 10
                        print(f"  {percent}% ({format_size(downloaded)}/{format_size(total_size)})")


def extract_with_progress(archive_path, dest_dir):
    """Extract lz4 archive with progress indication."""
    # Get archive size for progress estimation
    archive_size = archive_path.stat().st_size

    # Use pv if available for progress, otherwise show spinner
    pv_check = subprocess.run("which pv", shell=True, capture_output=True)

    if pv_check.returncode == 0:
        # pv is available - use it for progress
        cmd = f'pv -p -e "{archive_path}" | lz4 -dc | tar -C "{dest_dir}/" -xf -'
        subprocess.run(cmd, shell=True, check=True)
    else:
        # No pv - show a simple progress indicator
        print(f"  Extracting {format_size(archive_size)} archive...")

        # Run extraction in background and show spinner
        process = subprocess.Popen(
            f'lz4 -dc "{archive_path}" | tar -C "{dest_dir}/" -xf -',
            shell=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )

        spinner = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏']
        i = 0
        while process.poll() is None:
            print(f"\r  {spinner[i % len(spinner)]} Extracting...", end='', flush=True)
            i += 1
            try:
                process.wait(timeout=0.1)
            except subprocess.TimeoutExpired:
                pass

        if process.returncode != 0:
            print("\r  ✗ Extraction failed!")
            stderr = process.stderr.read().decode() if process.stderr else ""
            if stderr:
                print(f"  Error: {stderr}")
            sys.exit(1)

        print("\r  ✓ Extraction complete!    ")


def compute_md5_with_progress(file_path, chunk_size=8192 * 128):
    """Compute MD5 checksum with progress indicator."""
    md5_hash = hashlib.md5()
    file_size = file_path.stat().st_size
    processed = 0
    last_percent_printed = -10

    with open(file_path, "rb") as f:
        for chunk in iter(lambda: f.read(chunk_size), b""):
            md5_hash.update(chunk)
            processed += len(chunk)

            percent = int((processed / file_size) * 100)
            # Print every 10%
            if percent >= last_percent_printed + 10:
                last_percent_printed = (percent // 10) * 10
                print(f"  {percent}%")

    return md5_hash.hexdigest()


def verify_checksum(file_path, expected_md5):
    """Verify MD5 checksum of downloaded file."""
    print(f"  Computing MD5 checksum...")
    computed_md5 = compute_md5_with_progress(file_path)

    if computed_md5 == expected_md5:
        print("  ✓ Checksum verification passed!")
        return True
    else:
        print("  ✗ Checksum verification FAILED!")
        print(f"    Expected: {expected_md5}")
        print(f"    Got:      {computed_md5}")
        return False

# ============================================================================
# Main Script
# ============================================================================

def main():
    print("=" * 60)
    print("  ZetaChain Testnet Snapshot Download")
    print("=" * 60)

    # Check if cache already exists
    if SNAPSHOT_CACHE_DIR.exists() and any(SNAPSHOT_CACHE_DIR.iterdir()):
        print(f"\nWarning: Cached snapshot already exists at {SNAPSHOT_CACHE_DIR}")
        response = input("Do you want to re-download and overwrite? (yes/no): ")
        if response.lower() not in ['yes', 'y']:
            print("Aborted. Using existing cache.")
            sys.exit(0)
        print("Removing existing cache...")
        run_command(f'rm -rf "{SNAPSHOT_CACHE_DIR}"', silent=True)

    # Fetch snapshot info
    print("\n[1/5] Fetching snapshot information...")
    try:
        snapshot_json = requests.get(SNAPSHOT_JSON_URL, timeout=30).json()
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
    print(f"\n[2/5] Downloading snapshot...")
    try:
        download_with_progress(snapshot_link, snapshot_path)
        print("  ✓ Download complete!")
    except Exception as e:
        print(f"\n  Error downloading snapshot: {e}")
        sys.exit(1)

    # Verify checksum
    print(f"\n[3/5] Verifying snapshot integrity...")
    if expected_md5:
        if not verify_checksum(snapshot_path, expected_md5):
            print("\n  Checksum verification failed. Aborting...")
            snapshot_path.unlink()
            sys.exit(1)
    else:
        print("  Skipping checksum verification (no checksum available)")

    # Create temp directory and extract
    print(f"\n[4/5] Extracting snapshot...")
    if TEMP_EXTRACT_DIR.exists():
        run_command(f'rm -rf "{TEMP_EXTRACT_DIR}"', silent=True)
    TEMP_EXTRACT_DIR.mkdir(parents=True, exist_ok=True)

    extract_with_progress(snapshot_path, TEMP_EXTRACT_DIR)

    # Move data directory to cache location
    print(f"\n[5/5] Moving snapshot to cache...")
    temp_data_dir = TEMP_EXTRACT_DIR / "data"
    if not temp_data_dir.exists():
        print(f"  Error: Expected data directory not found at {temp_data_dir}")
        sys.exit(1)

    # Move the data directory to the cache location
    run_command(f'mv "{temp_data_dir}" "{SNAPSHOT_CACHE_DIR}"', silent=True)

    # Cleanup
    print("  Cleaning up temporary files...")
    run_command(f'rm -rf "{TEMP_EXTRACT_DIR}"', silent=True)
    snapshot_path.unlink()

    print("\n" + "=" * 60)
    print("  ✓ Snapshot cached successfully!")
    print(f"  Location: {SNAPSHOT_CACHE_DIR}")
    print("=" * 60)

if __name__ == "__main__":
    main()
