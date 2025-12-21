package services

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type s3Client struct {
	client     *minio.Client
	bucketName string
}

func NewS3Client(client *minio.Client, bucketName string) S3Client {
	return &s3Client{
		client:     client,
		bucketName: bucketName,
	}
}

func (s *s3Client) GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	if bucketName == "" {
		bucketName = s.bucketName
	}

	obj, err := s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}
