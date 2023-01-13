package writers

import (
	"bytes"
	"fmt"
	"log"
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

type s3Config struct {
	bucket string `default:""`
	key    string `default:""`
}

/* Parse a fully qualified S3 URI into a Bucket and Key */
func ParseUri(uri string) (s3Config, error) {
	config := s3Config{}

	parts := strings.Split(
		strings.ReplaceAll(uri, "s3://", ""),
		"/",
	)

	if !strings.HasPrefix(uri, "s3://") || len(parts) < 2 {
		return s3Config{}, fmt.Errorf("URI '%s' is not a valid S3 URI", uri)
	}

	config.bucket = parts[0]
	config.key = strings.Join(parts[1:], "/")

	return config, nil
}

/* Upload to S3 */
func S3Upload(uri string, data string) (*s3.PutObjectOutput, error) {
	config, err := ParseUri(uri)
	if err != nil {
		return nil, err
	}
	session := getSession()

	fmt.Println(config.key)
	res, err := s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.bucket),
		Key:    aws.String(config.key),
		Body:   aws.ReadSeekCloser(bytes.NewReader([]byte(data))),
	})

	return res, err
}
