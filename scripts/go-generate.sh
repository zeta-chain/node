#!/usr/bin/env bash

set -e

# Install mockery.
go install github.com/vektra/mockery/v2@v2.53.3

DIRS=(
    "./testutil"
    "./zetaclient"
)

for dir in "${DIRS[@]}"; do
    (cd "$dir" && go generate ./... > /dev/null 2>&1)
done
