#!/bin/bash

# --exact-match will ensure the tag is only returned if our commit is a tag
version=$(git describe --exact-match --tags 2>/dev/null)
if [[ $? -eq 0 ]]; then
    echo $version
    exit
fi

# use current timestamp for dirty builds
if ! git diff --no-ext-diff --quiet --exit-code ; then
    current_timestamp=$(date +"%s")
    echo "v0.0.${current_timestamp}-dirty"
    exit
fi

# otherwise assume we are on a develop build and use commit timestamp for version
commit_timestamp=$(git show --no-patch --format=%at)

echo "v0.0.${commit_timestamp}-develop"