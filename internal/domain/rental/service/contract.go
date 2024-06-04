package service

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func (s *service) CreateContract(data *dto.CreateContract) (*model.ContractModel, error) {
	return s.rRepo.CreateContract(context.Background(), data)
}

func (s *service) GetRentalContractsOfUser(userId uuid.UUID, query *dto.GetRentalContracts) ([]model.ContractModel, error) {
	if query.Limit == nil {
		query.Limit = types.Ptr[int32](math.MaxInt32)
	}
	if query.Offset == nil {
		query.Offset = types.Ptr[int32](0)
	}
	rs, err := s.rRepo.GetRentalContractsOfUser(context.Background(), userId, query)
	if err != nil {
		return nil, err
	}

	return s.rRepo.GetContractsByIds(context.Background(), rs, query.Fields)
}

func (s *service) GetRentalContract(id int64) (*model.ContractModel, error) {
	return s.rRepo.GetContractByRentalID(context.Background(), id)
}

func (s *service) PingRentalContract(id int64) (any, error) {
	return s.rRepo.PingRentalContract(context.Background(), id)
}

func (s *service) GetContract(id int64) (*model.ContractModel, error) {
	return s.rRepo.GetContractByID(context.Background(), id)
}

func (s *service) UpdateContract(data *dto.UpdateContract) error {
	return s.rRepo.UpdateContract(context.Background(), data)
}

func (s *service) UpdateContractContent(data *dto.UpdateContractContent) error {
	return s.rRepo.UpdateContractContent(context.Background(), data)
}
