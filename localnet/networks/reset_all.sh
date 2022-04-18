#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

rm -rf polygon/data/*
rm -rf zetachain/storage/*
rm -rf bsc/data/*

docker network create localnet --subnet 172.24.0.0/16

cd bsc/
./build.sh

cd ../zetachain
./generate_new_genesis_files.sh

