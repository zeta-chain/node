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


# if [ ! -d "zeta-node" ]; then
#     echo "Local zeta-node source code directory not found inside the localnet/zetachain directory"
#     echo "If you already have the zeta-node repo saved locally it must be symbolically linked to this directory. Alternatively a new copy can be downloaded from GitHub"
#     echo "Press 'Y' if you want to download the source code from github.com/zeta-chain/zeta-node"
#     read -r -p  "Download from Github? (Y/n)" INPUT
#     echo $INPUT
#     if [[ $INPUT = [Yy] ]]; then
#         git clone git@github.com:zeta-chain/zeta-node.git
#         git pull
#     else 
#         echo "You must add a symbolic link to your local zeta-node source code inside the localnet/zetachain directory"
#         echo "The command will be similar to this 'ln -s ../../your/dir/here/zeta-node zeta-node'"
#         echo "Exiting - Try again when you have your local zeta-node repo linked to the localnet/zetachain directory"
#         exit 0
#     fi
# fi

# Create Docker Image
../../../build/build.sh


