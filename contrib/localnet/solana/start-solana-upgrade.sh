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
solana program deploy gateway.so
solana program deploy connected.so
# upgrade to the new program, gateway-upgrade.so . The new program is identical to the old program, but has an extra field
solana program deploy gateway-upgrade.so --program-id 94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d

# leave some time for debug if validator exits due to errors
sleep 1000

