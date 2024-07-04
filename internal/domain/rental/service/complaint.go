package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

func (s *service) PreCreateRentalComplaint(data *dto.PreCreateRentalComplaint, creatorID uuid.UUID) error {
	for i := range data.Media {
		m := &data.Media[i]
		// split file name and extension
		ext := filepath.Ext(m.Name)
		fname := m.Name[:len(m.Name)-len(ext)]
		// key = creatorID + "/" + "/property" + filename
		objKey := fmt.Sprintf("%s/rental-complaints/%s_%v%s", creatorID.String(), fname, time.Now().Unix(), ext)

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

var ErrUnauthorizedToCreateComplaint = fmt.Errorf("unauthorized to create complaint")

func (s *service) CreateRentalComplaint(data *dto.CreateRentalComplaint) (model.RentalComplaint, error) {
	rental, err := s.domainRepo.RentalRepo.GetRental(context.Background(), data.RentalID)
	if err != nil {
		return model.RentalComplaint{}, err
	}

	canCreate := false
	if rental.TenantID == data.CreatorID {
		canCreate = true
	} else {
		managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), rental.PropertyID)
		if err != nil {
			return model.RentalComplaint{}, err
		}
		for _, m := range managers {
			if m.ManagerID == data.CreatorID {
				canCreate = true
				break
			}
		}
	}
	if !canCreate {
		return model.RentalComplaint{}, ErrUnauthorizedToCreateComplaint
	}

	res, err := s.domainRepo.RentalRepo.CreateRentalComplaint(context.Background(), data)
	if err != nil {
		return model.RentalComplaint{}, err
	}

	err = s.notifyCreateRentalComplaint(&res, &rental)
	if err != nil {
		// TODO: log the error
	}

	return res, nil
}

func (s *service) GetRentalComplaint(id int64) (model.RentalComplaint, error) {
	return s.domainRepo.RentalRepo.GetRentalComplaint(context.Background(), id)
}

func (s *service) GetRentalComplaintsByRentalId(rid int64, limit, offset int32) ([]model.RentalComplaint, error) {
	return s.domainRepo.RentalRepo.GetRentalComplaintsByRentalId(context.Background(), rid, limit, offset)
}

func (s *service) PreCreateRentalComplaintReply(data *dto.PreCreateRentalComplaint, creatorID uuid.UUID) error {
	for i := range data.Media {
		m := &data.Media[i]
		// split file name and extension
		ext := filepath.Ext(m.Name)
		fname := m.Name[:len(m.Name)-len(ext)]
		// key = creatorID + "/" + "/property" + filename
		objKey := fmt.Sprintf("%s/rental-complaints/%s_%v%s", creatorID.String(), fname, time.Now().Unix(), ext)

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

var (
	ErrUnauthorizedToCreateComplaintReply = fmt.Errorf("unauthorized to create complaint reply")
)

func (s *service) CreateRentalComplaintReply(data *dto.CreateRentalComplaintReply) (model.RentalComplaintReply, error) {
	complaint, err := s.domainRepo.RentalRepo.GetRentalComplaint(context.Background(), data.ComplaintID)
	if err != nil {
		return model.RentalComplaintReply{}, err
	}
	if complaint.Status != database.RENTALCOMPLAINTSTATUSPENDING {
		return model.RentalComplaintReply{}, ErrUnauthorizedToCreateComplaintReply
	}
	rental, err := s.domainRepo.RentalRepo.GetRental(context.Background(), complaint.RentalID)
	if err != nil {
		return model.RentalComplaintReply{}, err
	}

	res, err := s.domainRepo.RentalRepo.CreateRentalComplaintReply(context.Background(), data)
	if err != nil {
		return model.RentalComplaintReply{}, err
	}
	err = s.domainRepo.RentalRepo.UpdateRentalComplaint(context.Background(), &dto.UpdateRentalComplaint{
		ID:     data.ComplaintID,
		UserID: data.ReplierID,
	})
	if err != nil {
		return res, err
	}

	err = s.notifyCreateComplaintReply(&complaint, &res, &rental)
	if err != nil {
		// TODO: log the error
	}

	return res, nil
}

func (s *service) GetRentalComplaintReplies(rid int64, limit, offset int32) ([]model.RentalComplaintReply, error) {
	return s.domainRepo.RentalRepo.GetRentalComplaintReplies(context.Background(), rid, limit, offset)
}

func (s *service) GetRentalComplaintsOfUser(userId uuid.UUID, query dto.GetRentalComplaintsOfUserQuery) ([]model.RentalComplaint, error) {
	return s.domainRepo.RentalRepo.GetRentalComplaintsOfUser(context.Background(), userId, query)
}

func (s *service) UpdateRentalComplaintStatus(data *dto.UpdateRentalComplaintStatus) error {
	complaint, err := s.domainRepo.RentalRepo.GetRentalComplaint(context.Background(), data.ID)
	if err != nil {
		return err
	}
	rental, err := s.domainRepo.RentalRepo.GetRental(context.Background(), complaint.RentalID)
	if err != nil {
		return err
	}

	err = s.domainRepo.RentalRepo.UpdateRentalComplaint(context.Background(), &dto.UpdateRentalComplaint{
		ID:     data.ID,
		Status: data.Status,
		UserID: data.UserID,
	})
	if err != nil {
		return err
	}

	err = s.notifyUpdateComplaintStatus(&complaint, &rental, data.Status, data.UserID)
	if err != nil {
		// TODO: log the error
	}

	return nil
}
