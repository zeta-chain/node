#!/bin/bash

echo "making an id"
solana-keygen new -o /root/.config/solana/id.json --no-bip39-passphrase

solana config set --url localhost
echo "starting solana test validator..."
solana-test-validator &

sleep 5
# airdrop to e2e sol account and rent payer (used to generate atas for withdraw spl receivers if they don't exist)
solana airdrop 100
solana airdrop 100 37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ
solana airdrop 100 C6KPvGDYfNusoE4yfRP21F8wK35bxCBMT69xk4xo3X79
solana program deploy gateway.so


# leave some time for debug if validator exits due to errors
sleep 1000