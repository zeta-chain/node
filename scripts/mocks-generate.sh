#!/usr/bin/env bash

set -e

# Install mockery
go install github.com/vektra/mockery/v2@v2.53.3

MOCK_DIRS=(
    "./testutil/keeper/mocks"
    "./zetaclient/chains/bitcoin/client"
    "./zetaclient/chains/evm/observer"
    "./zetaclient/chains/ton/observer"
    "./zetaclient/chains/ton/signer"
    "./zetaclient/testutils/mocks"
)

for dir in "${MOCK_DIRS[@]}"; do
    (cd "$dir" && go generate > /dev/null 2>&1)
done

# Print a message to indicate completion
echo "Mocks generated."
