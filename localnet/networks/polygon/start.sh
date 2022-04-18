#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"
cd $DIR

docker compose up -d --remove-orphans