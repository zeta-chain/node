#!/usr/bin/env bash

# Generic function to generate bindings
function bindings() {
    cd $1
    go generate > /dev/null 2>&1
    echo "Generated bindings for $1"
}

# List of bindings to generate
bindings $(pwd)/precompiles/prototype

