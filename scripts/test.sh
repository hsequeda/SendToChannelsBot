#! /bin/bash
set -e

readonly env_file="$1"

env $(cat "./$env_file" | grep -Ev '^#' | xargs) go test --race -count=1 -p=8 -parallel=8 ./...
