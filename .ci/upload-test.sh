#!/bin/bash

set -eu

file=$(mktemp)
echo "Test file" > ${file}
ALLOW_PUBLIC_UPLOADER=true IMAGE_UPLOADER=test gh-pr-upload ${file} \
  && upload=OK || upload=Failed
echo

gh-pr-comment "✔ Test post" \
  "- environment: $(uname -o && echo -n " " && uname -m || echo "No uname command")
- upload test: ${upload}
"

[ ${upload} == "OK" ]
