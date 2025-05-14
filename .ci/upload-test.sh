#!/bin/bash

set -eu

export PATH="${PATH}:$(go env GOPATH)/bin"
go install ./...

file=$(mktemp)
echo "Test file" > ${file}

echo "Testing test uploader"
ALLOW_PUBLIC_UPLOADER=true \
  IMAGE_UPLOADER=test \
  gh-pr-upload ${file} \
    && upload_test=OK \
    || upload_test=Failed
echo

if [ "$(uname -o)" = "GNU/Linux" ]; then
  echo "Testing s3 uploader"
  IMAGE_UPLOADER=s3 \
    AWS_REGION=ap-northeast-1 \
    AWS_S3_BUCKET=test-bucket \
    AWS_ENDPOINT_URL=http://localhost:4566 \
    AWS_S3_USE_PATH_STYLE=true \
    AWS_ACCESS_KEY_ID=test \
    AWS_SECRET_ACCESS_KEY=test \
    gh-pr-upload ${file} \
      && upload_s3=OK \
      || upload_s3=Failed
fi

rm ${file}

gh-pr-comment "✔ $1" "test uploader: ${upload_test}
s3 uploader: ${upload_s3:-n/a}
$(go version)
"

[ ${upload_test} = "OK" ]
[ ${upload_s3:-OK} = "OK" ]
