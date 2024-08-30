#!/usr/bin/env bash

# Check if abigen is installed
if ! command -v abigen &> /dev/null
then
    echo "abigen could not be found, installing..."
    go install github.com/ethereum/go-ethereum/cmd/abigen@latest
fi

# Check if solc is installed and at version 0.8.26
if command -v solc &> /dev/null
then
    SOLC_VERSION=$(solc --version | grep -o "Version: 0.8.26")
    if [ "$SOLC_VERSION" == "Version: 0.8.26" ]; then
        echo "solc version 0.8.26 is already installed."
    else
        echo "solc is installed but not version 0.8.26. Checking for solc-select..."
        if command -v solc-select &> /dev/null
        then
            echo "solc-select found, installing and using solc 0.8.26..."
            solc-select install 0.8.26
            solc-select use 0.8.26
        else
            echo "solc-select not found. Please install solc-select or ensure solc 0.8.26 is available."
            exit 1
        fi
    fi
else
    echo "solc is not installed. Checking for solc-select..."
    if command -v solc-select &> /dev/null
    then
        echo "solc-select found, installing and using solc 0.8.26..."
        solc-select install 0.8.26
        solc-select use 0.8.26
    else
        echo "solc or solc-select could not be found. Please install one of them to proceed."
        exit 1
    fi
fi

# Generic function to generate bindings
function bindings() {
    cd $1
    go generate > /dev/null 2>&1
    echo "Generated bindings for $1"
}

# List of bindings to generate
bindings ./precompiles/prototype

