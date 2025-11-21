#!/bin/sh

set -e

# Create log directory
mkdir -p /var/log/zetaclient

# Create necessary directories (from start-zetaclientd.sh lines 67-69, 71)
mkdir -p /home/zetachain/.zetacored/relayer-keys
mkdir -p /home/zetachain/.zetacored/config
mkdir -p /home/zetachain/.zetacored/os_info
mkdir -p /home/zetachain/.tss

# Determine keyring backend (from start-zetaclientd.sh lines 50-54)
BACKEND="test"
if [ "$HOTKEY_BACKEND" == "file" ]; then
    BACKEND="file"
fi

HOTKEY_NAME=${HOTKEY_NAME:-"hotkey"}

# Create hotkey if it doesn't exist (adapted from add-keys.sh)
if ! zetacored keys show "$HOTKEY_NAME" --keyring-backend="$BACKEND" --home /home/zetachain/.zetacored >/dev/null 2>&1; then
    if [ "$BACKEND" == "file" ]; then
        # File backend requires password
        printf "%s\n%s\n" "${HOTKEY_PASSWORD:-}" "${HOTKEY_PASSWORD:-}" | \
            zetacored keys add "$HOTKEY_NAME" --algo=secp256k1 --keyring-backend="$BACKEND" --home /home/zetachain/.zetacored
    else
        # Test backend doesn't need password
        printf '\n' | zetacored keys add "$HOTKEY_NAME" --algo=secp256k1 --keyring-backend="$BACKEND" --home /home/zetachain/.zetacored
    fi
fi

# Initialize restricted addresses config (from start-zetaclientd.sh line 127)
echo "[]" > /home/zetachain/.zetacored/config/zetaclient_restricted_addresses.json

# Start zetaclientd (pipe empty passwords for observer mode)
echo -e '\n\n\n' | exec zetaclientd start --home /home/zetachain/.zetacored
