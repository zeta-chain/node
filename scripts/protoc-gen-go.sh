#!/usr/bin/env bash

set -eo pipefail

echo "Generating gogo proto code"
echo "$PWD"
echo */
cd proto
echo "$PWD"
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  echo $dir
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

cd ..
echo "$PWD"
echo */

# Move proto files to the right places.
cp -r github.com/zeta-chain/zetacore/* ./
rm -rf github.com