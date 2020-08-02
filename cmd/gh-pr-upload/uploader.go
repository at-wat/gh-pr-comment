package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

type uploaderType string

const (
	uploaderInvalid uploaderType = ""
	uploaderImgur   uploaderType = "imgur"
	uploaderS3      uploaderType = "s3"
	uploaderTest    uploaderType = "test"
)

const defaultUploader = uploaderImgur

func newUploaderType(typ string) (uploaderType, error) {
	switch uploaderType(typ) {
	case uploaderImgur, uploaderS3, uploaderTest:
		return uploaderType(typ), nil
	default:
		return uploaderInvalid, errors.New("unknown uploader")
	}
}

func (u uploaderType) Public() bool {
	switch u {
	case uploaderImgur:
		return true
	case uploaderTest:
		return true
	default:
		return false
	}
}

func (u uploaderType) Uploader() uploader {
	switch u {
	case uploaderImgur:
		return uploadImgur
	case uploaderS3:
		return uploadS3
	case uploaderTest:
		return uploadTest
	default:
		panic("invalid uploader")
	}
}

type uploader func(string) error

func (u uploader) Upload(filename string) error {
	return u(filename)
}

type imgurResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}

func uploadImgur(filename string) error {
	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image?type=file", nil)
	if err != nil {
		return err
	}

	var clientID = "dd2b80c72f01f10"
	if id, ok := os.LookupEnv("IMGUR_CLIENT_ID"); ok {
		clientID = id
	}
	req.Header.Add("Authorization", "Client-ID "+clientID)

	req.Body, err = os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer req.Body.Close()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to post: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		fmt.Fprintf(os.Stderr, "error: failed to upload: %s\n", resp.Status)
		os.Exit(1)
	}

	respDec := json.NewDecoder(resp.Body)
	respJSON := &imgurResponse{}
	if err := respDec.Decode(respJSON); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to parse response: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(respJSON.Data.Link)

	return nil
}

func uploadS3(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to detect mime type: %v\n", err)
		os.Exit(1)
	}
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to seek file: %v\n", err)
		os.Exit(1)
	}

	var bucketName string
	if bucket, ok := os.LookupEnv("AWS_S3_BUCKET"); !ok {
		fmt.Fprintf(os.Stderr, "error: AWS_S3_BUCKET is not specified\n")
		os.Exit(1)
	} else {
		bucketName = bucket
	}

	var region string
	if r, ok := os.LookupEnv("AWS_REGION"); ok {
		region = r
	} else if r, ok := os.LookupEnv("AWS_DEFAULT_REGION"); ok {
		region = r
	}

	var objectKey string
	if name, err := uuid.NewRandom(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	} else {
		objectKey = name.String() + filepath.Ext(filename)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	ul := s3manager.NewUploader(sess)
	out, err := ul.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String(mimeType.String()),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to upload: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(out.Location)

	return nil
}

func uploadTest(filename string) error {
	req, err := http.NewRequest("POST", "http://devnull-as-a-service.com/dev/null", nil)
	if err != nil {
		return err
	}

	req.Body, err = os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer req.Body.Close()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to post: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		fmt.Fprintf(os.Stderr, "error: failed to upload: %s\n", resp.Status)
		os.Exit(1)
	}

	fmt.Printf("file://%s", filename)

	return nil
}
