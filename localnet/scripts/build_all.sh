#!/bin/bash

LOCALNET_DIR="$( cd "$( dirname "$0" )" && pwd )/.."
cd "$LOCALNET_DIR" || exit
source scripts/rpc_commands
source .env
docker network create localnet --subnet 172.24.0.0/16 >> /dev/null 2>&1

# Download ZetaChain Mono Repo 
git clone -b "$ZETA_MONOREPO_TAG" https://github.com/zeta-chain/zetachain.git zetachain-monorepo
cd zetachain-monorepo || exit
yarn 
yarn compile 
cd .. || exit

# Build Images for each chain
cd chains || exit
for d in $(ls -d */); do 
    echo "$d"
    cd "$d" || exit
    ./build.sh
    cd ..
done

# Generate ZetaChain Genesis file if needed
FILE="${LOCALNET_DIR}/zetachain/config/genesis/genesis.json"
if  [ -f "$FILE" ]; then
    echo "Zetachain Genesis File Already Exists. Not Regenerating"
else
    echo "Zetachain Genesis File NOT Found"
    echo "Generating New Zetachain Genesis/Config Files"
    ${LOCALNET_DIR}/zetachain/create-genesis-files.sh
fi

cd "$LOCALNET_DIR" || exit
yarn install
