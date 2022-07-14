#!/bin/bash
LOCALNET_DIR="$( cd "$( dirname "$0" )" && pwd )/.."
cd "$LOCALNET_DIR" || exit

cd chains
for d in $(ls -d */); do 
  if [ "$d" != "node_modules/" ]; then
        echo "$d"
        cd "$d" || exit
        ./stop.sh
        cd ..
    fi
done

# docker container prune -f
