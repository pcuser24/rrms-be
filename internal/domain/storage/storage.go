package storage

import (
	"time"

	"github.com/user2410/rrms-backend/internal/domain/storage/dto"
	"github.com/user2410/rrms-backend/internal/domain/storage/object"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
)

type Storage interface {
	GetPutObjectPresignURL(object *dto.PutObjectPresignRequest, lifeTime time.Duration) (*object.PresignURL, error)
}

type s3Storage struct {
	storage       *s3.S3Client
	imgBucketName string
}

func NewStorage(
	s *s3.S3Client,
	imgBucketName string,
) *s3Storage {
	return &s3Storage{
		storage:       s,
		imgBucketName: imgBucketName,
	}
}

func (s *s3Storage) GetPutObjectPresignURL(
	o *dto.PutObjectPresignRequest, lifeTime time.Duration) (*object.PresignURL, error) {
	presignedRequest, err := s.storage.GetPutObjectPresignedURL(
		s.imgBucketName,
		o.Name, o.Type, o.Size,
		lifeTime,
	)
	if err != nil {
		return nil, err
	}
	return &object.PresignURL{
		URL:          presignedRequest.URL,
		Method:       presignedRequest.Method,
		SignedHeader: presignedRequest.SignedHeader,
		LifeTime:     lifeTime,
	}, nil
}
