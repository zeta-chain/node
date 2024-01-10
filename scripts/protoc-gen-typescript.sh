#!/usr/bin/env bash

DIR="typescript"

rm -rf $DIR

echo "Generate using buf."
(cd proto && buf generate --template buf.ts.yaml)

echo "create package json $DIR/package.json"
cat <<EOL > $DIR/package.json
{
  "name": "@zetachain/node-types",
  "version": "0.0.0-set-on-publish",
  "description": "",
  "main": "",
  "keywords": [],
  "author": "ZetaChain",
  "license": "MIT"
}
EOL
cat $DIR/package.json

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
