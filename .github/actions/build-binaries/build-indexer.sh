#!/bin/bash

echo "Building indexer"
# go build -mod=readonly ./cmd/indexer

make install-indexer
cp "$HOME"/go/bin/* ./