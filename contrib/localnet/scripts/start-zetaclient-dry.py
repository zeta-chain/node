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

# Configuration from environment variables
ZETACORE_HOST = os.environ.get('ZETACORE_HOST')

# Chain ID will be fetched from zetacore node
CHAIN_ID = None

# Validate required environment variables
def validate_env_vars():
    """Validate that required environment variables are set."""
    if not ZETACORE_HOST:
        print("Error: Required environment variable not set: ZETACORE_HOST")
        print("Please set the following environment variable:")
        print("  ZETACORE_HOST - Hostname of the zetacore node (e.g., zetacore0, testnet-node)")
        sys.exit(1)


# Static configuration
ZETACLIENT_HOME = "/root/.zetacored"
CONFIG_FILE = f"{ZETACLIENT_HOME}/config/zetaclient_config.json"
RESTRICTED_ADDR_FILE = f"{ZETACLIENT_HOME}/config/zetaclient_restricted_addresses.json"
PREPARAMS_PATH = "/root/preparams/zetaclient-dry.json"
CLIENT_MODE = 1  # Dry mode (read-only)

# Required ports for zetaclient to connect to zetacore
REQUIRED_PORTS = {
    1317: "Cosmos REST API",
    9090: "Cosmos gRPC",
    26657: "CometBFT RPC",
}

# Chain-specific configurations for external chain RPC endpoints
# Dry mode doesn't need external chains, so testnet/mainnet configs are empty
# Localnet uses Docker network hostnames to reach other containers
CHAIN_CONFIGS = {
    "athens_101-1": {  # localnet - use Docker network hostnames
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
    "athens_7001-1": {  # testnet
        "EVMChainConfigs": {},
        "BTCChainConfigs": {},
        "SolanaConfig": {},
        "SuiConfig": {},
        "TONConfig": {},
    },
    "zetachain_7000-1": {  # mainnet
        "EVMChainConfigs": {},
        "BTCChainConfigs": {},
        "SolanaConfig": {},
        "SuiConfig": {},
        "TONConfig": {},
    },
}


def get_chain_config():
    """Get chain-specific configuration or error if chain ID is unknown."""
    if CHAIN_ID not in CHAIN_CONFIGS:
        print(f"Error: Unknown chain ID: {CHAIN_ID}")
        print(f"Supported chain IDs: {', '.join(CHAIN_CONFIGS.keys())}")
        sys.exit(1)
    return CHAIN_CONFIGS[CHAIN_ID]


def check_port(host, port, timeout=2):
    """Check if a port is open on the given host."""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        result = sock.connect_ex((host, port))
        sock.close()
        return result == 0
    except socket.error:
        return False


def check_required_ports():
    """Check that all required ports are accessible on testnet-node."""
    print("Checking required ports on zetacore")

    all_ports_open = True
    for port, description in REQUIRED_PORTS.items():
        if check_port(ZETACORE_HOST, port):
            print(f"  ✓ Port {port} ({description}) is open")
        else:
            print(f"  ✗ Port {port} ({description}) is NOT reachable")
            all_ports_open = False

    if not all_ports_open:
        sys.exit(1)

    print("All required ports are accessible.")


def run_command(cmd, capture_output=True, check=True, input_data=None):
    """Run a shell command and return the result."""
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
    """Wait for zetacore node to be ready and fetch chain ID."""
    global CHAIN_ID
    print(f"Waiting for {ZETACORE_HOST} to be ready...")

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
                        return
                except json.JSONDecodeError:
                    pass
        except Exception:
            pass

        print(f"Waiting for {ZETACORE_HOST}...")
        time.sleep(5)


def fetch_operator_address():
    """Fetch operator address from observer set."""
    print(f"Fetching operator address from {ZETACORE_HOST}...")

    # Try to fetch observer set from zetacore node
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

    # If no observers, use a default address for testnet
    print("No observers found, using default testnet operator address")
    sys.exit(1)


def init_zetaclient(operator_address):
    """Initialize zetaclient if config doesn't exist."""
    if not os.path.exists(CONFIG_FILE):
        print("Initializing zetaclient for testnet...")
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
        print("Config already exists, updating...")


def update_config(operator_address):
    """Update config for zetacore connection."""
    print("Updating zetaclient config...")

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
    config['SolanaConfig'] = chain_config['SolanaConfig']
    config['SuiConfig'] = chain_config['SuiConfig']
    config['TONConfig'] = chain_config['TONConfig']

    # Write updated config
    with open(CONFIG_FILE, 'w') as f:
        json.dump(config, f, indent=2)


def setup_hotkey():
    """Setup hotkey from mnemonic."""
    print("Setting up hotkey...")

    # Clean up the entire keyring directory to avoid corruption issues
    keyring_path = "/root/.zetacored/keyring-test"
    if os.path.exists(keyring_path):
        import shutil
        shutil.rmtree(keyring_path)

    print("Creating hotkey from mnemonic...")
    mnemonic = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
    run_command(
        f'echo "{mnemonic}" | zetacored keys add hotkey --algo=secp256k1 --recover --keyring-backend=test --output json > /dev/null 2>&1',
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
    """Create restricted addresses config."""
    print("Creating restricted addresses config...")
    with open(RESTRICTED_ADDR_FILE, 'w') as f:
        json.dump([], f)


def wait_for_tss():
    """Wait for TSS to be initialized before starting dry mode client."""
    print("Waiting for TSS to be initialized...")

    while True:
        result = run_command(
            f"zetacored q observer get-tss-address --node tcp://{ZETACORE_HOST}:26657 --output json",
            check=False
        )
        if result:
            try:
                data = json.loads(result)
                # Check if we got a valid TSS address
                tss_address = data.get('eth') or data.get('btc')
                if tss_address:
                    print(f"TSS is ready (eth: {data.get('eth', 'N/A')})")
                    return
            except json.JSONDecodeError:
                pass

        print("Waiting for TSS...")
        time.sleep(5)


def start_zetaclientd():
    """Start zetaclientd with proper parameters."""
    print("Starting zetaclientd in dry mode for testnet...")

    # Create password file with empty values
    # Hotkey uses test keyring (no password), TSS and Solana not used in dry mode
    passwords = "\n\n\n"

    with open("/root/password.file", 'w') as f:
        f.write(passwords)

    # Start zetaclientd with passwords piped from file
    print("Starting zetaclientd...")

    # Use exec to replace the current process
    os.execlp(
        "bash", "bash", "-c",
        "exec zetaclientd start < /root/password.file"
    )


def main():
    """Main execution."""
    print("========================================")
    print("  Starting ZetaClient Dry Mode         ")
    print("========================================")

    # Validate required environment variables first
    validate_env_vars()

    print(f"Zetacore Host: {ZETACORE_HOST}")

    # Wait for zetacore and fetch chain ID from node status
    wait_for_zetacore()
    check_required_ports()
    operator_address = fetch_operator_address()
    init_zetaclient(operator_address)
    setup_hotkey()
    update_config(operator_address)
    create_restricted_addresses()

    # Wait for TSS to be initialized before starting dry mode client
    wait_for_tss()

    print("Configuration complete, starting zetaclientd...")
    start_zetaclientd()


if __name__ == "__main__":
    main()
