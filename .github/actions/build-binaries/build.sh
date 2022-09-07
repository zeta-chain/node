#!/bin/bash

echo "Building zetacored and zetaclients"
# go build -mod=readonly ./cmd/zetacored
# go build -mod=readonly ./cmd/zetaclientd

make install
cp "$HOME"/go/bin/* ./
