package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"

	"github.com/google/uuid"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
)

var ErrPropertyVerificationRequestAlreadyExists = errors.New("property verification request already exists")

func (s *service) PreCreatePropertyVerificationRequest(data *property_dto.PreCreatePropertyVerificationRequest, creatorID uuid.UUID) error {
	_, err := s.domainRepo.PropertyRepo.GetPropertyVerificationRequests(context.Background(), &property_dto.GetPropertyVerificationRequestsQuery{
		PropertyID: data.PropertyID,
		Status:     []database.PROPERTYVERIFICATIONSTATUS{database.PROPERTYVERIFICATIONSTATUSPENDING},
		Limit:      1,
	})
	if err != nil && !errors.Is(err, database.ErrRecordNotFound) {
		return err
	}

	getPresignUrl := func(m *property_dto.Media) error {
		ext := filepath.Ext(m.Name)
		fname := m.Name[:len(m.Name)-len(ext)]
		// key = creatorID + "/" + "/property" + filename
		objKey := fmt.Sprintf("%s/property_verification/%s_%v%s", creatorID.String(), fname, time.Now().Unix(), ext)

		url, err := s.s3Client.GetPutObjectPresignedURL(
			s.imageBucketName, objKey, m.Type, m.Size, UPLOAD_URL_LIFETIME*time.Minute,
		)
		if err != nil {
			return err
		}
		m.Url = url.URL
		return nil
	}
	if data.HouseOwnershipCertificate != nil {
		if err := getPresignUrl(data.HouseOwnershipCertificate); err != nil {
			return err
		}
	}
	if data.CertificateOfLanduseRight != nil {
		if err := getPresignUrl(data.CertificateOfLanduseRight); err != nil {
			return err
		}
	}
	if err := getPresignUrl(&data.FrontIdcard); err != nil {
		return err
	}
	if err := getPresignUrl(&data.BackIdcard); err != nil {
		return err
	}
	return nil
}

func (s *service) CreatePropertyVerificationRequest(data *property_dto.CreatePropertyVerificationRequest) (property_model.PropertyVerificationRequest, error) {
	_, err := s.domainRepo.PropertyRepo.GetPropertyVerificationRequests(context.Background(), &property_dto.GetPropertyVerificationRequestsQuery{
		PropertyID: data.PropertyID,
		Status:     []database.PROPERTYVERIFICATIONSTATUS{database.PROPERTYVERIFICATIONSTATUSPENDING},
		Limit:      1,
	})
	if err != nil && !errors.Is(err, database.ErrRecordNotFound) {
		return property_model.PropertyVerificationRequest{}, err
	}

	return s.domainRepo.PropertyRepo.CreatePropertyVerificationRequest(context.Background(), data)
}

func (s *service) GetPropertyVerificationRequests(filter *property_dto.GetPropertyVerificationRequestsQuery) (*property_dto.GetPropertyVerificationRequestsResponse, error) {
	return s.domainRepo.PropertyRepo.GetPropertyVerificationRequests(context.Background(), filter)
}

func (s *service) GetPropertyVerificationRequest(id int64) (property_model.PropertyVerificationRequest, error) {
	return s.domainRepo.PropertyRepo.GetPropertyVerificationRequest(context.Background(), id)
}

func (s *service) GetPropertyVerificationRequestsOfProperty(pid uuid.UUID, limit, offset int32) ([]property_model.PropertyVerificationRequest, error) {
	return s.domainRepo.PropertyRepo.GetPropertyVerificationRequestsOfProperty(context.Background(), pid, limit, offset)
}

func (s *service) GetPropertiesVerificationStatus(ids []uuid.UUID) ([]property_dto.GetPropertyVerificationStatus, error) {
	return s.domainRepo.PropertyRepo.GetPropertiesVerificationStatus(context.Background(), ids)
}

func (s *service) UpdatePropertyVerificationRequestStatus(id int64, data *property_dto.UpdatePropertyVerificationRequestStatus) error {
	err := s.domainRepo.PropertyRepo.UpdatePropertyVerificationRequestStatus(context.Background(), id, data)
	if err != nil {
		return err
	}

	request, err := s.domainRepo.PropertyRepo.GetPropertyVerificationRequest(context.Background(), id)
	if err != nil {
		return err
	}

	// update es document
	doc := map[string]interface{}{
		"property.verification_status": data.Status,
	}
	docByte, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	client := s.esClient.GetTypedClient()
	_, err = client.Update(string(es.LISTINGINDEX), request.PropertyID.String()).
		Request(&update.Request{
			Doc: json.RawMessage(docByte),
		}).
		Do(context.Background())

	return err
}
