#!/usr/bin/env bash

set -eo pipefail

# Check if abigen is installed
if ! command -v abigen &> /dev/null
then
    echo "abigen could not be found, installing..."
    go install github.com/ethereum/go-ethereum/cmd/abigen@latest
fi

# Generic function to generate bindings
function bindings() {
    cd $1
    go generate > /dev/null 2>&1
    echo "Generated bindings for $1"
}

# List of bindings to generate
bindings ./precompiles/prototype

