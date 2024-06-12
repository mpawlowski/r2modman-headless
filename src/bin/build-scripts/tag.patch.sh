#!/bin/bash

git fetch --tags --force
latestTag=$(git tag | grep v | sort -V | tail -n 1)

# Split the version into major, minor, patch
major=$(echo $latestTag | awk -F. '{print $1}')
minor=$(echo $latestTag | awk -F. '{print $2}')
patch=$(echo $latestTag | awk -F. '{print $3}')

# Check if major is empty
if [ -z "$major" ]; then
    major="v0"
fi

# Check if minor is empty
if [ -z "$minor" ]; then
    minor=0
fi

# Check if patch is empty
if [ -z "$patch" ]; then
    patch=0
fi

# Increment the minor version
patch=$((patch + 1))

# Assemble the new version
newTag="$major.$minor.$patch"

echo "New Version ($major.$minor.$patch) -> ($newTag)"

read -p "Verify new version: " inputTag

if [ "$inputTag" != "$newTag" ]; then
    echo "Incorrect tag entered. Exiting..."
    exit 1
fi

git tag $newTag
git push origin $newTag
