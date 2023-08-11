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

#!/bin/bash

# Loop through all directories recursively
find "$DIR" -type d | while read -r dir; do
    # Check if there are any .d.ts files in the directory
    if ls "$dir"/*.d.ts &> /dev/null; then
        # Create or clear index.d.ts in the directory
        > "$dir/index.d.ts"

        # Loop through all .d.ts files in the directory
        for file in "$dir"/*.d.ts; do
            # Extract the base filename without the .d.ts extension
            base_name=$(basename "$file" .d.ts)
            # If the base name is not 'index', append the export line to index.d.ts
            if [ "$base_name" != "index" ]; then
                echo "export * from \"./$base_name\";" >> "$dir/index.d.ts"
            fi
        done
    fi
done
