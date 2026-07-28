#!/usr/bin/env python3
"""
Snapshot Export Script
======================
Produces a `zetacored export` genesis JSON for testnet or mainnet, exporting the
modules needed for a downstream snapshot (e.g. a ZETA -> Solana migration snapshot).

Unlike the devnet fork flow, this does NOT sync or start the node: `zetacored export`
reads committed state directly from the data dir, so we only init a home, drop in the
network config, load a snapshot's `data/` into it, and run the export.

Run from the root zeta-node directory:
    python3 contrib/snapshot/export_snapshot.py --network testnet
    python3 contrib/snapshot/export_snapshot.py --network mainnet
Or:
    make snapshot-export NETWORK=mainnet
    make snapshot-export NETWORK=testnet SNAPSHOT_EXPORT_ARGS="--height 1000000 --modules bank,staking"

Examples:
    # Testnet, latest height, default modules (bank,staking,distribution)
    python3 contrib/snapshot/export_snapshot.py --network testnet

    # Mainnet, stamping a version label on the locally-built binary
    python3 contrib/snapshot/export_snapshot.py --network mainnet --node-version v36.0.1

    # Reuse an already-downloaded snapshot cache
    python3 contrib/snapshot/export_snapshot.py --network mainnet --skip-download

    # Dry run: print every resolved command, execute nothing
    python3 contrib/snapshot/export_snapshot.py --network mainnet --dry-run

Caveats:
    1. Binary-version compatibility: the `zetacored` binary must be able to open the
       snapshot's store, so build it from source checked out at the release the snapshot
       was produced with (e.g. `git checkout v36.0.1`) BEFORE running. `--node-version`
       only stamps the version string into the binary via `make install`'s linker flags;
       it does NOT fetch or check out that release. Opening a snapshot with an
       incompatible binary can fail or read incorrectly.
    2. Pruned height: `--height <H>` for a specific old height only works if that height
       is still un-pruned in the snapshot store. The default (-1 / latest) always works.
"""

import argparse
import subprocess
import sys
from pathlib import Path

import requests

# ============================================================================
# Configuration Constants
# ============================================================================

HOME_DIR = Path.home()

# Default export home — deliberately NOT the user's ~/.zetacored.
DEFAULT_EXPORT_HOME = HOME_DIR / ".zetacored-snapshot-export"

# network-config raw base + the files we pull into <HOME>/config/.
NETWORK_CONFIG_BASE = "https://raw.githubusercontent.com/zeta-chain/network-config/main"
CONFIG_FILES = ["genesis.json", "app.toml", "config.toml", "client.toml"]

# Snapshot downloader is reused via subprocess rather than reimplemented.
DOWNLOAD_SNAPSHOT_SCRIPT = (
    Path(__file__).resolve().parent.parent
    / "localnet"
    / "scripts_python"
    / "download_snapshot.py"
)

NETWORKS = {
    "testnet": {
        "chain_id": "athens_7001-1",
        "config_dir": "athens3",
        "cache_dir": HOME_DIR / "zetacored_snapshot_testnet",
        "name": "Testnet",
    },
    "mainnet": {
        "chain_id": "zetachain_7000-1",
        "config_dir": "mainnet",
        "cache_dir": HOME_DIR / "zetacored_snapshot_mainnet",
        "name": "Mainnet",
    },
}

# ============================================================================
# Helper Functions
# ============================================================================

def run_command(cmd, dry_run=False, shell=True, check=True):
    """Run a shell command, or just print it in dry-run mode."""
    print(f"Running: {cmd}")
    if dry_run:
        return
    try:
        subprocess.run(cmd, shell=shell, check=check, text=True)
    except subprocess.CalledProcessError as e:
        print(f"Error running command: {cmd}")
        print(f"Error: {e}")
        sys.exit(1)


def download_file(url, dest_path, dry_run=False):
    """Download a file from URL to destination path, or print it in dry-run mode."""
    print(f"Downloading {url} to {dest_path}")
    if dry_run:
        return
    try:
        response = requests.get(url, timeout=300)
        response.raise_for_status()
        with open(dest_path, 'wb') as f:
            f.write(response.content)
        print(f"Successfully downloaded {dest_path}")
    except Exception as e:
        print(f"Error downloading {url}: {e}")
        sys.exit(1)

def strip_preamble(raw_path, out_path):
    """Copy raw_path to out_path, dropping any lines before the JSON object.

    `zetacored export` prints the genesis JSON to stdout but the node logger can
    emit warnings there first, so keep only from the first line starting with '{'.
    """
    with open(raw_path) as src, open(out_path, 'w') as dst:
        started = False
        for line in src:
            if not started and not line.lstrip().startswith('{'):
                continue
            started = True
            dst.write(line)
    if not started:
        print(f"Error: no JSON object found in export output {raw_path}")
        sys.exit(1)

