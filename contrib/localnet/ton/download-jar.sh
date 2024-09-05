#!/bin/bash

script_dir=$(cd -- "$(dirname -- "$0")" &> /dev/null && pwd)

# Downloads JAR outside of the Dockerfile to avoid re-downloading it every time during rebuilds.
jar_version="v120"
jar_url="https://github.com/neodix42/MyLocalTon/releases/download/${jar_version}/MyLocalTon-x86-64.jar"
jar_file="$script_dir/my-local-ton.jar"

if [ -f "$jar_file" ]; then
    echo "File $jar_file already exists. Skipping download."
    exit 0
fi

echo "File not found. Downloading..."
echo "URL: $jar_url"
wget -q --show-progress -O "$jar_file" "$jar_url"
