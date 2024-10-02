#!/bin/bash

# run rpcimportable test with the required build flags

set -eo pipefail

cd "$(dirname "$0")"

# --allow-multiple-definitions need to be set when you are importing both cosmos-sdk
# and go-ethereum: https://github.com/cosmos/cosmos-sdk/tree/release/v0.47.x/crypto/keys/secp256k1/internal/secp256k1
#
# enable libsecp256k1_sdk to bypass the btcec breaking change:
# https://github.com/btcsuite/btcd/issues/2243
go test -tags libsecp256k1_sdk -ldflags="-extldflags=-Wl,--allow-multiple-definition" ./...
