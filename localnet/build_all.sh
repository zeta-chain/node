#!/bin/bash

DIR="$( cd "$( dirname "$0" )" && pwd )"

cd $DIR
source rpc_commands
docker network create localnet --subnet 172.24.0.0/16

for d in $(ls -d */); do 
    echo $d
    cd $d
    ./build.sh
    cd ..
done

FILE="$DIR/zetachain/config/genesis/genesis.json"
if  [ -f "$FILE" ]; then
    echo "Zetachain Genesis File Already Exists. Not Regenerating"
else
    echo "Zetachain Genesis File NOT Found"
    echo "Generating New Zetachain Genesis/Config Files"
    $DIR/zetachain/generate_new_genesis_files.sh
fi

cd $DIR
yarn install
