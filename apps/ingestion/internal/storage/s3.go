package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/replay/platform/apps/ingestion/internal/config"
)

// S3Client uploads content-addressed payloads to MinIO/S3.
type S3Client struct {
	client *s3.Client
	bucket string
	prefix string
}

// NewS3Client builds a client from config.
func NewS3Client(cfg config.Config) (*S3Client, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...any) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{URL: cfg.S3Endpoint, HostnameImmutable: true}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown service")
	})
	awsCfg := aws.Config{
		Region:                      "us-east-1",
		Credentials:                 credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		EndpointResolverWithOptions: resolver,
	}
	return &S3Client{
		client: s3.NewFromConfig(awsCfg, func(o *s3.Options) { o.UsePathStyle = true }),
		bucket: cfg.S3Bucket,
		prefix: cfg.BucketPrefix,
	}, nil
}

// PayloadKey returns org/project/sha256 object key layout.
func PayloadKey(orgID, projectID, sha256 string) string {
	return fmt.Sprintf("%s/%s/%s/%s", "payloads", orgID, projectID, sha256)
}

// PutObject stores bytes at the computed key.
func (c *S3Client) PutObject(ctx context.Context, key string, body []byte) (string, error) {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("s3://%s/%s", c.bucket, key), nil
}

// GetObject reads an object for tests.
func (c *S3Client) GetObject(ctx context.Context, key string) ([]byte, error) {
	out, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}
