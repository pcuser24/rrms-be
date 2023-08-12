package rental

import (
	"context"

	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	GetAllRentalPolicies(ctx context.Context) ([]model.RentalPolicyModel, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(dao database.DAO) Repo {
	return &repo{
		dao: dao,
	}
}

func (r *repo) GetAllRentalPolicies(ctx context.Context) ([]model.RentalPolicyModel, error) {
	resDB, err := r.dao.GetAllRentalPolicies(ctx)
	if err != nil {
		return nil, err
	}

	var res []model.RentalPolicyModel
	for _, i := range resDB {
		res = append(res, model.RentalPolicyModel(i))
	}

	return res, nil
}
