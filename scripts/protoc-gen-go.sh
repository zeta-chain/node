#!/usr/bin/env bash

set -eo pipefail

export BUF_CACHE_DIR="/tmp/buf-cache"

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

cd ..

# Move proto files to the right places.
cp -r github.com/zeta-chain/node/* ./
rm -rf github.com