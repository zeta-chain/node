#!/usr/bin/env bash

go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.16.2

go mod download

COSMOS_SDK="github.com/cosmos/cosmos-sdk"
ETHERMINT="github.com/evmos/ethermint"
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

DIR_PATH_ETHERMINT=$(go list -f '{{ .Dir }}' -m $ETHERMINT)

# Find the first OpenAPI YAML file and output its path
COSMOS_OPENAPI_PATH=$(find "$DIR_PATH" -type f -name "*.yaml" -exec grep -q "swagger:" {} \; -exec echo {} \; | head -n 1)

ETHERMINT_OPENAPI_PATH=$(find "$DIR_PATH_ETHERMINT" -type f -name "*.yaml" -exec grep -q "swagger:" {} \; -exec grep -l "/ethermint/evm/v1/" {} + | head -n 1)

# Generate OpenAPI YAML file using buf
buf generate --template $PROTO_TEMPLATE --output=$OUTPUT_DIR $ZETACHAIN_PROTO_PATH

# Initialize the merged swagger file with header info
echo "$OPENAPI_HEADER" > $OUTPUT_DIR/$MERGED_SWAGGER_FILE

# Extract paths and definitions from the cosmos swagger file and ZetaChain swagger file
# Merge them using yq and write to the merged swagger file
{
  echo "paths:"
  yq ea 'select(fileIndex == 0).paths * select(fileIndex == 1).paths * select(fileIndex == 2).paths' $COSMOS_OPENAPI_PATH $OUTPUT_DIR/$ZETACHAIN_OPENAPI $ETHERMINT_OPENAPI_PATH | sed -e 's/^/  /'
  echo "definitions:"
  yq ea 'select(fileIndex == 0).definitions * select(fileIndex == 1).definitions * select(fileIndex == 2).definitions' $COSMOS_OPENAPI_PATH $OUTPUT_DIR/$ZETACHAIN_OPENAPI $ETHERMINT_OPENAPI_PATH | sed -e 's/^/  /'
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
    # Extract NAME
    name=$(echo $key | cut -d '/' -f 3)
    # Check if the standard module is not imported in the app.go
    if ! grep -q "github.com/cosmos/cosmos-sdk/x/$name" $APP_GO; then
      # Keep the standard "base" and "tx" endpoints
      if [[ $name == "base" || $name == "tx" ]]; then
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
