#!/bin/bash

set -eu

gh-pr-comment "Test post" \
  "testing comment post
- UTF-8 text: \"bœuf/牛\"
"
