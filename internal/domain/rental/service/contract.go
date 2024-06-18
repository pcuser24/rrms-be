package service

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/rental/dto"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var ErrUnauthorizedToCreateContract = errors.New("unauthorized to create contract")

func (s *service) CreateContract(data *dto.CreateContract) (*model.ContractModel, error) {
	rental, err := s.domainRepo.RentalRepo.GetRental(context.Background(), data.RentalID)
	if err != nil {
		return nil, err
	}
	managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), rental.PropertyID)
	if err != nil {
		return nil, err
	}
	// check if the user is a manager of the property
	isManager := false
	for _, m := range managers {
		if m.ManagerID == data.UserID {
			isManager = true
			break
		}
	}
	if !isManager {
		return nil, ErrUnauthorizedToCreateContract
	}

	contract, err := s.domainRepo.RentalRepo.CreateContract(context.Background(), data)
	if err != nil {
		return nil, err
	}

	err = s.notifyCreateContract(contract, &rental)
	if err != nil {
		// TODO: log error
	}

	return contract, nil
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

var ErrUnauthorizedToUpdateContract = errors.New("unauthorized to update contract")

func (s *service) UpdateContract(data *dto.UpdateContract) error {
	cs, err := s.domainRepo.RentalRepo.GetContractsByIds(context.Background(), []int64{data.ID}, []string{"rental_id", "updated_by"})
	if err != nil {
		return err
	}
	if len(cs) == 0 {
		return database.ErrRecordNotFound
	}
	rental, err := s.domainRepo.RentalRepo.GetRental(context.Background(), cs[0].RentalID)
	if err != nil {
		return err
	}
	// check if the user is eligible to update the contract
	var (
		updaterSide, lastUpdaterSide string
		canUpdate                    bool = false
	)
	if cs[0].UpdatedBy == rental.TenantID {
		lastUpdaterSide = "B"
	} else {
		lastUpdaterSide = "A"
	}

	if rental.TenantID == data.UserID {
		canUpdate = true
		updaterSide = "B"
	} else {
		managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), rental.PropertyID)
		if err != nil {
			return err
		}
		for _, m := range managers {
			if m.ManagerID == data.UserID {
				canUpdate = true
				updaterSide = "A"
				break
			}
		}
	}
	if !canUpdate {
		return ErrUnauthorizedToUpdateContract
	}

	err = s.domainRepo.RentalRepo.UpdateContract(context.Background(), data)
	if err != nil {
		return err
	}

	if lastUpdaterSide != updaterSide {
		s.notifyUpdateContract(&cs[0], &rental, updaterSide)
	}

	return nil
}

func (s *service) UpdateContractContent(data *dto.UpdateContractContent) error {
	cs, err := s.domainRepo.RentalRepo.GetContractsByIds(context.Background(), []int64{data.ID}, []string{"rental_id", "updated_by"})
	if err != nil {
		return err
	}
	if len(cs) == 0 {
		return database.ErrRecordNotFound
	}
	rental, err := s.domainRepo.RentalRepo.GetRental(context.Background(), cs[0].RentalID)
	if err != nil {
		return err
	}
	// check if the user is eligible to update the contract
	var (
		updaterSide, lastUpdaterSide string
		canUpdate                    bool = false
	)
	if cs[0].UpdatedBy == rental.TenantID {
		lastUpdaterSide = "B"
	} else {
		lastUpdaterSide = "A"
	}

	if rental.TenantID == data.UserID {
		canUpdate = true
		updaterSide = "B"
	} else {
		managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), rental.PropertyID)
		if err != nil {
			return err
		}
		for _, m := range managers {
			if m.ManagerID == data.UserID {
				canUpdate = true
				updaterSide = "A"
				break
			}
		}
	}
	if !canUpdate {
		return ErrUnauthorizedToUpdateContract
	}

	err = s.domainRepo.RentalRepo.UpdateContractContent(context.Background(), data)
	if err != nil {
		return err
	}

	if lastUpdaterSide != updaterSide {
		s.notifyUpdateContract(&cs[0], &rental, updaterSide)
	}

	return nil
}
