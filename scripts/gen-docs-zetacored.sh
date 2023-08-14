#!/bin/bash

DIR=docs/cli/zetacored

rm -rf $DIR

go install ./cmd/zetacored

zetacored docs --path $DIR

# Recursive function to process files
process_files() {
    local dir="$1"
    
    # Process all files in the directory
    for file in "$dir"/*; do
        if [ -f "$file" ]; then
            # Replace <...> with [...] in the file, otherwise Docusaurus thinks it's a link
            sed -i.bak 's/<\([^<>]*\)>/\[\1\]/g' "$file"

            # Modify the heading by replacing ## zetacored with #
            sed -i.bak 's/^## zetacored /# /g' "$file" 

            # Replace the (default "/path/to/ with (default "~/
            sed -i.bak 's/\(default "\)[^"]*\//\1~\//g' "$file"
            
            # Remove the backup files
            rm -f "$file.bak"
        elif [ -d "$file" ]; then
            # Recurse into subdirectory
            process_files "$file"
        fi
    done
}

# Start processing from the given directory
process_files $DIR
