#!/bin/bash

GIT_TAG=$(git describe --tags 2>/dev/null)
GIT_BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null)
GIT_HASH=$(git rev-parse --verify HEAD 2>/dev/null)
GIT_DIRTY=$(git diff --quiet; [ $? -eq 0 ] && echo false || echo true)

CGO_ENABLED=0 go build -ldflags "-X main.gitTag=$GIT_TAG -X main.gitBranch=$GIT_BRANCH -X main.gitHash=$GIT_HASH -X main.gitDirty=$GIT_DIRTY"
