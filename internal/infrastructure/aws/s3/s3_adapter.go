package s3

import (
	"os"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/storage"
	"github.com/user2410/rrms-backend/internal/domain/storage/dto"
)

type AWSS3StorageService struct {
	s3Client   *S3Client
	bucketName string
}

func NewAWSS3StorageService(region, keyID, secretKey, bucketName string) (*AWSS3StorageService, error) {
	os.Setenv("AWS_REGION", region)
	os.Setenv("AWS_ACCESS_KEY_ID", keyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", secretKey)

	s3Client, err := NewS3Client()
	if err != nil {
		return nil, err
	}

	return &AWSS3StorageService{
		s3Client:   s3Client,
		bucketName: bucketName,
	}, nil
}

func (s *AWSS3StorageService) GetPutObjectPresignURL(object *dto.PutObjectPresign, lifeTime time.Duration) (*storage.PresignURL, error) {
	presignedRequest, err := s.s3Client.PutObject(
		s.bucketName,
		object.Name, object.Type, object.Size,
		lifeTime,
	)
	if err != nil {
		return nil, err
	}
	return &storage.PresignURL{
		URL:          presignedRequest.URL,
		Method:       presignedRequest.Method,
		SignedHeader: presignedRequest.SignedHeader,
		LifeTime:     lifeTime,
	}, nil
}
