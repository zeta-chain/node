#!/bin/bash

echo "Running Build Tests"
go test -v -coverprofile coverage.out  $(go list ./... | grep -v /x/zetacore/)