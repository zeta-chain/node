#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

echo "You Entered $INPUT"

if  [ "$INPUT" == "" ]; then
    echo "Building zetacore Only"
    docker build -f ../Dockerfile.zetacore ../  -t zetacore
elif [ "$INPUT" == "zetaclient" ]; then
    echo "Building $INPUT Only"
    docker build -f ../Dockerfile.zetaclient ../  -t zetaclient
elif [ "$INPUT" == "mockmpi" ]; then
    echo "Building $INPUT Only"

    docker build -f ../Dockerfile.mockmpi ../ -t zeta-mockmpi
else 
    echo "Unknown Input"
    echo "Enter zetacore, zetaclient, mockmpi, or include no argument at all. If no argument is provided all three images will be built."
fi
