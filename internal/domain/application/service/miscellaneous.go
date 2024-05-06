package service

import (
	"context"

	"github.com/user2410/rrms-backend/pkg/ds/set"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	reminder_dto "github.com/user2410/rrms-backend/internal/domain/reminder/dto"
	reminder_model "github.com/user2410/rrms-backend/internal/domain/reminder/model"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
)

func (s *service) CreateReminder(aid int64, userId uuid.UUID, data *dto.CreateReminder) (reminder_model.ReminderModel, error) {
	// get members of the application
	members := set.NewSet[uuid.UUID]()
	application, err := s.aRepo.GetApplicationById(context.Background(), aid)
	if err != nil {
		return reminder_model.ReminderModel{}, err
	}
	if application.CreatorID != uuid.Nil { // not an anonymous applicant
		members.Add(application.CreatorID)
	}
	if application.CreatorID == userId { // current user is the applicant
		pManagers, err := s.pRepo.GetPropertyManagers(context.Background(), application.PropertyID)
		if err != nil {
			return reminder_model.ReminderModel{}, err
		}
		for _, m := range pManagers {
			members.Add(m.ManagerID)
		}
	} else { // current user is a manager
		members.Add(userId)
	}
	for m := range members {
		data.Members = append(data.Members, m)
	}

	res, err := s.rService.CreateReminder(&reminder_dto.CreateReminder{
		CreatorID:      userId,
		Title:          data.Title,
		StartAt:        data.StartAt,
		EndAt:          data.EndAt,
		Note:           data.Note,
		Location:       data.Location,
		Priority:       data.Priority,
		RecurrenceMode: data.RecurrenceMode,
		ResourceTag:    GetResourceName(aid),
		Members:        members.ToSlice(),
	})
	if err != nil {
		return reminder_model.ReminderModel{}, err
	}

	return res, err
}

func (s *service) GetRentalByApplicationId(aid int64) (rental_model.RentalModel, error) {
	return s.aRepo.GetRentalByApplicationId(context.Background(), aid)
}
