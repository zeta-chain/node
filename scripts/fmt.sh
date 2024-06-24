#!/bin/bash

# Exit on any error
set -e

if ! command -v golangci-lint &> /dev/null
then
    echo "golangci-lint is not found, installing..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
fi

if ! command -v golines &> /dev/null
then
    echo "golines could not be found, installing..."
    go install github.com/segmentio/golines@v0.12.2
fi

# Fix long lines
echo "Fixing long lines..."
golines -w --max-len=120 --ignore-generated --ignored-dirs=".git" --base-formatter="gofmt" .

# Gofmt, fix & order imports, remove whitespaces
echo "Formatting code..."
golangci-lint run --enable-only 'gci' --enable-only 'gofmt' --enable-only 'whitespace' --fix

echo "Code is formatted"
