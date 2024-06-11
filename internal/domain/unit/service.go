package unit

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/unit/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/unit/dto"
	"github.com/user2410/rrms-backend/internal/domain/unit/model"
)

const (
	MAX_IMAGE_SIZE      = 10 * 1024 * 1024 // 10MB
	UPLOAD_URL_LIFETIME = 5                // 5 minutes
)

type Service interface {
	PreCreateUnit(data *dto.PreCreateUnit, creatorID uuid.UUID) error
	CreateUnit(data *dto.CreateUnit) (*model.UnitModel, error)
	GetUnitById(id uuid.UUID) (*model.UnitModel, error)
	GetUnitsByIds(ids []uuid.UUID, fields []string, userId uuid.UUID) ([]model.UnitModel, error)
	SearchUnit(query *dto.SearchUnitCombinationQuery) (*dto.SearchUnitCombinationResponse, error)
	UpdateUnit(data *dto.UpdateUnit) error
	DeleteUnit(id uuid.UUID) error
	CheckVisibility(id uuid.UUID, uid uuid.UUID) (bool, error)
	CheckUnitManageability(id uuid.UUID, userId uuid.UUID) (bool, error)
	CheckUnitOfProperty(pid, uid uuid.UUID) (bool, error)
}

type service struct {
	repo repo.Repo

	s3Client        s3.S3Client
	imageBucketName string
}

func NewService(repo repo.Repo, s3Client s3.S3Client, imageBucketName string) Service {
	return &service{
		repo: repo,

		s3Client:        s3Client,
		imageBucketName: imageBucketName,
	}
}

func (s *service) PreCreateUnit(data *dto.PreCreateUnit, creatorID uuid.UUID) error {
	for i := range data.Media {
		m := &data.Media[i]
		// split file name and extension
		ext := filepath.Ext(m.Name)
		fname := m.Name[:len(m.Name)-len(ext)]
		// key = creatorID + "/" + "/property" + filename
		objKey := fmt.Sprintf("%s/units/%s_%v%s", creatorID.String(), fname, time.Now().Unix(), ext)

		url, err := s.s3Client.GetPutObjectPresignedURL(
			s.imageBucketName, objKey, m.Type, m.Size, UPLOAD_URL_LIFETIME*time.Minute,
		)
		if err != nil {
			return err
		}
		m.Url = url.URL
	}
	return nil
}

func (s *service) CreateUnit(data *dto.CreateUnit) (*model.UnitModel, error) {
	return s.repo.CreateUnit(context.Background(), data)
}

func (s *service) GetUnitById(id uuid.UUID) (*model.UnitModel, error) {
	return s.repo.GetUnitById(context.Background(), id)
}

func (s *service) GetUnitsByIds(ids []uuid.UUID, fields []string, userId uuid.UUID) ([]model.UnitModel, error) {
	var _ids []uuid.UUID
	for _, id := range ids {
		isVisible, err := s.CheckVisibility(id, userId)
		if err != nil {
			return nil, err
		}
		if isVisible {
			_ids = append(_ids, id)
		}
	}

	return s.repo.GetUnitsByIds(context.Background(), _ids, fields)
}

func (s *service) UpdateUnit(data *dto.UpdateUnit) error {
	return s.repo.UpdateUnit(context.Background(), data)
}

func (s *service) DeleteUnit(id uuid.UUID) error {
	return s.repo.DeleteUnit(context.Background(), id)
}

func (s *service) CheckUnitManageability(id uuid.UUID, userId uuid.UUID) (bool, error) {
	return s.repo.CheckUnitManageability(context.Background(), id, userId)
}

func (s *service) CheckVisibility(id uuid.UUID, uid uuid.UUID) (bool, error) {
	isPublic, err := s.repo.IsPublic(context.Background(), id)
	if err != nil {
		return false, err
	}
	if isPublic {
		return true, nil
	}
	return s.CheckUnitManageability(id, uid)
}

func (s *service) CheckUnitOfProperty(pid, uid uuid.UUID) (bool, error) {
	return s.repo.CheckUnitOfProperty(context.Background(), pid, uid)
}

func (s *service) SearchUnit(query *dto.SearchUnitCombinationQuery) (*dto.SearchUnitCombinationResponse, error) {
	return s.repo.SearchUnitCombination(context.Background(), query)
}
