#!/bin/bash

# run rpcimportable test with the required build flags

set -eo pipefail

cd "$(dirname "$0")"

go test -tags libsecp256k1_sdk -ldflags=all="-extldflags=-Wl,--allow-multiple-definition" ./...
