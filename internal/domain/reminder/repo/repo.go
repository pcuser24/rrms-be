package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/reminder/dto"
	"github.com/user2410/rrms-backend/internal/domain/reminder/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreateReminder(ctx context.Context, data *dto.CreateReminder) (model.ReminderModel, error)
	GetRemindersOfUser(ctx context.Context, userId uuid.UUID, query *dto.GetRemindersQuery) ([]model.ReminderModel, error)
	GetReminder(ctx context.Context, id int64) (model.ReminderModel, error)
	CheckReminderVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error)
	UpdateReminder(ctx context.Context, data *dto.UpdateReminder) (int, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(d database.DAO) Repo {
	return &repo{
		dao: d,
	}
}

func (r *repo) CreateReminder(ctx context.Context, data *dto.CreateReminder) (model.ReminderModel, error) {
	rmdb, err := r.dao.CreateReminder(ctx, data.ToCreateReminderDB())
	if err != nil {
		return model.ReminderModel{}, err
	}
	rmm := model.ToReminderModel(&rmdb)

	for _, m := range data.Members {
		mdb, err := r.dao.CreateReminderMember(ctx, database.CreateReminderMemberParams{
			ReminderID: rmdb.ID,
			UserID:     m,
		})
		if err != nil {
			_ = r.dao.DeleteReminder(ctx, rmdb.ID)
			return model.ReminderModel{}, err
		}
		rmm.ReminderMembers = append(rmm.ReminderMembers, model.ReminderMemberModel{
			ReminderID: mdb.ReminderID,
			UserID:     mdb.UserID,
		})
	}

	return rmm, nil
}

func (r *repo) GetRemindersOfUser(ctx context.Context, userId uuid.UUID, query *dto.GetRemindersQuery) ([]model.ReminderModel, error) {
	var (
		res []database.Reminder
		err error
	)

	if query != nil && query.ResourceTag != nil {
		res, err = r.dao.GetRemindersOfUserWithResourceTag(ctx, database.GetRemindersOfUserWithResourceTagParams{
			UserID:      userId,
			ResourceTag: *query.ResourceTag,
		})
	} else {
		res, err = r.dao.GetRemindersOfUser(ctx, userId)
	}
	if err != nil {
		return nil, err
	}

	var reminders []model.ReminderModel
	for _, rm := range res {
		reminders = append(reminders, model.ToReminderModel(&rm))
	}

	return reminders, nil
}

func (r *repo) GetReminder(ctx context.Context, id int64) (model.ReminderModel, error) {
	res, err := r.dao.GetReminderById(ctx, id)
	if err != nil {
		return model.ReminderModel{}, err
	}
	return model.ToReminderModel(&res), nil
}

func (r *repo) UpdateReminder(ctx context.Context, data *dto.UpdateReminder) (int, error) {
	res, err := r.dao.UpdateReminder(ctx, data.ToUpdateReminderDB())
	if err != nil {
		return 0, err
	}
	return len(res), nil
}

func (r *repo) CheckReminderVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error) {
	return r.dao.CheckReminderVisibility(ctx, database.CheckReminderVisibilityParams{
		ReminderID: id,
		UserID:     userId,
	})
}
