#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

cp env_vars .env

# Build Bor Docker Image
git clone https://github.com/maticnetwork/bor.git
cd bor && git pull && cd ..

docker build . -t bor

# Build Heimdall Docker Image
git clone https://github.com/maticnetwork/heimdall.git
cd heimdall
make install network=local
docker build . -t heimdall 
