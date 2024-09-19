#!/usr/bin/env bash
# set -x

ABIGEN_VERSION=1.14.8-stable
ABIGEN_VERSION_SOURCE=github.com/ethereum/go-ethereum/cmd/abigen@v1.14.8
ABIGEN_INSTALL_CMD="go install $(echo $ABIGEN_VERSION_SOURCE)"

SOLC_VERSION=0.8.26
SOLC_SELECT_CMD_INSTALL="solc-select install $(echo $SOLC_VERSION)"
SOLC_SELECT_CMD_USE="solc-select use $(echo $SOLC_VERSION)"

install_abigen() {
    echo "Installing abigen version $ABIGEN_VERSION..."
    $ABIGEN_INSTALL_CMD
    if [ $? -ne 0 ]; then
        echo "Error: Failed to install abigen."
        exit 1
    fi
    echo "abigen version $ABIGEN_VERSION installed successfully."
}

# Check if abigen is installed
if command -v abigen &> /dev/null; then
    INSTALLED_ABIGEN_VERSION=$(abigen --version | grep -o "$ABIGEN_VERSION")
    if [ "$INSTALLED_ABIGEN_VERSION" == "$ABIGEN_VERSION" ]; then
        echo "abigen version $ABIGEN_VERSION is already installed."
    else
        echo "abigen version $ABIGEN_VERSION not found, installing..."
        install_abigen
    fi
else
    echo "abigen not found, installing..."
    install_abigen
fi

# Check if solc is installed and at version 0.8.26
if command -v solc &> /dev/null
then
    INSTALLED_SOLC_VERSION=$(solc --version | grep -o "$SOLC_VERSION")
    if [ "$INSTALLED_SOLC_VERSION" == "$SOLC_VERSION" ]; then
        echo "solc version $SOLC_VERSION is already installed."
    else
        echo "solc is installed but not version $SOLC_VERSION. Checking for solc-select..."
        if command -v solc-select &> /dev/null
        then
            echo "solc-select found, installing and using solc $SOLC_VERSION."
            $SOLC_SELECT_CMD_INSTALL
            $SOLC_SELECT_CMD_USE
        else
            echo "solc-select not found. Please install solc-select or ensure solc $SOLC_VERSION is available."
            exit 1
        fi
    fi
else
    echo "solc is not installed. Checking for solc-select..."
    if command -v solc-select &> /dev/null
    then
        echo "solc-select found, installing and using solc $SOLC_VERSION."
        $SOLC_SELECT_CMD_INSTALL
        $SOLC_SELECT_CMD_USE
    else
        echo "solc-select not found. Please install solc-select or ensure solc $SOLC_VERSION is available."
        exit 1
    fi
fi

# Generic function to generate bindings
function bindings() {
    cd $1
    go generate > /dev/null 2>&1
    echo "Generated bindings for $1"
    cd - > /dev/null 2>&1
}

# List of bindings to generate
bindings ./precompiles/prototype
bindings ./precompiles/staking

