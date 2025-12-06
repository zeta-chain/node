#!/usr/bin/env python3
"""
Start script for ZetaClient in dry mode for testnet.
This script initializes and configures zetaclientd for read-only observation of a testnet.
"""

import json
import os
import socket
import subprocess
import sys
import time
import shutil
from time import sleep

# Configuration from environment variables
ZETACORE_HOST = os.environ.get('ZETACORE_HOST')
RPC_API_KEY_ALLTHATNODE = os.environ.get('RPC_API_KEY_ALLTHATNODE')

# Chain ID will be fetched from zetacore node
CHAIN_ID = None

# Validate required environment variables
def validate_env_vars():
    if not ZETACORE_HOST:
        print("Error: Required environment variable ZETACORE_HOST not set")
        print("  ZETACORE_HOST - Hostname of the zetacore node (e.g., zetacore0, testnet-node)")
        sys.exit(1)


def validate_rpc_api_key():
    """Validate RPC API key is set for non-localnet chains."""
    if CHAIN_ID != "athens_101-1" and not RPC_API_KEY_ALLTHATNODE:
        print("Error: RPC_API_KEY_ALLTHATNODE is required for non-localnet chains")
        print(f"  Chain ID: {CHAIN_ID}")
        print("  Please set RPC_API_KEY_ALLTHATNODE environment variable")
        sys.exit(1)


# Static configuration
ZETACLIENT_HOME = "/root/.zetacored"
CONFIG_FILE = f"{ZETACLIENT_HOME}/config/zetaclient_config.json"
RESTRICTED_ADDR_FILE = f"{ZETACLIENT_HOME}/config/zetaclient_restricted_addresses.json"
PREPARAMS_PATH = "/root/static-preparams/zetaclient-dry.json"
CLIENT_MODE = 1  # Dry mode

# Hotkey mnemonic for dry mode client
# Address: zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax (0x8E3C1898776e80A19a37546920AcE1935cCEE08E)
HOTKEY_MNEMONIC = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"

# Required ports for zetaclient to connect to zetacore
REQUIRED_PORTS = {
    1317: "Cosmos REST API",
    9090: "Cosmos gRPC",
    26657: "CometBFT RPC",
}