# ============================================================================
# Main Script
# ============================================================================

def main(args):
    network = NETWORKS[args.network]
    chain_id = network["chain_id"]
    cache_dir = network["cache_dir"]
    home = Path(args.home).expanduser() if args.home else DEFAULT_EXPORT_HOME
    config_dir = home / "config"

    height_label = "latest" if args.height == -1 else str(args.height)
    output = Path(args.output).expanduser() if args.output else (
        HOME_DIR / f"zeta-snapshot-{args.network}-{height_label}.json"
    )
    log_path = output.with_suffix(".log")

    print("=" * 80)
    print(f"ZetaChain Snapshot Export — {network['name']}")
    print("=" * 80)
    print(f"Chain ID:  {chain_id}")
    print(f"Home:      {home}")
    print(f"Modules:   {args.modules}")
    print(f"Height:    {height_label}")
    print(f"Output:    {output}")
    if args.dry_run:
        print("Mode:      DRY RUN (nothing will be executed)")
    print("=" * 80)

    # Step 1: Optionally build the binary at a specific version.
    print("\n[1/5] Building zetacored...")
    if args.node_version:
        run_command(f"NODE_VERSION={args.node_version} make install", dry_run=args.dry_run)
    else:
        print("No --node-version given; using zetacored already on PATH.")

    # Step 2: Initialize the export home.
    print("\n[2/5] Initializing export home...")
    run_command(
        f"zetacored init export-node --chain-id {chain_id} --home {home} -o",
        dry_run=args.dry_run,
    )

    # Step 3: Download network-config files into <HOME>/config/.
    print("\n[3/5] Downloading network configuration...")
    base = f"{NETWORK_CONFIG_BASE}/{network['config_dir']}"
    for filename in CONFIG_FILES:
        download_file(f"{base}/{filename}", config_dir / filename, dry_run=args.dry_run)

    # Step 4: Load the snapshot's data dir into the export home.
    print("\n[4/5] Loading snapshot data...")
    if args.skip_download:
        print("--skip-download set; reusing existing snapshot cache.")
    else:
        run_command(
            f"python3 {DOWNLOAD_SNAPSHOT_SCRIPT} --chain-id {chain_id}",
            dry_run=args.dry_run,
        )

    data_dir = home / "data"
    run_command(f'rm -rf "{data_dir}"', dry_run=args.dry_run)
    run_command(f'cp -r "{cache_dir / "data"}" "{data_dir}"', dry_run=args.dry_run)

    # Step 5: Run the export.
    print("\n[5/5] Exporting genesis...")
    raw_output = output.with_name(output.name + ".raw")
    export_cmd = f"zetacored export --home {home} --modules-to-export {args.modules}"
    if args.height != -1:
        export_cmd += f" --height {args.height}"
    if args.for_zero_height:
        export_cmd += " --for-zero-height"
    export_cmd += f' > "{raw_output}" 2> "{log_path}"'
    run_command(export_cmd, dry_run=args.dry_run)

    # The node logger emits an early warning to stdout ahead of the JSON, which
    # would make the output invalid JSON. Keep only from the first line that
    # opens the genesis object.
    if not args.dry_run:
        strip_preamble(raw_output, output)
        raw_output.unlink()

    print("\n" + "=" * 80)
    print(f"Export written to: {output}")
    print(f"Export log:        {log_path}")
    print("=" * 80)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Export a zetacored genesis snapshot for testnet or mainnet",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "--network",
        required=True,
        choices=["testnet", "mainnet"],
        help="Which network to export",
    )
    parser.add_argument(
        "--node-version",
        default=None,
        help="Build this node version before exporting (e.g. v36.0.1). "
             "If omitted, uses the zetacored already on PATH.",
    )
    parser.add_argument(
        "--modules",
        default="bank,staking,distribution",
        help="Comma-separated modules to export (default: bank,staking,distribution)",
    )
    parser.add_argument(
        "--height",
        type=int,
        default=-1,
        help="Height to export (default: -1 = latest)",
    )
    parser.add_argument(
        "--for-zero-height",
        action="store_true",
        help="Prepare the export for a fresh zero-height genesis",
    )
    parser.add_argument(
        "--output",
        default=None,
        help="Output path for the exported JSON "
             "(default: ~/zeta-snapshot-<network>-<height>.json)",
    )
    parser.add_argument(
        "--home",
        default=None,
        help=f"Export home directory (default: {DEFAULT_EXPORT_HOME})",
    )
    parser.add_argument(
        "--skip-download",
        action="store_true",
        help="Reuse the existing snapshot cache instead of downloading",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Print every resolved command but execute nothing",
    )
    args = parser.parse_args()

    main(args)
