#!/bin/bash

KEYRING_TEST="test"
KEYRING_FILE="file"
HOSTNAME=$(hostname)

# Operator key
zetacored keys add operator --algo=secp256k1 --keyring-backend=$KEYRING_TEST
operator_address=$(zetacored keys show operator -a --keyring-backend=$KEYRING_TEST)

# Hotkey key depending on the keyring-backend
if [ "$HOTKEY_BACKEND" = "$KEYRING_FILE" ]; then
    printf "%s\n%s\n" "$HOTKEY_PASSWORD" "$HOTKEY_PASSWORD" | zetacored keys add hotkey --algo=secp256k1 --keyring-backend=$KEYRING_FILE
    hotkey_address=$(printf "%s\n%s\n" "$HOTKEY_PASSWORD" "$HOTKEY_PASSWORD" | zetacored keys show hotkey -a --keyring-backend=$KEYRING_FILE)

    # Get hotkey pubkey, the command use keyring-backend in the cosmos config
    zetacored config keyring-backend "$KEYRING_FILE"
    pubkey=$(printf "%s\n%s\n" "$HOTKEY_PASSWORD" "$HOTKEY_PASSWORD" | zetacored get-pubkey hotkey | sed -e 's/secp256k1:"\(.*\)"/\1/' |sed 's/ //g' )
    zetacored config keyring-backend "$KEYRING_TEST"
else
    zetacored keys add hotkey --algo=secp256k1 --keyring-backend=$KEYRING
    hotkey_address=$(zetacored keys show hotkey -a --keyring-backend=$KEYRING)
    pubkey=$(zetacored get-pubkey hotkey|sed -e 's/secp256k1:"\(.*\)"/\1/' | sed 's/ //g' )
fi

is_observer="y"

echo "operator_address: $operator_address"
echo "hotkey_address: $hotkey_address"
echo "pubkey: $pubkey"
mkdir ~/.zetacored/os_info

# set key in file
jq -n --arg is_observer "$is_observer" --arg operator_address "$operator_address" --arg hotkey_address "$hotkey_address" --arg pubkey "$pubkey" '{"IsObserver":$is_observer,"ObserverAddress":$operator_address,"ZetaClientGranteeAddress":$hotkey_address,"ZetaClientGranteePubKey":$pubkey}' > ~/.zetacored/os_info/os.json
