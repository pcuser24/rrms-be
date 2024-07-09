package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
	"github.com/user2410/rrms-backend/internal/utils/types"

	"github.com/google/uuid"
	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	listing_dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	misc_dto "github.com/user2410/rrms-backend/internal/domain/misc/dto"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	html_util "github.com/user2410/rrms-backend/internal/utils/template/html"
	text_util "github.com/user2410/rrms-backend/internal/utils/template/text"
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

	res, err := s.domainRepo.PropertyRepo.CreatePropertyVerificationRequest(context.Background(), data)
	if err != nil {
		return property_model.PropertyVerificationRequest{}, err
	}

	// send notification about the creation of verification request
	err = s.asynctaskDistributor.DistributeTaskJSON(context.Background(), asynctask.PROPERTY_VERIFICATION_CREATE, data)

	return res, err
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

	// update corresponding listings documents in ES
	doc := map[string]interface{}{
		"property": map[string]interface{}{
			"verification_status": data.Status,
		},
	}
	docByte, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	listingIds, err := s.domainRepo.PropertyRepo.GetListingsOfProperty(context.Background(), request.PropertyID, &listing_dto.GetListingsOfPropertyQuery{
		Limit:  types.Ptr[int32](math.MaxInt32),
		Offset: types.Ptr[int32](0),
	})
	if err != nil {
		return err
	}

	// TODO: batch update
	client := s.esClient.GetTypedClient()
	for _, listingId := range listingIds {
		_, err = client.Update(string(es.LISTINGINDEX), listingId.String()).
			Request(&update.Request{
				Doc: json.RawMessage(docByte),
			}).
			Do(context.Background())
		if err != nil {
			return err
		}
	}

	// send notification about the update of verification status
	err = s.asynctaskDistributor.DistributeTaskJSON(context.Background(), asynctask.PROPERTY_VERIFICATION_UPDATE, property_dto.UpdatePropertyVerificationRequestStatusNotification{
		Request:    &request,
		UpdateData: data,
	})

	return err
}

const basePath = "internal/domain/property/service/templates"

func (s *service) NotifyCreatePropertyVerificationRequestStatus(r *property_model.PropertyVerificationRequest) error {
	var (
		data = struct {
			FESite   string
			Property *property_model.PropertyModel
			Request  *property_model.PropertyVerificationRequest
			Creator  *auth_model.UserModel
		}{
			FESite:  s.feSite,
			Request: r,
		}
		err error
	)
	{
		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{r.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		data.Property = &ps[0]

		us, err := s.domainRepo.AuthRepo.GetUsersByIds(context.Background(), []uuid.UUID{r.CreatorID}, []string{"email"})
		if err != nil {
			return err
		}
		if len(us) == 0 {
			return database.ErrRecordNotFound
		}
		data.Creator = &us[0]
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/create_verification.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/create_verification.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}

	adminIds, err := s.domainRepo.AuthRepo.GetAdminUsers(context.Background())
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_CREATEPROPERTYVERIFICATIONSTATUS,
			"propertyId":       r.PropertyID.String(),
			"requestId":        r.ID,
		},
		Targets: make([]misc_dto.CreateNotificationTarget, len(adminIds)),
	}
	for i, adminId := range adminIds {
		ds, err := s.miscService.GetNotificationDevice(adminId, uuid.Nil, "", "")
		if err != nil {
			return err
		}
		tokens := make([]string, len(ds))
		for i, d := range ds {
			tokens[i] = d.Token
		}
		cn.Targets[i] = misc_dto.CreateNotificationTarget{
			UserId: adminId,
			Tokens: tokens,
		}
	}

	return s.miscService.SendNotification(&cn)
}

func (s *service) NotifyUpdatePropertyVerificationRequestStatus(
	request *property_model.PropertyVerificationRequest,
	updateData *property_dto.UpdatePropertyVerificationRequestStatus,
) error {
	var (
		data = struct {
			FESite     string
			Property   *property_model.PropertyModel
			Request    *property_model.PropertyVerificationRequest
			UpdateData *property_dto.UpdatePropertyVerificationRequestStatus
		}{
			FESite:     s.feSite,
			Request:    request,
			UpdateData: updateData,
		}
		err error
	)
	{
		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{request.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		data.Property = &ps[0]
	}
	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/update_verification.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/update_verification.gohtml", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/update_verification.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}

	// send email
	targets, err := s.miscService.GetNotificationManagersTargets(request.PropertyID)
	if err != nil {
		return err
	}
	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_UPDATEPROPERTYVERIFICATIONSTATUS,
			"propertyId":       request.PropertyID.String(),
			"requestId":        request.ID,
		},
		Targets: func() []misc_dto.CreateNotificationTarget {
			var ts []misc_dto.CreateNotificationTarget
			for _, t := range targets {
				ts = append(ts, misc_dto.CreateNotificationTarget{
					UserId: t.UserId,
					Emails: t.Emails,
				})
			}
			return ts
		}(),
	}
	err = s.miscService.SendNotification(&cn)
	if err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = func() []misc_dto.CreateNotificationTarget {
		var ts []misc_dto.CreateNotificationTarget
		for _, t := range targets {
			ts = append(ts, misc_dto.CreateNotificationTarget{
				UserId: t.UserId,
				Tokens: t.Tokens,
			})
		}
		return ts
	}()
	err = s.miscService.SendNotification(&cn)

	return err
}
