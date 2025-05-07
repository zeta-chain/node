#!/bin/bash

service ssh start
echo "making an id"
solana-keygen new -o /root/.config/solana/id.json --no-bip39-passphrase

solana config set --url localhost
echo "starting solana test validator..."
solana-test-validator --limit-ledger-size 50000000 &

sleep 5
# airdrop to e2e sol and spl accounts
solana airdrop 1000
solana airdrop 1000 37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ
solana airdrop 1000 BZRrLRu7VktRkZt7ZihxP9PLXjBf8vPdVb9dQU4Bj6my

# Deploy initial programs
solana program deploy gateway.so
solana program deploy connected.so
solana program deploy connected_spl.so

# Get program ID from gateway keypair
GATEWAY_PROGRAM_ID=$(solana-keygen pubkey gateway-keypair.json)


echo "Gateway program ID: $GATEWAY_PROGRAM_ID"
echo "Starting upgrade loop"
# Execute upgrade when execute-update file is found.
# This file is created by the orchestrator when trying to upgrade the program
while true; do
    if [ -f "/data/execute-update" ]; then
        echo "Found execute-update file, performing upgrade"
        solana program deploy gateway_upgrade.so --program-id "$GATEWAY_PROGRAM_ID"
        rm /data/execute-update
        echo "Upgrade completed and execute-update file removed"
    fi
    sleep 2
done

# leave some time for debug if validator exits due to errors
sleep 1000