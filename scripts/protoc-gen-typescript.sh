#!/usr/bin/env bash

cd proto
rm -rf typescript
buf generate --template buf.ts.yaml

# Set the target directory
DIR="../typescript"

# Loop through each file in the directory recursively
find "$DIR" -type f -name "*_pb.*" | while read -r file; do
    # Compute the new filename by removing '_pb'
    new_file="${file/_pb/}"
    
    # Rename the file
    mv "$file" "$new_file"
done