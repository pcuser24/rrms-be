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
	return s.domainRepo.RentalRepo.CreateContract(context.Background(), data)
}

func (s *service) GetRentalContractsOfUser(userId uuid.UUID, query *dto.GetRentalContracts) ([]model.ContractModel, error) {
	if query.Limit == nil {
		query.Limit = types.Ptr[int32](math.MaxInt32)
	}
	if query.Offset == nil {
		query.Offset = types.Ptr[int32](0)
	}
	rs, err := s.domainRepo.RentalRepo.GetRentalContractsOfUser(context.Background(), userId, query)
	if err != nil {
		return nil, err
	}

	return s.domainRepo.RentalRepo.GetContractsByIds(context.Background(), rs, query.Fields)
}

func (s *service) GetRentalContract(id int64) (*model.ContractModel, error) {
	return s.domainRepo.RentalRepo.GetContractByRentalID(context.Background(), id)
}

func (s *service) PingRentalContract(id int64) (any, error) {
	return s.domainRepo.RentalRepo.PingRentalContract(context.Background(), id)
}

func (s *service) GetContract(id int64) (*model.ContractModel, error) {
	return s.domainRepo.RentalRepo.GetContractByID(context.Background(), id)
}

func (s *service) UpdateContract(data *dto.UpdateContract) error {
	return s.domainRepo.RentalRepo.UpdateContract(context.Background(), data)
}

func (s *service) UpdateContractContent(data *dto.UpdateContractContent) error {
	return s.domainRepo.RentalRepo.UpdateContractContent(context.Background(), data)
}
