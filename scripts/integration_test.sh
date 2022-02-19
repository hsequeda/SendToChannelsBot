#! /bin/bash
set -e

readonly env_file="$1"

find . -wholename '.*\/adapter\/*test.go' -print0 | while read -d $'\0' file
do
  env $(cat "./$env_file" | grep -Ev '^#' | xargs) go test -v --race -count=1 -p=8 -parallel=8 $file
done
