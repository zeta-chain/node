#!/usr/bin/env bash

# Define a shell function for iterating and formatting proto files
format_proto() {
  local dir="$1"
  for file in "$dir"/*.proto; do
    if grep -q go_package "$file"; then
      if command -v buf >/dev/null 2>&1; then
        buf format -w "$file"
      else
        echo "Error: buf command not found. See https://docs.buf.build/installation" >&2
        exit 1
      fi
    fi
  done
}

# Format Gogo proto code.
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  format_proto "$dir"
done
cd ..