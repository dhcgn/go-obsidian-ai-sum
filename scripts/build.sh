#!/bin/bash

# Parse arguments
OUTPUT_FILE=${1:-"go-obsidian-ai-sum"}  # Default to "go-obsidian-ai-sum" if no argument provided

# Get version info (format: v1.0.0-5-g1234567 or v1.0.0)
GIT_DESC=$(git describe --tags 2>/dev/null || echo "v0.0.0")

# Parse version components
BASE_TAG=$(echo $GIT_DESC | cut -d- -f1)
COMMITS_SINCE=$(echo $GIT_DESC | grep -o '[0-9]*-g[0-9a-f]*$' | cut -d- -f1)

# Build final version string
if [ -n "$COMMITS_SINCE" ]; then
  VERSION="${BASE_TAG}+${COMMITS_SINCE}"
else
  VERSION="${BASE_TAG}"
fi

COMMIT=$(git rev-parse --short HEAD)
BUILDTIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

# Build the main package instead of all packages
go build -v -o "${OUTPUT_FILE}" -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILDTIME}" ./main.go