#!/bin/bash

echo "making an id"
solana-keygen new -o /root/.config/solana/id.json --no-bip39-passphrase

solana config set --url localhost
echo "starting solana test validator..."
solana-test-validator &

sleep 5
# airdrop to e2e sol account
solana airdrop 100
solana airdrop 100 37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ
solana program deploy gateway.so


# leave some time for debug if validator exits due to errors
sleep 1000