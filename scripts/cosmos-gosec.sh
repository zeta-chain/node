#!/usr/bin/env bash

# Install gosec
go install github.com/cosmos/gosec/v2/cmd/gosec@latest

# Run gosec
gosec ./... -include=G701,G703,G704