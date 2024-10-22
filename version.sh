#!/bin/bash

# if NODE_VERSION is set, just return it immediately
# this allows you to do docker builds without the git directory
if [[ -n $NODE_VERSION ]]; then
    echo $NODE_VERSION
    exit
fi

# --exact-match will ensure the tag is only returned if our commit is a tag
version=$(git describe --exact-match --tags 2>/dev/null)
if [[ $? -eq 0 ]]; then
    # do not return v prefix on the version
    # to match the goreleaser logic
    echo ${version#v}
    exit
fi

# develop build and use commit timestamp for version
commit_timestamp=$(git show --no-patch --format=%at)
short_commit=$(git rev-parse --short HEAD)

# append -dirty for dirty builds
if ! git diff --no-ext-diff --quiet --exit-code ; then
    echo "0.0.${commit_timestamp}-${short_commit}-dirty"
    exit
fi

echo "0.0.${commit_timestamp}-${short_commit}"