#!/bin/bash

# This script allows to add keys for operator and hotkey and create the required json structure for os_info

KEYRING_TEST="test"
KEYRING_FILE="file"
HOSTNAME=$(hostname)

# Check if is_observer flag is provided
if [ -z "$1" ]; then
    is_observer="y" # Default value if not provided
else
    is_observer="$1"
fi


zetacored keys add operator --algo=secp256k1 --keyring-backend=$KEYRING_TEST

operator_address=$(zetacored keys show operator --keyring-backend=$KEYRING_TEST | sed -n 's/.*address: \(zeta[a-zA-Z0-9]*\).*/\1/p')

# Hotkey key depending on the keyring-backend
if [ "$HOTKEY_BACKEND" == "$KEYRING_FILE" ]; then
    printf "%s\n%s\n" "$HOTKEY_PASSWORD" "$HOTKEY_PASSWORD" | zetacored keys add hotkey --algo=secp256k1 --keyring-backend=$KEYRING_FILE
    hotkey_address=$(printf "%s\n%s\n" "$HOTKEY_PASSWORD" "$HOTKEY_PASSWORD" | zetacored keys show hotkey --keyring-backend=$KEYRING_FILE | sed -n 's/.*address: \(zeta[a-zA-Z0-9]*\).*/\1/p')

    # TODO: remove after v50 upgrade
    # Get hotkey pubkey, the command use keyring-backend in the cosmos config
    if ! zetacored config set client keyring-backend "$KEYRING_FILE"; then
        zetacored config keyring-backend "$KEYRING_FILE"
    fi
    pubkey=$(printf "%s\n%s\n" "$HOTKEY_PASSWORD" "$HOTKEY_PASSWORD" | zetacored get-pubkey hotkey | sed -n 's/secp256k1:"\(zetapub[a-zA-Z0-9]*\)".*/\1/p')
    if ! zetacored config set client keyring-backend "$KEYRING_TEST"; then
        zetacored config keyring-backend "$KEYRING_TEST"
    fi
else
    zetacored keys add hotkey --algo=secp256k1 --keyring-backend=$KEYRING_TEST
    hotkey_address=$(zetacored keys show hotkey --keyring-backend=$KEYRING_TEST | sed -n 's/.*address: \(zeta[a-zA-Z0-9]*\).*/\1/p')
    pubkey=$(zetacored get-pubkey hotkey | sed -n 's/secp256k1:"\(zetapub[a-zA-Z0-9]*\)".*/\1/p')
fi

echo "operator_address: $operator_address"
echo "hotkey_address: $hotkey_address"
echo "pubkey: $pubkey"
echo "is_observer: $is_observer"
mkdir -p ~/.zetacored/os_info

# set key in file
jq -n --arg is_observer "$is_observer" --arg operator_address "$operator_address" --arg hotkey_address "$hotkey_address" --arg pubkey "$pubkey" '{"IsObserver":$is_observer,"ObserverAddress":$operator_address,"ZetaClientGranteeAddress":$hotkey_address,"ZetaClientGranteePubKey":$pubkey}' > ~/.zetacored/os_info/os.json