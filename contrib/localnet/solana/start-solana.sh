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

# test accounts that can be used in e2e tests (eg. for withdraw and call)
E2E_AIRDROP_ADDRS=(
  4C2kkMnqXMfPJ8PPK5v6TCg42k5z16f2kzqVocXoSDcq
  3tqcVCVz5Jztwnku1H9zpvjaWshSpHakMuX4xUJQuuuA
  7duGsuv6nB3yr15EuWuHEDD7rWovpAnjuveXJ5ySZuFV
  8vjuCrCKVfnBGinWjc33zLRnG8iy53wj3YWHhKqvTE7o
  bzkoxG5YMeWxKfNcjzbEHb3XaGTY4NLKfejjDmxVhhY
  GUjKWPmnXNwPLR6kcrkLSBARmQdYnySPpxVNEUGFLs72
  5oqdTyA78hpeP8RTwRBmoxCvp1V7DFicKj7T2DvtDDQM
  C481t79gpWbsWwPD9eJZTAo5TSaTBet8icEkiPhwKLDx
  EJvNNovWkfQYrmyMncqVvHNde2QJpSA4EJk355vyQWph
  EkUpd7HFbSYPJEDbXZeDsCG19Hj5vbKTUt4rzpPYKsTM
  7c7TqqdbKRWDNLVAxNa481F355AAij1fdSRttzUNVNeD
  FuefjNTywey57U2zW6SmBWaJsCx7E84jmWUCbk52sBHR
  fczqc5N5arnKbvMj1kgg9P1FpYPQBmJsyTXcrmK9bbp
  GyH4mpobR6g2npNo5vRNcs2Cv8CxDKbK82p34kNCB4p2
  9XYF8U1srAkUETkywtFrZsApipENiMU3C8Rnz5aPib94
  ABnm3PMB4onvCFriWg7eBcNSmiye9iq9rRdeBViJqWif
  AEydLk3RXZv67wry7EMZmS1uHLYdH8ia6xdsYri4hyB2
  ArDNFdmDzrRP13UTJ4nyP11NfjV6aiQrxNwFUnwm3h8N
  FzAt9aPKFUy1D2Qq8u7myYW8HFqKLzbQS2paaWz8iAmZ
  9RybduN4CJHaZXvUiHoZ7KsHS9dgv1NCSAdJZaRJDW5U
  C54jMgtk2umaJYoPD8aF3hmuH8XkAz2xA4sxr2ZtJABV
  8kRqLbezvj4apyaK6fanurQJjhDQwn6wnW1Yr9H96gsT
  2fddFSJoGJ2YuAZWxEvK9pRXXWLSKJ44rfJNVZg5WHBn
  GW9oi4yqAUFBcUNUHqz56FRdfHU6Md1t9o7i2svL1XcG
  Gre1nqrE1KyBbBH2Xb4qVQWFFdyXetLR5rruJpJkNhkV
  BUrJRTsFVqnuLeq4Vfs7NrrU6EsL957n4tBZGoRRxkza
  3tE2kKyPfuwC5rZUpJ9NaKMBf5og7G3rVTzE9CNMm9Q9
  CFEQ79VSAupXWmdzvmjNtef5BybDVYXBKeh8frNHZiYe
  bTnsajQuybXV6Wf9V7a8wQrwaWZ1WskkmUUChimkWmc
  4dpxhhomWY3A9ey9g7EfxRQKfXnikfM4tRrFAgQf8Y8n
  aQJPrcj4LNNHh9UK41sfcACFspaFR7wgcUTgmSKRXiB
  5ndhaFZ48eKyU7f66vq7WSbjZRh9WhpnbgdMwDRrvgj4
  G7m7dSWH5tb1WC2g86vqA1UvdKJVNCU9td1TR8j8wQXo
  3tf2MkQzmLHBjnsmRwKnJQgASrmUggxK2Q3PiFd99tDn
  3X31YYsRw8We2YhsK29QtVwXXk783HbbYDGAV1HBjcBD
  FhDjU5r4MWx6KdfY7MVx6w1YJf6tvRJyj6mbxqU8N3F7
  57xqRiBeQjHgrYqhXUHRKh9WU9Ukya79U3hM7QPr13Gy
  7pjSLC42Er4KPdVLZW7VGkxU9tLZKBqt4apwyZzuWyYU
  jyDrCsnuvxGM9H7YE2rLsBYoFoSHWABFoi7P61FaQ3Q
  GUHXNHugMc22rkX6Mz4GMU4Vj1hbPa3DCrfCDRQFWQ2b
  2azGMpfp91pqd5gXpZJWK8egdpgUxDkhXK8UYHtRjiZa
  CPH4QdmL4yNB9KBpr3bQwUQxdQbMeUKsUyos2ViHaNTB
  GKNqfGFsK1Th32GSRkT9kaiA7w89GKJJmGVV5ibo8xn2
  GzuWB5nf2NH15Ssk15n3Zd72iykMoPm8Qx5TPCUS99LC
  7eAojuq3vcrev41DuVFgpZ1yagQhpUNhUn2uXnKo7A41
  BCuFo9AhTREJ5bgJzCrhXzckmAoyscrDDkggixX6t3c5
  FUnpGc7v43bvBvC584gQXiRdcuCMnDoXLXbJmMNkg3wQ
  FtFPeHGXZhgacdNoXh2dYKBZkgmTq9YYoUZX97hyjgh4
  2f3V4h5z9jds59EeFVqViVKuZrMoYM3xb3eq8fWWuN7Y
  BEvgtgRX7DdUrZ8Jrw5SMLctA7pQ76ScGry73mEzH869
  FD8pHBAwhq2VtHQQSTpdpnmiEoMNfFzJhQXGhSQVkdcQ
  HqQuQ9wF3QE7RwiYdgB88SnHwi1n5Q2ogidy7WfZJGgb
  Dg2fDYcuvRxCBtZtf1rB2bgbc4KTmy8KMC22Y4JFX7Qd
  p6pGSE2rLDH7yiiZ7bKoNZcB3YaszRsRDQ5rVRuTXiz
  HEg2w4Ev5ouoZB51Tmhj4DBPG7jrxKaTrf9GfKYubBbG
  A5mcmJHSMARvaQcYXGQ96Nx1h4sFeReJNTCBxqxxMqrF
)

for addr in "${E2E_AIRDROP_ADDRS[@]}"; do
    solana airdrop 1 "$addr"
done

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