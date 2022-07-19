#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

if [ ! -f ".env" ]; then
    cp env_vars .env
fi

# Create Directories for Local Config Files
mkdir -p config/genesis
mkdir -p config/node0/data config/node0/config config/node0/keyring-test/
mkdir -p config/node1/data config/node1/config config/node1/keyring-test/
mkdir -p config/node2/data config/node2/config config/node2/keyring-test/
mkdir -p config/node3/data config/node3/config config/node3/keyring-test/

# Create Docker Image
../../../build/build.sh


