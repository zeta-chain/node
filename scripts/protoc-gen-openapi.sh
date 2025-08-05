#!/usr/bin/env bash

set -eo pipefail

go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.26.1

go mod download

COSMOS_SDK="github.com/cosmos/cosmos-sdk"
PROTO_TEMPLATE="proto/buf.openapi.yaml"
OUTPUT_DIR="./docs/openapi"
MERGED_SWAGGER_FILE="openapi.swagger.yaml"
ZETACHAIN_OPENAPI="zetachain.swagger.yaml"
ZETACHAIN_PROTO_PATH="proto/"
APP_GO="app/app.go"
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

keys=$(yq e '.paths | keys' $OUTPUT_DIR/$MERGED_SWAGGER_FILE)

# Loop through each path
for key in $keys; do
  # Check if key starts with '/cosmos/NAME'
  if [[ $key == "/cosmos/"* ]]; then
    # Exclude paths starting with /cosmos/gov/v1beta1
    # these endpoints are broken post v0.47 upgrade
    if [[ $key == "/cosmos/gov/v1beta1"* ]]; then
      yq e "del(.paths.\"$key\")" -i $OUTPUT_DIR/$MERGED_SWAGGER_FILE
      continue
    fi

    # Extract NAME
    name=$(echo $key | cut -d '/' -f 3)
    # Check if the standard module is not imported in the app.go
    if ! grep -q "github.com/cosmos/cosmos-sdk/x/$name" $APP_GO; then
      # Keep the standard "base", "tx", and "upgrade" endpoints
      if [[ $name == "base" || $name == "tx" || $name == "upgrade" ]]; then
        continue
      fi
      # If not found, delete the key from the YAML file in-place
      yq e "del(.paths.\"$key\")" -i $OUTPUT_DIR/$MERGED_SWAGGER_FILE
    fi
  fi
done

# Remove extra white lines
awk '/^$/{if (++n <= 1) print; next}; {n=0;print}' $OUTPUT_DIR/$MERGED_SWAGGER_FILE > $OUTPUT_DIR/$MERGED_SWAGGER_FILE.temp
rm $OUTPUT_DIR/$MERGED_SWAGGER_FILE
mv $OUTPUT_DIR/$MERGED_SWAGGER_FILE.temp $OUTPUT_DIR/$MERGED_SWAGGER_FILE
