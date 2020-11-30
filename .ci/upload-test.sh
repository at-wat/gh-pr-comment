#!/bin/bash

set -eu

file=$(mktemp)
echo "Test file" > ${file}
ALLOW_PUBLIC_UPLOADER=true IMAGE_UPLOADER=test gh-pr-upload ${file} \
  && upload=OK || upload=Failed
echo

gh-pr-comment "✔ Test post" \
  "- environment: $(uname -a || echo "No uname command")
- upload test: ${upload}
"

[ ${upload} == "OK" ]
