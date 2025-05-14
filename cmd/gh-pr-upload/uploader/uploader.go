package uploader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

type UploaderType string

const (
	UploaderInvalid UploaderType = ""
	UploaderImgur   UploaderType = "imgur"
	UploaderS3      UploaderType = "s3"
	UploaderTest    UploaderType = "test"
)

const DefaultUploader = UploaderImgur

func NewUploaderType(typ string) (UploaderType, error) {
	switch UploaderType(typ) {
	case UploaderImgur, UploaderS3, UploaderTest:
		return UploaderType(typ), nil
	default:
		return UploaderInvalid, errors.New("unknown Uploader")
	}
}

func (u UploaderType) Public() bool {
	switch u {
	case UploaderImgur:
		return true
	case UploaderTest:
		return true
	default:
		return false
	}
}

func (u UploaderType) Uploader() Uploader {
	switch u {
	case UploaderImgur:
		return uploadImgur
	case UploaderS3:
		return uploadS3
	case UploaderTest:
		return uploadTest
	default:
		panic("invalid Uploader")
	}
}

type Uploader func(context.Context, string) error

func (u Uploader) Upload(ctx context.Context, filename string) error {
	return u(ctx, filename)
}

type imgurResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}

func uploadImgur(ctx context.Context, filename string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.imgur.com/3/image?type=file", nil)
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

func uploadS3(ctx context.Context, filename string) error {
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

	var usePathStyle bool
	if v, ok := os.LookupEnv("AWS_S3_USE_PATH_STYLE"); ok && v == "true" {
		usePathStyle = true
	}

	var objectKey string
	if name, err := uuid.NewRandom(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	} else {
		objectKey = name.String() + filepath.Ext(filename)
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load AWS config: %v\n", err)
		os.Exit(1)
	}
	s3Cli := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = usePathStyle
	})
	ul := s3manager.NewUploader(s3Cli)
	out, err := ul.Upload(ctx, &s3.PutObjectInput{
		ACL:         s3types.ObjectCannedACLPublicRead,
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

func uploadTest(ctx context.Context, filename string) error {
	fmt.Printf("file://%s", filename)

	return nil
}
