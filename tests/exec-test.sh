#!/bin/bash

set -eu

file=$(mktemp)
echo "Test file" > ${file}
ALLOW_PUBLIC_UPLOADER=true IMAGE_UPLOADER=test gh-pr-upload ${file} \
  && upload=OK || upload=Failed

gh-pr-comment "Test post (python ${TRAVIS_PYTHON_VERSION})" "testing comment post on python ${TRAVIS_PYTHON_VERSION}
- UTF-8 text: \"bœuf/牛\"
- upload test: ${upload}
"
[ ${upload} == "OK" ]
