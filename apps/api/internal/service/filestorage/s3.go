package filestorage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Storage struct {
	Client     *s3.Client
	BucketName string
	Region     string
}

func NewS3Storage(client *s3.Client, bucketName, region string) *S3Storage {
	return &S3Storage{
		Client:     client,
		BucketName: bucketName,
		Region:     region,
	}
}

func (s *S3Storage) Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	key := fmt.Sprintf("uploads/%d_%s_%s", time.Now().UnixNano(), uuid.New().String(), filename)

    // Only set ContentType to avoid "not open" errors if file needs to be seeked (though Reader doesn't support it)
    // S3 PutObject input expects io.Reader which is fine
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

    // Construct URL
    // Assuming standard S3 URL structure
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.BucketName, s.Region, key)
	return url, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	return err
}
