#!/usr/bin/env bash

# Install the required protoc execution tools.
go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest

# Define a shell function for generating proto code.
generate_proto() {
  local dir="$1"
  for file in "$dir"/*.proto; do
    if grep -q go_package "$file"; then
      if command -v buf >/dev/null 2>&1; then
        buf generate --template buf.gen.gogo.yaml "$file"
      else
        echo "Error: buf command not found. See https://docs.buf.build/installation" >&2
        exit 1
      fi
    fi
  done
}

# Generate Gogo proto code.
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  generate_proto "$dir"
done
cd ..

# Move proto files to the right places.
cp -r github.com/zeta-chain/zetacore/* ./
rm -rf github.com