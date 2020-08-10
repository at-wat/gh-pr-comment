#!/bin/bash

set -eu

file=$(mktemp)
echo "Test file" > ${file}
ALLOW_PUBLIC_UPLOADER=true IMAGE_UPLOADER=test gh-pr-upload ${file} \
  && upload=OK || upload=Failed
echo

gh-pr-comment "Test post ${TRAVIS_OS_NAME} ${TRAVIS_CPU_ARCH}" \
  "testing comment post
- UTF-8 text: \"bœuf/牛\"
- upload test: ${upload}
"

[ ${upload} == "OK" ]
