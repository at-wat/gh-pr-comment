#!/bin/bash

set -eu

gh-pr-comment "Test post (python ${TRAVIS_PYTHON_VERSION})" "testing comment post on python ${TRAVIS_PYTHON_VERSION}
- UTF-8 text: \"bœuf/牛\"
"