# Chain-specific configurations for external chain RPC endpoints
CHAIN_CONFIGS = {
    "athens_101-1": {
        "EVMChainConfigs": {
            "1337": {"Endpoint": "http://eth:8545"},
        },
        "BTCChainConfigs": {
            "18444": {
                "RPCUsername": "smoketest",
                "RPCPassword": "123",
                "RPCHost": "bitcoin:18443",
                "RPCParams": "regtest",
            },
        },
        "SolanaConfig": {"Endpoint": "http://solana:8899"},
        "SuiConfig": {"Endpoint": "http://sui:9000"},
        "TONConfig": {"Endpoint": "http://ton:8081"},
    },
    "athens_7001-1": {
        "EVMChainConfigs": {
            "11155111": {"Endpoint": f"https://ethereum-sepolia.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},  # Ethereum Sepolia
            "97": {"Endpoint": f"https://bsc-testnet.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},              # BSC Testnet
            "80002": {"Endpoint": f"https://polygon-amoy.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},          # Polygon Amoy
            "84532": {"Endpoint": f"https://base-sepolia.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},          # Base Sepolia
            "421614": {"Endpoint": f"https://arbitrum-sepolia.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},     # Arbitrum Sepolia
            "43113": {"Endpoint": f"https://avalanche-fuji.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}/ext/bc/C/rpc"},  # Avalanche Fuji
        },
        "BTCChainConfigs": {},
        "SolanaConfig": {"Endpoint": f"https://solana-devnet.g.allthatnode.com/archive/json_rpc/{RPC_API_KEY_ALLTHATNODE}"},
    },
    "zetachain_7000-1": {  # mainnet
        "EVMChainConfigs": {
            "1": {"Endpoint": f"https://ethereum-mainnet.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},    # Ethereum
            "56": {"Endpoint": f"https://bsc-mainnet.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},        # BSC
            "137": {"Endpoint": f"https://polygon-mainnet.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},   # Polygon
            "8453": {"Endpoint": f"https://base-mainnet.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},     # Base
            "42161": {"Endpoint": f"https://arbitrum-one.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}"},    # Arbitrum
            "43114": {"Endpoint": f"https://avalanche-mainnet.g.allthatnode.com/full/evm/{RPC_API_KEY_ALLTHATNODE}/ext/bc/C/rpc"},  # Avalanche
        },
        "BTCChainConfigs": {},
        "SolanaConfig": {"Endpoint": f"https://solana-mainnet.g.allthatnode.com/full/json_rpc/{RPC_API_KEY_ALLTHATNODE}"},
    },
}


def get_chain_config():
    if CHAIN_ID not in CHAIN_CONFIGS:
        print(f"Error: Unknown chain ID: {CHAIN_ID}")
        print(f"Supported chain IDs: {', '.join(CHAIN_CONFIGS.keys())}")
        sys.exit(1)
    return CHAIN_CONFIGS[CHAIN_ID]


def check_port(host, port, timeout=2):
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        result = sock.connect_ex((host, port))
        sock.close()
        return result == 0
    except socket.error:
        return False


def check_required_ports():
    all_ports_open = True
    for port, description in REQUIRED_PORTS.items():
        if check_port(ZETACORE_HOST, port):
            print(f"  ✓ Port {port} ({description}) is open")
        else:
            print(f"  ✗ Port {port} ({description}) is NOT reachable")
            all_ports_open = False

    if not all_ports_open:
        sys.exit(1)


def run_command(cmd, capture_output=True, check=True, input_data=None):
    try:
        result = subprocess.run(
            cmd,
            shell=True,
            capture_output=capture_output,
            text=True,
            check=check,
            input=input_data
        )
        return result.stdout.strip() if capture_output else None
    except subprocess.CalledProcessError as e:
        if check:
            raise
        return None


def wait_for_zetacore():
    global CHAIN_ID
    print(f"Waiting for {ZETACORE_HOST} to be ready")

    while True:
        try:
            result = subprocess.run(
                f"curl -s http://{ZETACORE_HOST}:26657/status",
                shell=True,
                capture_output=True,
                text=True
            )
            if result.returncode == 0 and result.stdout:
                # Parse the status response to get chain ID
                try:
                    status = json.loads(result.stdout)
                    network = status.get('result', {}).get('node_info', {}).get('network')
                    if network:
                        CHAIN_ID = network
                        print(f"{ZETACORE_HOST} is ready")
                        print(f"Detected chain ID: {CHAIN_ID}")
                        sleep(6) # wait a bit for few block to be produced
                        return
                except json.JSONDecodeError:
                    pass
        except Exception:
            pass

        print(f"Waiting for {ZETACORE_HOST}...")
        time.sleep(5)


def fetch_operator_address():
    try:
        result = run_command(
            f"zetacored q observer list-observer-set --node tcp://{ZETACORE_HOST}:26657 --output json",
            check=False
        )
        if result:
            data = json.loads(result)
            observers = data.get('observers', [])
            if observers:
                operator_address = observers[0]
                print(f"Found operator address: {operator_address}")
                return operator_address
    except (json.JSONDecodeError, Exception):
        pass

    sys.exit(1)


def init_zetaclient(operator_address):
    if not os.path.exists(CONFIG_FILE):
        run_command(
            f'zetaclientd init '
            f'--zetacore-url "{ZETACORE_HOST}" '
            f'--chain-id "{CHAIN_ID}" '
            f'--operator "{operator_address}" '
            f'--log-format text '
            f'--public-ip "127.0.0.1" '
            f'--keyring-backend test '
            f'--pre-params "{PREPARAMS_PATH}"'
        )
    else:
        print("Config already exists, updating")


def update_config(operator_address):
    print("Updating zetaclient config")

    # Get chain-specific configuration
    chain_config = get_chain_config()

    # Read existing config
    with open(CONFIG_FILE, 'r') as f:
        config = json.load(f)

    # Update config values
    config['ChainID'] = CHAIN_ID
    config['ZetaCoreURL'] = ZETACORE_HOST
    config['AuthzGranter'] = operator_address
    config['ClientMode'] = CLIENT_MODE
    config['AuthzHotkey'] = 'hotkey'

    # Apply chain-specific external chain configs
    config['EVMChainConfigs'] = chain_config['EVMChainConfigs']
    config['BTCChainConfigs'] = chain_config['BTCChainConfigs']
    config['SolanaConfig'] = chain_config.get('SolanaConfig', {})
    config['SuiConfig'] = chain_config.get('SuiConfig', {})
    config['TONConfig'] = chain_config.get('TONConfig', {})

    # Write updated config
    with open(CONFIG_FILE, 'w') as f:
        json.dump(config, f, indent=2)


def setup_hotkey():
    """Setup hotkey from mnemonic."""
    # Clean up the entire keyring directory to avoid corruption issues
    keyring_path = "/root/.zetacored/keyring-test"
    if os.path.exists(keyring_path):
        shutil.rmtree(keyring_path)

    run_command(
        f'echo "{HOTKEY_MNEMONIC}" | zetacored keys add hotkey --algo=secp256k1 --recover --keyring-backend=test --output json > /dev/null 2>&1',
        capture_output=False,
        check=False
    )

    # Get the hotkey address
    result = run_command(
        "zetacored keys show hotkey --keyring-backend=test --output json",
        check=False
    )
    if not result:
        raise RuntimeError("Failed to create hotkey from mnemonic")

    data = json.loads(result)
    hotkey_address = data.get('address')
    if not hotkey_address:
        raise RuntimeError("Failed to get hotkey address")

    print(f"Hotkey address: {hotkey_address}")


def create_restricted_addresses():
    with open(RESTRICTED_ADDR_FILE, 'w') as f:
        json.dump([], f)


def wait_for_tss():
    while True:
        result = run_command(
            f"zetacored q observer get-tss-address --node tcp://{ZETACORE_HOST}:26657 --output json",
            check=False
        )
        if result:
            try:
                data = json.loads(result)
                tss_address = data.get('eth') or data.get('btc')
                if tss_address:
                    print(f"TSS is ready (eth: {data.get('eth', 'N/A')})")
                    return
            except json.JSONDecodeError:
                pass

        print("Waiting for TSS")
        time.sleep(5)


def start_zetaclientd():
    with open("/root/password.file", 'w') as f:
        f.write("\n\n\n")
    os.execlp(
        "bash", "bash", "-c",
        "exec zetaclientd start < /root/password.file"
    )


def main():
    print("========================================")
    print("  Starting ZetaClient Dry Mode         ")
    print("========================================")

    validate_env_vars()

    print(f"Zetacore Host: {ZETACORE_HOST}")
    wait_for_zetacore()
    validate_rpc_api_key()
    check_required_ports()
    operator_address = fetch_operator_address()
    init_zetaclient(operator_address)
    setup_hotkey()
    update_config(operator_address)
    create_restricted_addresses()
    wait_for_tss()

    start_zetaclientd()


if __name__ == "__main__":
    main()
