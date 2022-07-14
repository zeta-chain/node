#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

echo "You Entered $INPUT"

if  [ "$INPUT" == "" ]; then
    echo "Building zetacore Only"
    docker build -f ../Dockerfile.zetacore ../  -t zetacore
elif  [ "$INPUT" == "cp-local-binary" ]; then
    echo "Building zetacore Only"
    docker build -f ../Dockerfile.zetacore_binary ../  -t zetacore
elif [ "$INPUT" == "mockmpi" ]; then
    echo "Building $INPUT Only"
    docker build -f ../Dockerfile.mockmpi ../ -t zeta-mockmpi
else 
    echo "Unknown Input"
    echo "Enter zetacore, mockmpi, or include no argument at all"
fi
