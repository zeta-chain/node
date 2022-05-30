#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

for d in $(ls -d */); do 
  if [ $d != "node_modules/" ]; then
        echo $d
        cd $d
        ./stop.sh
        cd ..
    fi
done

# docker container prune -f
