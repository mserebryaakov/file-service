package file

import (
	"context"

	"github.com/mserebryaakov/file-service/pkg/logger"
)

type service struct {
	storage Storage
	logger  *logger.Logger
}

func NewService(storage Storage, logger *logger.Logger) *service {
	return &service{
		storage: storage,
		logger:  logger,
	}
}

type Service interface {
	GetFile(ctx context.Context, bucketName, fileId string) (f *File, err error)
	GetFilesByBucketName(ctx context.Context, bucketName string) ([]*File, error)
	Create(ctx context.Context, bucketName string, dto CreateFileDTO) (string, error)
	Delete(ctx context.Context, bucketName, fileId string) error
}

func (s *service) GetFile(ctx context.Context, bucketName, fileId string) (f *File, err error) {
	f, err = s.storage.GetFile(ctx, bucketName, fileId)
	if err != nil {
		return f, err
	}
	return f, nil
}

func (s *service) GetFilesByBucketName(ctx context.Context, bucketName string) ([]*File, error) {
	files, err := s.storage.GetFilesByBucketName(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *service) Create(ctx context.Context, bucketName string, dto CreateFileDTO) (string, error) {
	dto.NormalizeName()
	file, err := NewFile(dto)
	if err != nil {
		return "", err
	}
	id, err := s.storage.CreateFile(ctx, bucketName, file)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *service) Delete(ctx context.Context, bucketName, fileName string) error {
	err := s.storage.DeleteFile(ctx, bucketName, fileName)
	if err != nil {
		return err
	}
	return nil
}
