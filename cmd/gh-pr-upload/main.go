package main

import (
	"context"
	"fmt"
	"os"

	"github.com/at-wat/gh-pr-comment/cmd/gh-pr-upload/uploader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s filename\n", os.Args[0])
		fmt.Fprint(os.Stderr, `env:
  IMAGE_UPLOADER: imgur(default), s3
  ALLOW_PUBLIC_UPLOADER: set it to enable public uploader

env for imgur:
  IMGUR_CLIENT_ID: custom-client-id

env for s3:
  AWS_DEFAULT_REGION
  AWS_ACCESS_KEY_ID
  AWS_SECRET_ACCESS_KEY
  AWS_S3_BUCKET

return: image url
`)
		os.Exit(1)
	}
	filename := os.Args[1]

	var ut = uploader.DefaultUploader
	if t, ok := os.LookupEnv("IMAGE_UPLOADER"); ok {
		var err error
		ut, err = uploader.NewUploaderType(t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to parse IMAGE_UPLOADER: %v\n", err)
			os.Exit(1)
		}
	}

	if ut.Public() {
		if _, ok := os.LookupEnv("ALLOW_PUBLIC_UPLOADER"); !ok {
			fmt.Fprintf(os.Stderr, "error: public uploader is not allowed, set ALLOW_PUBLIC_UPLOADER to enable\n")
			os.Exit(1)
		}
	}

	u := ut.Uploader()
	if err := u.Upload(context.Background(), filename); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to upload: %v\n", err)
		os.Exit(1)
	}
}
