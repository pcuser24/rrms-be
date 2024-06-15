package reminder

import (
	"context"
	"errors"

	"github.com/google/uuid"
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	"github.com/user2410/rrms-backend/internal/domain/reminder/dto"
	"github.com/user2410/rrms-backend/internal/domain/reminder/model"
)

type Service interface {
	CreateReminder(data *dto.CreateReminder) (model.ReminderModel, error)
	GetRemindersOfUser(userId uuid.UUID, query *dto.GetRemindersQuery) ([]model.ReminderModel, error)
	GetReminderById(id int64) (model.ReminderModel, error)
	CheckReminderVisibility(id int64, userId uuid.UUID) (bool, error)
}

type service struct {
	domainRepo repos.DomainRepo
}

func NewService(
	domainRepo repos.DomainRepo,
) Service {
	return &service{
		domainRepo: domainRepo,
	}
}

var ErrOverlappingReminder = errors.New("overlapping reminder")

func (s *service) CreateReminder(data *dto.CreateReminder) (model.ReminderModel, error) {
	isOverlapping, err := s.domainRepo.ReminderRepo.CheckOverlappingReminder(context.Background(), data.CreatorID, data.StartAt, data.EndAt)
	if err != nil {
		return model.ReminderModel{}, err
	}
	if isOverlapping {
		return model.ReminderModel{}, ErrOverlappingReminder
	}

	return s.domainRepo.ReminderRepo.CreateReminder(context.Background(), data)
}

func (s *service) GetRemindersOfUser(userId uuid.UUID, query *dto.GetRemindersQuery) ([]model.ReminderModel, error) {
	return s.domainRepo.ReminderRepo.GetRemindersOfUser(context.Background(), userId, query)
}

func (s *service) GetReminderById(id int64) (model.ReminderModel, error) {
	return s.domainRepo.ReminderRepo.GetReminder(context.Background(), id)
}

func (s *service) CheckReminderVisibility(id int64, userId uuid.UUID) (bool, error) {
	return s.domainRepo.ReminderRepo.CheckReminderVisibility(context.Background(), id, userId)
}
