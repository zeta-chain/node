#!/usr/bin/env bash

image=ghcr.io/zeta-chain/gosec:2.21.4-zeta2

docker run -it --rm -w /node -v "$(pwd):/node" \
  -e GO111MODULE=on -e GOTOOLCHAIN=auto \
  $image -exclude-generated -exclude-dir testutil ./...

