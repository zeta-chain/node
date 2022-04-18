#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

for d in $(ls -d */); do 
    echo $d
    cd $d
    ./stop.sh
    cd ..
done

# docker container prune -f
