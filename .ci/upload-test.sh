#!/bin/bash

set -eu

export PATH="${PATH}:$(go env GOPATH)/bin"
go install ./...

file=$(mktemp)
echo "Test file" > ${file}
ALLOW_PUBLIC_UPLOADER=true IMAGE_UPLOADER=test gh-pr-upload ${file} \
  && upload=OK || upload=Failed
echo

gh-pr-comment "✔ $1" "upload test: ${upload}"

[ ${upload} == "OK" ]
