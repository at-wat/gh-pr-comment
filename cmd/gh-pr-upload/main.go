package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s filename\n", os.Args[0])
		os.Exit(1)
	}
	filename := os.Args[1]

	var ut = defaultUploader
	if t, ok := os.LookupEnv("IMAGE_UPLOADER"); ok {
		var err error
		ut, err = newUploaderType(t)
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
	if err := u.Upload(filename); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to upload: %v\n", err)
		os.Exit(1)
	}
}
