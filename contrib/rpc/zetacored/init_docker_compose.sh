#!/bin/bash

NETWORK=${1}
TYPE=${2}
DOCKER_TAG=${3}

if [ "$TYPE" == "image" ]; then
  echo "Source Environment File."
  SOURCE_FILE_NAME="networks/.${NETWORK}"
  if [ ! -f "$SOURCE_FILE_NAME" ]; then
    echo "Environment file $SOURCE_FILE_NAME does not exist."
    exit 1
  fi
  source ${SOURCE_FILE_NAME}
elif [ "$TYPE" == "localbuild" ]; then
  echo "Source Environment File."
  SOURCE_FILE_NAME="networks/.${NETWORK}-localbuild"
  if [ ! -f "$SOURCE_FILE_NAME" ]; then
    echo "Environment file $SOURCE_FILE_NAME does not exist."
    exit 1
  fi
  source ${SOURCE_FILE_NAME}
fi

# Define the path to the Docker Compose file
FILE_PATH="${NETWORK}-docker-compose.yml"
cp docker-compose.yml ${FILE_PATH}

# Determine the appropriate Docker Compose configuration based on TYPE
if [ "$TYPE" == "image" ]; then
    IMAGE_BLOCK="image: zetachain/zetacored:\${DOCKER_TAG:-ubuntu-v14.0.1.0}"
    NAME="zetacored-rpc-${NETWORK}"
elif [ "$TYPE" == "localbuild" ]; then
    IMAGE_BLOCK=$(cat << 'EOF'
build:
      context: ../../..
      dockerfile: Dockerfile
EOF
)
  NAME="zetacored-rpc-${NETWORK}-localbuild"
else
    echo "Invalid TYPE. Please specify 'image' or 'localbuild'."
    exit 1
fi

IMAGE_BLOCK_ESCAPED=$(echo "$IMAGE_BLOCK" | sed 's/[&/]/\\&/g; s/$/\\/')
IMAGE_BLOCK_ESCAPED=${IMAGE_BLOCK_ESCAPED%?}

# Replace placeholders in the Docker Compose file
sed -i '' "s|-=name=-|$NAME|g" $FILE_PATH
sed -i '' "s|-=image_block=-|$IMAGE_BLOCK_ESCAPED|g" $FILE_PATH

echo "DEBUG ENV VARS"
printenv
echo "================"

echo "Placeholders have been replaced in $FILE_PATH."
cat $FILE_PATH
echo "================"

if [ "$TYPE" == "image" ]; then
  docker-compose -f ${FILE_PATH} up
elif [ "$TYPE" == "localbuild" ]; then
  docker-compose -f ${FILE_PATH} up --build
fi
