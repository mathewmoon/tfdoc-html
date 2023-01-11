package writers

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getSession() *session.Session {
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err)
	}
	return session
}

func ParseUri(uri string) (string, string) {
	if !strings.HasPrefix(uri, "s3://") {
		fmt.Printf("URI '%s' is not a valid S3 URI", uri)
		os.Exit(3)
	}
	uri = strings.ReplaceAll(uri, "s3://", "")
	parts := strings.Split(uri, "/")

	bucket := parts[0]
	key := strings.Join(parts[1:], "/")

	return bucket, key
}

func S3Upload(uri string, data string) (*s3.PutObjectOutput, error) {
	bucket, key := ParseUri(uri)

	session := getSession()

	res, err := s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(bytes.NewReader([]byte(data))),
	})

	return res, err
}
