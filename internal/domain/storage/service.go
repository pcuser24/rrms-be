package storage

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/storage/dto"
)

type PresignURL struct {
	URL          string        `json:"url"`
	Method       string        `json:"method"`
	SignedHeader http.Header   `json:"signedHeader"`
	LifeTime     time.Duration `json:"lifeTime"`
}

type StorageService interface {
	GetPutObjectPresignURL(object *dto.PutObjectPresign, lifeTime time.Duration) (*PresignURL, error)
}

type Service interface {
	GetPresignUrl(data *dto.PutObjectPresign, userID uuid.UUID) (*PresignURL, error)
}

type service struct {
	storageService StorageService
}

func NewService(storageService StorageService) *service {
	return &service{
		storageService: storageService,
	}
}

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

func (s *service) GetPresignUrl(data *dto.PutObjectPresign, userID uuid.UUID) (*PresignURL, error) {
	data.Name = fmt.Sprintf("%s/%s", userID, processFilename(data.Name))
	return s.storageService.GetPutObjectPresignURL(data, 5*time.Minute)
}
