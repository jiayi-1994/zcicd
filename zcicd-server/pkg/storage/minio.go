package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zcicd/zcicd-server/pkg/config"
)

type Client struct {
	client *minio.Client
}

func NewMinIOClient(cfg *config.Config) (*Client, error) {
	mc, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}
	return &Client{client: mc}, nil
}

func (c *Client) EnsureBucket(ctx context.Context, bucket string) error {
	exists, err := c.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		if err := c.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
		}
	}
	return nil
}

func (c *Client) UploadFile(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload %s: %w", objectName, err)
	}
	return nil
}

func (c *Client) DownloadFile(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	obj, err := c.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", objectName, err)
	}
	return obj, nil
}

func (c *Client) DeleteFile(ctx context.Context, bucket, objectName string) error {
	err := c.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", objectName, err)
	}
	return nil
}

func (c *Client) GetPresignedURL(ctx context.Context, bucket, objectName string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := c.client.PresignedGetObject(ctx, bucket, objectName, expires, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL for %s: %w", objectName, err)
	}
	return presignedURL.String(), nil
}
