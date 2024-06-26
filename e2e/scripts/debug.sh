#!/bin/bash

# Exit on any error
set -e

# Check if dlv is installed
if ! command -v dlv &> /dev/null
then
    echo "dlv could not be found, installing..."
    go install github.com/go-delve/delve/cmd/dlv@latest
fi

# Check if at least one argument is provided
if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <test> [args...]"
  exit 1
fi

# Extract the test argument
test=$1
shift

dlv_opts="--headless --listen=:2345 --api-version=2 --accept-multiclient"

# Collect additional arguments
e2e_test_args=$(echo "$@" | tr ' ' ',')
e2e_opts="--config cmd/zetae2e/config/local.yml"

# Echo commands
# shellcheck disable=SC2086
set -x
dlv debug ./cmd/zetae2e/ $dlv_opts -- run $test:$e2e_test_args $e2e_opts
