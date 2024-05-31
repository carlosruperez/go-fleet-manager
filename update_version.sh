#!/bin/bash

# Check if a version is provided as an argument
if [ -z "$1" ]; then
    echo "No version provided. Usage: ./update_version.sh <new_version>"
    exit 1
fi

# New version passed as an argument
new_version=$1

# Path to the version file
version_file="version.yaml"

# Check if the file exists
if [ ! -f "$version_file" ]; then
    echo "Version file does not exist: $version_file"
    exit 1
fi

# Update the version in the version.go file
echo $new_version
sed -i.bak "s/^version: .*/version: $new_version/" $version_file

# Check if the sed command was successful
if [ $? -ne 0 ]; then
    echo "Failed to update the version."
    exit 1
fi

# Remove the backup file created by sed
rm "${version_file}.bak"
