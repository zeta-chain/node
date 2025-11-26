#!/bin/bash

# Setup script for local zetaclient to connect to dockerized localnet
# This script automates the initialization and configuration of zetaclient

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration variables
ZETACLIENT_HOME="${HOME}/.zetacored"
CONFIG_FILE="${ZETACLIENT_HOME}/config/zetaclient_config.json"
RESTRICTED_ADDR_FILE="${ZETACLIENT_HOME}/config/zetaclient_restricted_addresses.json"
CHAIN_ID="athens_101-1"
OPERATOR_ADDRESS=""  # Will be fetched dynamically from observer set
DRY_MODE=1  # Set to 1 for dry mode, 0 for normal mode

# Get the script directory to find the preparams file relative to the script location
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PREPARAMS_PATH="${SCRIPT_DIR}/contrib/localnet/preparams/zetaclient-dry.json"

# Optional: Override operator address (leave empty to fetch from observer set)
OPERATOR_ADDRESS_OVERRIDE=""  # Set to a specific address if you don't want to use auto-detected

# Optional: Set to import a specific hotkey (leave empty to create new or use existing)
# Example private key from localnet for testing doc
HOTKEY_IMPORT_KEY=""  # Set to a private key hex string if you want to import a specific key

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if required binaries are installed
check_required_binaries() {
    local missing_binaries=()

    # Check for zetaclientd
    if ! command -v zetaclientd &> /dev/null; then
        missing_binaries+=("zetaclientd")
    else
        print_info "Found zetaclientd at: $(which zetaclientd)"
    fi

    # Check for zetacored
    if ! command -v zetacored &> /dev/null; then
        missing_binaries+=("zetacored")
    else
        print_info "Found zetacored at: $(which zetacored)"
    fi

    # If any binaries are missing, provide installation instructions
    if [ ${#missing_binaries[@]} -gt 0 ]; then
        print_error "Missing required binaries: ${missing_binaries[*]}"
        print_info "Please install them by running:"
        echo ""
        echo "  make install"
        echo ""
        print_info "This will build and install both zetacored and zetaclientd"
        exit 1
    fi

    print_info "All required binaries are installed"
}

# Check if Docker is running and containers are up
check_docker() {
    if ! docker info &> /dev/null; then
        print_error "Docker is not running!"
        print_info "Please start Docker Desktop and run: make start-localnet"
        exit 1
    fi

    if ! docker ps | grep -q zetacore0; then
        print_error "zetacore0 container is not running!"
        print_info "Please run: make start-localnet"
        exit 1
    fi
    print_info "Docker containers are running"
}

# Fetch operator address from observer set
fetch_operator_address() {
    # Check if override is set
    if [ -n "$OPERATOR_ADDRESS_OVERRIDE" ]; then
        OPERATOR_ADDRESS="$OPERATOR_ADDRESS_OVERRIDE"
        print_info "Using override operator address: $OPERATOR_ADDRESS"
        return
    fi

    print_info "Fetching operator address from observer set..."

    # Try to fetch observer set
    OBSERVERS=$(zetacored q observer list-observer-set --output json 2>/dev/null | jq -r '.observers[]' 2>/dev/null)

    if [ -z "$OBSERVERS" ]; then
        print_warning "Could not fetch observer set, trying alternative method..."
        # Try without json output
        OBSERVERS=$(zetacored q observer list-observer-set 2>/dev/null | grep -E "^- zeta" | sed 's/^- //')
    fi

    if [ -z "$OBSERVERS" ]; then
        print_error "Could not fetch observer set from zetacore"
        print_info "Using default operator address as fallback"
        OPERATOR_ADDRESS="zeta1tjhrwxw2ltt4j22scmht98lux22wzdyxng2pr9"
    else
        # Get the first observer address
        OPERATOR_ADDRESS=$(echo "$OBSERVERS" | head -n1)
        print_info "Found $(echo "$OBSERVERS" | wc -l | tr -d ' ') observer(s) in the set"

        # Show all observers for reference
        print_info "Available observers:"
        echo "$OBSERVERS" | while read -r obs; do
            echo "    - $obs"
        done
    fi

    print_info "Using operator address: $OPERATOR_ADDRESS"
}

# Check if pre-params file exists
check_preparams() {
    print_info "Looking for pre-params at: $PREPARAMS_PATH"
    if [ ! -f "$PREPARAMS_PATH" ]; then
        print_error "Pre-params file not found at: $PREPARAMS_PATH"
        print_info "Script directory: $SCRIPT_DIR"
        print_info "Please ensure the preparams directory exists in the repository"
        exit 1
    fi
    print_info "Found pre-params file at: $PREPARAMS_PATH"
}

# Backup existing config if it exists
backup_config() {
    if [ -f "$CONFIG_FILE" ]; then
        BACKUP_FILE="${CONFIG_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
        cp "$CONFIG_FILE" "$BACKUP_FILE"
        print_info "Backed up existing config to: $BACKUP_FILE"
    fi
}

# Initialize zetaclientd
init_zetaclient() {
    print_info "Initializing zetaclientd..."

    zetaclientd init \
        --zetacore-url 127.0.0.1 \
        --chain-id "$CHAIN_ID" \
        --operator "$OPERATOR_ADDRESS" \
        --log-format text \
        --public-ip "127.0.0.1" \
        --keyring-backend test \
        --pre-params "$PREPARAMS_PATH"

    if [ $? -eq 0 ]; then
        print_info "zetaclientd initialized successfully"
    else
        print_error "Failed to initialize zetaclientd"
        exit 1
    fi
}

# Update RPC endpoints for Docker services
update_rpc_endpoints() {
    print_info "Updating RPC endpoints to use localhost..."

    # Create temporary file for jq operations
    TMP_FILE=$(mktemp)

    # Update all endpoints
    jq '.EVMChainConfigs."1337".Endpoint = "http://127.0.0.1:8545" |
        .BTCChainConfigs."18444".RPCHost = "127.0.0.1:18443" |
        .SolanaConfig.Endpoint = "http://127.0.0.1:8899" |
        .SuiConfig.Endpoint = "http://127.0.0.1:9000" |
        .TONConfig.Endpoint = "http://127.0.0.1:8081"' \
        "$CONFIG_FILE" > "$TMP_FILE"

    if [ $? -eq 0 ]; then
        mv "$TMP_FILE" "$CONFIG_FILE"
        print_info "RPC endpoints updated successfully"
    else
        print_error "Failed to update RPC endpoints"
        rm -f "$TMP_FILE"
        exit 1
    fi
}

# Set dry mode if enabled
set_dry_mode() {
    if [ "$DRY_MODE" -eq 1 ]; then
        print_info "Setting dry mode (read-only)..."
        jq '.ClientMode = 1' "$CONFIG_FILE" > "${CONFIG_FILE}.tmp" && \
            mv "${CONFIG_FILE}.tmp" "$CONFIG_FILE"
        print_info "Dry mode enabled"
    fi
}

# Create restricted addresses config
create_restricted_addresses() {
    print_info "Creating restricted addresses config..."
    echo "[]" > "$RESTRICTED_ADDR_FILE"
    print_info "Restricted addresses config created"
}

# Add or import hotkey to keyring
setup_hotkey() {
    print_info "Setting up hotkey in keyring..."

    # Check if hotkey already exists
    if zetacored keys show hotkey --keyring-backend=test &> /dev/null; then
        print_info "Hotkey already exists in keyring"
        # Show the address
        HOTKEY_ADDRESS=$(zetacored keys show hotkey --keyring-backend=test --output json | jq -r '.address')
        print_info "Hotkey address: $HOTKEY_ADDRESS"
    else
        # If import key is provided, import it
        if [ -n "$HOTKEY_IMPORT_KEY" ]; then
            print_info "Importing hotkey from provided private key..."
            # Create a temporary file for the key
            echo "$HOTKEY_IMPORT_KEY" | zetacored keys import hotkey /dev/stdin --keyring-backend=test

            if [ $? -eq 0 ]; then
                HOTKEY_ADDRESS=$(zetacored keys show hotkey --keyring-backend=test --output json | jq -r '.address')
                print_info "Hotkey imported successfully!"
                print_info "Hotkey address: $HOTKEY_ADDRESS"
            else
                print_error "Failed to import hotkey, trying to create new one instead..."
                # Fall back to creating a new key
                zetacored keys add hotkey --algo=secp256k1 --keyring-backend=test
            fi
        else
            print_info "Creating new hotkey..."
            # Create new hotkey (this will show mnemonic)
            zetacored keys add hotkey --algo=secp256k1 --keyring-backend=test

            if [ $? -eq 0 ]; then
                HOTKEY_ADDRESS=$(zetacored keys show hotkey --keyring-backend=test --output json | jq -r '.address')
                print_info "Hotkey created successfully!"
                print_info "Hotkey address: $HOTKEY_ADDRESS"
                print_warning "IMPORTANT: Save the mnemonic phrase shown above in a secure location!"
            else
                print_error "Failed to create hotkey"
                return
            fi
        fi
    fi
}

# Import Solana relayer key (optional, not needed for dry mode)
import_relayer_key() {
    if [ "$DRY_MODE" -eq 0 ]; then
        print_info "Importing Solana relayer key..."
        zetaclientd relayer import-key \
            --network=7 \
            --private-key="3EMjCcCJg53fMEGVj13UPQpo6py9AKKyLE2qroR4yL1SvAN2tUznBvDKRYjntw7m6Jof1R2CSqjTddL27rEb6sFQ" \
            --password=pass_relayerkey
        print_info "Relayer key imported"
    else
        print_info "Skipping relayer key import (not needed in dry mode)"
    fi
}

# Test connectivity to services
test_connectivity() {
    print_info "Testing connectivity to blockchain services..."

    # Test ETH
    if curl -s http://127.0.0.1:8545 -X POST -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' | grep -q "1337"; then
        print_info "✓ ETH RPC connected (Chain ID: 1337)"
    else
        print_warning "✗ ETH RPC connection failed"
    fi

    # Test Solana
    if curl -s http://127.0.0.1:8899 -X POST -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"getHealth","id":1}' | grep -q "ok"; then
        print_info "✓ Solana RPC connected"
    else
        print_warning "✗ Solana RPC connection failed"
    fi

    # Test ZetaCore
    if curl -s http://127.0.0.1:26657/status | grep -q "athens_101-1"; then
        print_info "✓ ZetaCore RPC connected"
    else
        print_warning "✗ ZetaCore RPC connection failed"
    fi
}

# Display final configuration
display_config() {
    print_info "Final configuration:"
    echo -e "${GREEN}Configuration Summary:${NC}"
    echo "  Config Path: $CONFIG_FILE"
    echo "  Chain ID: $CHAIN_ID"
    echo "  Operator: $OPERATOR_ADDRESS"
    echo "  Mode: $([ "$DRY_MODE" -eq 1 ] && echo "Dry (read-only)" || echo "Normal")"
    echo ""
    echo -e "${GREEN}RPC Endpoints:${NC}"
    jq -r '
        "  ZetaCore: " + .ZetaCoreURL + "\n" +
        "  EVM: " + .EVMChainConfigs."1337".Endpoint + "\n" +
        "  Bitcoin: " + .BTCChainConfigs."18444".RPCHost + "\n" +
        "  Solana: " + .SolanaConfig.Endpoint + "\n" +
        "  Sui: " + .SuiConfig.Endpoint + "\n" +
        "  TON: " + .TONConfig.Endpoint
    ' "$CONFIG_FILE" 2>/dev/null || print_warning "Could not display all endpoints"
}

# Main execution
main() {
    echo -e "${GREEN}============================================${NC}"
    echo -e "${GREEN}  ZetaClient Local Setup Script${NC}"
    echo -e "${GREEN}============================================${NC}"
    echo ""

    # Run all setup steps
    check_required_binaries
    check_docker
    fetch_operator_address
    check_preparams
    backup_config

    # Ask user if they want to reinitialize if config already exists
    if [ -f "$CONFIG_FILE" ]; then
        read -p "Config already exists. Reinitialize? (y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Skipping initialization, updating existing config..."
        else
            init_zetaclient
        fi
    else
        init_zetaclient
    fi

    update_rpc_endpoints
    set_dry_mode
    create_restricted_addresses
    setup_hotkey
    import_relayer_key
    test_connectivity
    display_config

    echo ""
    echo -e "${GREEN}============================================${NC}"
    echo -e "${GREEN}  Setup Complete!${NC}"
    echo -e "${GREEN}============================================${NC}"
    echo ""
    echo -e "${GREEN}To start zetaclientd:${NC}"
    echo "  zetaclientd start"
    echo ""
    echo -e "${YELLOW}============================================${NC}"
    echo -e "${YELLOW}  IMPORTANT: Initialize Hotkey on ZetaCore${NC}"
    echo -e "${YELLOW}============================================${NC}"
    echo ""
    echo "The hotkey created by this script needs to be initialized on ZetaCore."
    echo "The simplest way is to transfer some tokens to the hotkey address."
    echo ""
    echo -e "${GREEN}Step 1:${NC} Get your hotkey address"
    echo "  zetacored keys show hotkey --keyring-backend=test"
    echo ""
    echo -e "${GREEN}Step 2:${NC} Connect to zetacore0 container"
    echo "  docker exec -it zetacore0 /bin/bash"
    echo ""
    echo -e "${GREEN}Step 3:${NC} Inside the container, get the operator address"
    echo "  zetacored keys show operator --keyring-backend=test"
    echo ""
    echo -e "${GREEN}Step 4:${NC} Send tokens to your hotkey"
    echo "  zetacored tx bank send <OPERATOR_ADDRESS> <HOTKEY_ADDRESS> \\"
    echo "    10000000000000000000000azeta \\"
    echo "    --keyring-backend=test \\"
    echo "    --fees=2000000000000000azeta \\"
    echo "    --chain-id=athens_101-1 \\"
    echo "    --yes"
    echo ""
    echo -e "${YELLOW}Note:${NC} Replace <OPERATOR_ADDRESS> and <HOTKEY_ADDRESS> with actual addresses"
    echo -e "${YELLOW}============================================${NC}"
    echo ""
}

# Run main function
main "$@"