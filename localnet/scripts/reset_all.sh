#!/bin/bash
LOCALNET_DIR="$( cd "$( dirname "$0" )" && pwd )/.."
cd "$LOCALNET_DIR" || exit

echo "Sourcing Environment Variables from .env"
source .env

cd chains
rm -rf polygon/data/*
rm -rf zetachain/storage/*
rm -rf bsc/data/*
rm -rf ganache/storage/*


docker network create localnet --subnet 172.24.0.0/16

if [ "$USE_GANACHE" == false ]; then
    echo "Launching Ganache Development Networks (Not Forked)"
    cd bsc/ || exit
    ./build.sh
    cd ..
fi

cd zetachain || exit
./create-genesis-files.sh

