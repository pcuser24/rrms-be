package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
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

func (s *service) CreateRentalComplaint(data *dto.CreateRentalComplaint) (model.RentalComplaint, error) {
	return s.rRepo.CreateRentalComplaint(context.Background(), data)
}

func (s *service) GetRentalComplaint(id int64) (model.RentalComplaint, error) {
	return s.rRepo.GetRentalComplaint(context.Background(), id)
}

func (s *service) GetRentalComplaintsByRentalId(rid int64) ([]model.RentalComplaint, error) {
	return s.rRepo.GetRentalComplaintsByRentalId(context.Background(), rid)
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

func (a *service) CreateRentalComplaintReply(data *dto.CreateRentalComplaintReply) (model.RentalComplaintReply, error) {
	res, err := a.rRepo.CreateRentalComplaintReply(context.Background(), data)
	if err != nil {
		return model.RentalComplaintReply{}, err
	}
	err = a.rRepo.UpdateRentalComplaint(context.Background(), &dto.UpdateRentalComplaint{
		ID:     data.ComplaintID,
		UserID: data.ReplierID,
	})
	return res, err
}

func (s *service) GetRentalComplaintReplies(rid int64, limit, offset int32) ([]model.RentalComplaintReply, error) {
	return s.rRepo.GetRentalComplaintReplies(context.Background(), rid, limit, offset)
}

func (s *service) GetRentalComplaintsOfUser(userId uuid.UUID, query dto.GetRentalComplaintsOfUserQuery) ([]model.RentalComplaint, error) {
	return s.rRepo.GetRentalComplaintsOfUser(context.Background(), userId, query)
}

func (s *service) UpdateRentalComplaint(data *dto.UpdateRentalComplaint) error {
	return s.rRepo.UpdateRentalComplaint(context.Background(), data)
}
