#!/usr/bin/env bash

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
  description: Documentation for the API of ZetaChain."

# Get the directory path for the cosmos-sdk
DIR_PATH=$(go list -f '{{ .Dir }}' -m $COSMOS_SDK)

# Find the first OpenAPI YAML file and output its path
COSMOS_OPENAPI_PATH=$(find "$DIR_PATH" -type f -name "*.yaml" -exec grep -q "swagger:" {} \; -exec echo {} \; | head -n 1)

# Generate OpenAPI YAML file using buf
buf generate --template $PROTO_TEMPLATE --output=$OUTPUT_DIR $ZETACHAIN_PROTO_PATH

# Initialize the merged swagger file with header info
echo "$OPENAPI_HEADER" > $OUTPUT_DIR/$MERGED_SWAGGER_FILE

# Extract paths and definitions from the cosmos swagger file and ZetaChain swagger file
# Merge them using yq and write to the merged swagger file
{
  echo "paths:"
  yq ea 'select(fileIndex == 0).paths * select(fileIndex == 1).paths' $COSMOS_OPENAPI_PATH $OUTPUT_DIR/$ZETACHAIN_OPENAPI | sed -e 's/^/  /'
  echo "definitions:"
  yq ea 'select(fileIndex == 0).definitions * select(fileIndex == 1).definitions' $COSMOS_OPENAPI_PATH $OUTPUT_DIR/$ZETACHAIN_OPENAPI | sed -e 's/^/  /'
} >> $OUTPUT_DIR/$MERGED_SWAGGER_FILE

# Check if the merged swagger file was created successfully
if [ -f $OUTPUT_DIR/$MERGED_SWAGGER_FILE ]; then
  # Remove the ZETACHAIN_OPENAPI file
  rm $OUTPUT_DIR/$ZETACHAIN_OPENAPI
fi
