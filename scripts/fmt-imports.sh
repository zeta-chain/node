#!/usr/bin/env bash

# Install goimports-revise
go install github.com/incu6us/goimports-reviser/v3@v3.6.4

# Run goimports-revise on all Go files
find . -name '*.go' -exec goimports-reviser -project-name github.com/zeta-chain/zetacore -file-path {} \; > /dev/null 2>&1

# Print a message to indicate completion
echo "Go imports formatted."

