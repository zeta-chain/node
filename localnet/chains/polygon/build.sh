#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd "$DIR" || exit

cp env_vars .env

# Build Bor Docker Image
git clone https://github.com/maticnetwork/bor.git >> /dev/null 2>&1
cd bor && git checkout tags/v0.2.16 >> /dev/null 2>&1 && cd ..
docker build . -t bor

# # Build Heimdall Docker Image
# git clone https://github.com/maticnetwork/heimdall.git
# cd heimdall || exit
# make install network=local
# docker build . -t heimdall 
