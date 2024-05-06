package service

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
)

func (s *service) CreateRentalComplaint(data *dto.CreateRentalComplaint) (model.RentalComplaint, error) {
	return s.rRepo.CreateRentalComplaint(context.Background(), data)
}

func (s *service) GetRentalComplaint(id int64) (model.RentalComplaint, error) {
	return s.rRepo.GetRentalComplaint(context.Background(), id)
}

func (s *service) GetRentalComplaintsByRentalId(rid int64) ([]model.RentalComplaint, error) {
	return s.rRepo.GetRentalComplaintsByRentalId(context.Background(), rid)
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

func (s *service) GetRentalComplaintReplies(id int64) ([]model.RentalComplaintReply, error) {
	return s.rRepo.GetRentalComplaintReplies(context.Background(), id)
}

func (s *service) UpdateRentalComplaint(data *dto.UpdateRentalComplaint) error {
	return s.rRepo.UpdateRentalComplaint(context.Background(), data)
}
