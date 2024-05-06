package repo

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
)

func (r *repo) CreateRentalComplaint(ctx context.Context, data *dto.CreateRentalComplaint) (model.RentalComplaint, error) {
	res, err := r.dao.CreateRentalComplaint(ctx, data.ToCreateRentalComplaintDB())
	if err != nil {
		return model.RentalComplaint{}, err
	}
	return model.ToRentalComplaintModel(&res), nil
}

func (r *repo) GetRentalComplaint(ctx context.Context, id int64) (model.RentalComplaint, error) {
	res, err := r.dao.GetRentalComplaint(ctx, id)
	if err != nil {
		return model.RentalComplaint{}, err
	}
	return model.ToRentalComplaintModel(&res), nil
}

func (r *repo) GetRentalComplaintsByRentalId(ctx context.Context, rid int64) ([]model.RentalComplaint, error) {
	res, err := r.dao.GetRentalComplaintsByRentalId(ctx, rid)
	if err != nil {
		return nil, err
	}
	var result []model.RentalComplaint
	for _, v := range res {
		result = append(result, model.ToRentalComplaintModel(&v))
	}
	return result, nil
}

func (r *repo) CreateRentalComplaintReply(ctx context.Context, data *dto.CreateRentalComplaintReply) (model.RentalComplaintReply, error) {
	res, err := r.dao.CreateRentalComplaintReply(ctx, data.ToCreateRentalComplaintReplyDB())
	if err != nil {
		return model.RentalComplaintReply{}, err
	}
	return model.RentalComplaintReply(res), nil
}

func (r *repo) GetRentalComplaintReplies(ctx context.Context, id int64) ([]model.RentalComplaintReply, error) {
	res, err := r.dao.GetRentalComplaintReplies(ctx, id)
	if err != nil {
		return nil, err
	}
	var result []model.RentalComplaintReply
	for _, v := range res {
		result = append(result, model.RentalComplaintReply(v))
	}
	return result, nil
}

func (r *repo) UpdateRentalComplaint(ctx context.Context, data *dto.UpdateRentalComplaint) error {
	return r.dao.UpdateRentalComplaint(ctx, data.ToUpdateRentalComplaintDB())
}
