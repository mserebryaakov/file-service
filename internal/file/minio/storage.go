package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/mserebryaakov/file-service/internal/file"
	"github.com/mserebryaakov/file-service/pkg/logger"
	"github.com/mserebryaakov/file-service/pkg/minio"
)

type minioStorage struct {
	client *minio.Client
	logger *logger.Logger
}

func NewStorage(endpoint, accessKeyID, secretAccessKey string, useSSL bool, logger *logger.Logger) (file.Storage, error) {
	client, err := minio.NewClient(endpoint, accessKeyID, secretAccessKey, useSSL, *logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}
	return &minioStorage{
		client: client,
		logger: logger,
	}, nil
}

func (m *minioStorage) GetFile(ctx context.Context, bucketName, fileID string) (*file.File, error) {
	obj, err := m.client.GetFile(ctx, bucketName, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file. err: %w", err)
	}
	defer obj.Close()
	objectInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file. err: %w", err)
	}
	buffer := make([]byte, objectInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}
	f := file.File{
		ID:    objectInfo.Key,
		Name:  objectInfo.UserMetadata["Name"],
		Size:  objectInfo.Size,
		Bytes: buffer,
	}
	return &f, nil
}

func (m *minioStorage) GetFilesByBucketName(ctx context.Context, bucketName string) ([]*file.File, error) {
	objects, err := m.client.GetBucketFiles(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}
	if len(objects) == 0 {
		return nil, fmt.Errorf("Not found")
	}

	var files []*file.File
	for _, obj := range objects {
		stat, err := obj.Stat()
		if err != nil {
			m.logger.Errorf("failed to get objects. err: %v", err)
			continue
		}
		buffer := make([]byte, stat.Size)
		_, err = obj.Read(buffer)
		if err != nil && err != io.EOF {
			m.logger.Errorf("failed to get objects. err: %v", err)
			continue
		}
		f := file.File{
			ID:    stat.Key,
			Name:  stat.UserMetadata["Name"],
			Size:  stat.Size,
			Bytes: buffer,
		}
		files = append(files, &f)
		obj.Close()
	}

	return files, nil
}

func (m *minioStorage) CreateFile(ctx context.Context, bucketName string, file *file.File) error {
	err := m.client.UploadFile(ctx, file.ID, file.Name, bucketName, file.Size, bytes.NewBuffer(file.Bytes))
	if err != nil {
		return err
	}
	return nil
}

func (m *minioStorage) DeleteFile(ctx context.Context, bucketName, fileId string) error {
	err := m.client.DeleteFile(ctx, bucketName, fileId)
	if err != nil {
		return err
	}
	return nil
}
