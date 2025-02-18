#!/bin/bash

echo "making an id"
solana-keygen new -o /root/.config/solana/id.json --no-bip39-passphrase

solana config set --url localhost
echo "starting solana test validator..."
solana-test-validator &

sleep 5
# airdrop to e2e sol account
solana airdrop 1000
solana airdrop 1000 37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ

# Deploy initial programs
solana program deploy gateway.so
solana program deploy connected.so
solana program deploy connected_spl.so

# Get program ID from gateway keypair
GATEWAY_PROGRAM_ID=$(solana-keygen pubkey gateway-keypair.json)

# upgrade to the gateway-upgrade program using the program ID from keypair
solana program deploy gateway_upgrade.so --program-id "$GATEWAY_PROGRAM_ID"

# leave some time for debug if validator exits due to errors
sleep 1000
