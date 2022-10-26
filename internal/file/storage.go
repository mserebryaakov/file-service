package file

import (
	"context"
)

type Storage interface {
	GetFile(ctx context.Context, bucketName, fileID string) (*File, error)
	GetFilesByBucketName(ctx context.Context, bucketName string) ([]*File, error)
	CreateFile(ctx context.Context, bucketName string, file *File) error
	DeleteFile(ctx context.Context, bucketName, fileId string) error
}
