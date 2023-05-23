#!/usr/bin/env bash

set -eo pipefail

go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

COSMOS_SDK="github.com/cosmos/cosmos-sdk"
PROTO_TEMPLATE="proto/buf.openapi.yaml"
OUTPUT_DIR="./docs/openapi"
MERGED_SWAGGER_FILE="openapi.swagger.yaml"
ZETACHAIN_OPENAPI="zetachain.swagger.yaml"
ZETACHAIN_PROTO_PATH="proto/"
OPENAPI_HEADER="swagger: '2.0'
info:
  title: ZetaChain - gRPC Gateway docs
  description: Documentation for the API of ZetaChain.
paths:"

# Get the directory path for the cosmos-sdk
DIR_PATH=$(go list -f '{{ .Dir }}' -m $COSMOS_SDK)

# Find the first OpenAPI YAML file and output its path
COSMOS_OPENAPI_PATH=$(find "$DIR_PATH" -type f -name "*.yaml" -exec grep -q "swagger:" {} \; -exec echo {} \; | head -n 1)

# Generate OpenAPI YAML file using buf
buf generate --template $PROTO_TEMPLATE --output=$OUTPUT_DIR $ZETACHAIN_PROTO_PATH

# Initialize the merged swagger file
echo "$OPENAPI_HEADER" > $OUTPUT_DIR/$MERGED_SWAGGER_FILE

# Extract paths from the single swagger file
yq e '.paths' $COSMOS_OPENAPI_PATH | sed -e 's/^/  /' >> $OUTPUT_DIR/$MERGED_SWAGGER_FILE

# Append the generated OpenAPI YAML file to the merged file
yq e '.paths' $OUTPUT_DIR/$ZETACHAIN_OPENAPI | sed -e 's/^/  /' >> $OUTPUT_DIR/$MERGED_SWAGGER_FILE

# Check if the merged swagger file was created successfully
if [ -f $OUTPUT_DIR/$MERGED_SWAGGER_FILE ]; then
  # Remove the ZETACHAIN_OPENAPI file
  rm $OUTPUT_DIR/$ZETACHAIN_OPENAPI
fi
