#!/bin/bash

# Exit on any error
set -e

# Make sure that dlv is installed!
# go install github.com/go-delve/delve/cmd/dlv@latest

# Check if at least one argument is provided
if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <test> [args...]"
  exit 1
fi

# Extract the test argument
test=$1
shift

# Collect additional arguments
e2e_test_args=$(echo "$@" | tr ' ' ',')
e2e_opts="--config cmd/zetae2e/config/local.yml --verbose"

# Echo commands
# shellcheck disable=SC2086
set -x
go run ./cmd/zetae2e/ run $test:$e2e_test_args $e2e_opts
