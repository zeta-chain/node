#!/bin/bash

# Add keys
hermes keys add --chain zetachain --mnemonic "label they thank fitness pond noble honey friend another medal hedgehog door awake shoot walk stereo bubble attend tired front goat entire spot quick"
hermes keys add --chain mars --mnemonic "label they thank fitness pond noble honey friend another medal hedgehog door awake shoot walk stereo bubble attend tired front goat entire spot quick"

# Create connection and channels
hermes create connection --a-chain zetachain --b-chain mars
hermes create channel --a-chain zetachain --a-port transfer --b-chain mars --b-port transfer --order unordered --version ics20-1

# Start Hermes
hermes start