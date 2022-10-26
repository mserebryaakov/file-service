package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mserebryaakov/file-service/pkg/logger"
)

type Client struct {
	log         logger.Logger
	minioClient *minio.Client
}

// Инициализация MinIO Client
func NewClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool, logger logger.Logger) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}

	return &Client{
		log:         logger,
		minioClient: minioClient,
	}, nil
}

// Получение файла по fileId из bucketName
func (c *Client) GetFile(ctx context.Context, bucketName, fileId string) (*minio.Object, error) {
	reader, err := c.minioClient.GetObject(ctx, bucketName, fileId, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file with id: %s from minio bucket %s. err: %w", fileId, bucketName, err)
	}

	return reader, nil
}

// Получение всех файлов из bucketName
func (c *Client) GetBucketFiles(ctx context.Context, bucketName string) ([]*minio.Object, error) {
	var files []*minio.Object
	for lobj := range c.minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{WithMetadata: true}) {
		if lobj.Err != nil {
			c.log.Errorf("failed to list object from minio bucket %s. err: %v", bucketName, lobj.Err)
			continue
		}
		object, err := c.minioClient.GetObject(ctx, bucketName, lobj.Key, minio.GetObjectOptions{})
		if err != nil {
			c.log.Errorf("failed to get object key=%s from minio bucket %s. err: %v", lobj.Key, bucketName, lobj.Err)
			continue
		}
		files = append(files, object)
	}
	return files, nil
}

// Загрузка файла в bucketName (если backet не найдена - ошибка)
func (c *Client) UploadFile(ctx context.Context, fileId, fileName, bucketName string, fileSize int64, reader io.Reader) (string, error) {
	found, err := c.minioClient.BucketExists(ctx, bucketName)
	if err != nil || !found {
		err = c.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", fmt.Errorf("Failed to create Bucket with err: %w", err)
		}
	}

	c.log.Debugf("Put new object %s to bucket %s", fileName, bucketName)
	info, err := c.minioClient.PutObject(ctx, bucketName, fileId, reader, fileSize,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": fileName,
			},
			ContentType: "application/octet-stream",
		})
	if err != nil {
		return "", fmt.Errorf("Failed to upload file. err: %w", err)
	}

	return info.Key, nil
}

// Удаление файла
func (c *Client) DeleteFile(ctx context.Context, noteUUID, fileId string) error {
	err := c.minioClient.RemoveObject(ctx, noteUUID, fileId, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file. err: %w", err)
	}
	return nil
}
