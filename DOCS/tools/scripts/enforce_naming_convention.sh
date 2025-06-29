#!/bin/bash

# Copyright (c) 2025 A Bit of Help, Inc.

# enforce_naming_convention.sh
# This script enforces the naming convention for .puml and .svg files in the DOCS/diagrams directory.
# The naming convention is to use underscores and lowercase for all file names.

# Set the directory to check
DIAGRAMS_DIR="DOCS/diagrams"

# Function to convert a filename to the correct naming convention
convert_filename() {
    local filename="$1"
    # Convert to lowercase and replace spaces and hyphens with underscores
    echo "$filename" | tr '[:upper:]' '[:lower:]' | tr ' -' '_'
}

# Function to check and rename files if needed
check_and_rename() {
    local file="$1"
    local dir=$(dirname "$file")
    local filename=$(basename "$file")
    local correct_filename=$(convert_filename "$filename")
    
    # If the filename doesn't match the correct naming convention, rename it
    if [ "$filename" != "$correct_filename" ]; then
        echo "Renaming $file to $dir/$correct_filename"
        mv "$file" "$dir/$correct_filename"
        
        # If this is a .puml file, also update the corresponding .svg file if it exists
        if [[ "$file" == *.puml ]]; then
            local svg_file="${file%.puml}.svg"
            local correct_svg_file="${dir}/${correct_filename%.puml}.svg"
            if [ -f "$svg_file" ]; then
                echo "Renaming $svg_file to $correct_svg_file"
                mv "$svg_file" "$correct_svg_file"
            fi
        fi
    fi
}

# Main function to process all files in the directory
process_directory() {
    local dir="$1"
    
    # Check if the directory exists
    if [ ! -d "$dir" ]; then
        echo "Error: Directory $dir does not exist."
        exit 1
    fi
    
    # Process all .puml and .svg files in the directory
    find "$dir" -type f \( -name "*.puml" -o -name "*.svg" \) | while read file; do
        check_and_rename "$file"
    done
    
    echo "Naming convention check completed."
}

# Run the script
process_directory "$DIAGRAMS_DIR"

# Exit with success
exit 0