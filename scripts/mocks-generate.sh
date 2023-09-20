#!/usr/bin/env bash

# Install mockery
go install github.com/vektra/mockery/v2@latest

# Run generate command for mocks
cd ./testutil/keeper/mocks
go generate "mocks.go"

# Print a message to indicate completion
echo "Mocks generated."