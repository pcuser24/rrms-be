package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/storage/dto"
	"github.com/user2410/rrms-backend/internal/domain/storage/object"
)

type Service interface {
	// GetPresignUrl returns a presigned URL for a specific client to upload an object to the storage
	GetPresignUrl(data *dto.PutObjectPresignRequest, userID uuid.UUID) (*object.PresignURL, error)
}

type service struct {
	storage Storage
}

func NewService(s Storage) *service {
	return &service{
		storage: s,
	}
}

// processFilename processes the filename to remove spaces and dots and add a timestamp
func processFilename(fileName string) string {
	// remove extension
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		fileName = fileName[:pos]
	}
	// remove spaces
	fileName = strings.Replace(fileName, " ", "_", -1)
	// remove dots
	fileName = strings.Replace(fileName, ".", "-", -1)
	// add timestamp
	return fmt.Sprintf("%s-%d", fileName, time.Now().UnixMilli())
}

func (s *service) GetPresignUrl(data *dto.PutObjectPresignRequest, userID uuid.UUID) (*object.PresignURL, error) {
	data.Name = fmt.Sprintf("%s/%s", userID, processFilename(data.Name))
	return s.storage.GetPutObjectPresignURL(data, time.Minute) // default lifetime of 1 minutes
}
