#!/bin/bash

NETWORK=${1}
CLEAN=${2}
FILE_PATH="${NETWORK}-docker-compose.yml"

if [ "${CLEAN}" == "true" ]; then
  docker-compose -f ${FILE_PATH} down -v
  rm -rf ${FILE_PATH}
else
  docker-compose -f ${FILE_PATH} down
  rm -rf ${FILE_PATH}
fi

