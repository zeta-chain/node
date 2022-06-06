#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

echo "Sourcing Environment Variables from .env"
source .env

rm -rf polygon/data/*
rm -rf zetachain/storage/*
rm -rf bsc/data/*
rm -rf ganache/storage/*


docker network create localnet --subnet 172.24.0.0/16

if [ $USE_GANACHE == false ]; then
    echo "Launching Ganache Development Networks (Not Forked)"
    cd bsc/
    ./build.sh
    cd ..
fi

cd zetachain
./generate_new_genesis_files.sh

