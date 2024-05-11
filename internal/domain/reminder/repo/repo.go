package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/user2410/rrms-backend/internal/domain/reminder/dto"
	"github.com/user2410/rrms-backend/internal/domain/reminder/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreateReminder(ctx context.Context, data *dto.CreateReminder) (model.ReminderModel, error)
	GetRemindersOfUser(ctx context.Context, userId uuid.UUID, query *dto.GetRemindersQuery) ([]model.ReminderModel, error)
	GetReminder(ctx context.Context, id int64) (model.ReminderModel, error)
	CheckReminderVisibility(ctx context.Context, id int64, userId uuid.UUID) (bool, error)
	CheckOverlappingReminder(ctx context.Context, userID uuid.UUID, startTime, endTime time.Time) (bool, error)
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
		andExprs []string              = make([]string, 0)
		res      []model.ReminderModel = make([]model.ReminderModel, 0)
	)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "creator_id", "title", "start_at", "end_at", "note", "location", "recurrence_day", "recurrence_month", "recurrence_mode", "priority", "status", "resource_tag", "created_at", "updated_at")
	sb.From("reminders")
	if query.CreatorID != uuid.Nil {
		andExprs = append(andExprs, sb.Equal("creator_id", query.CreatorID))
	}
	if !query.MinStartAt.IsZero() {
		andExprs = append(andExprs, sb.GTE("start_at", query.MinStartAt))
	}
	if !query.MaxStartAt.IsZero() {
		andExprs = append(andExprs, sb.LTE("start_at", query.MaxStartAt))
	}
	if !query.MinEndAt.IsZero() {
		andExprs = append(andExprs, sb.GTE("end_at", query.MinEndAt))
	}
	if !query.MaxEndAt.IsZero() {
		andExprs = append(andExprs, sb.LTE("end_at", query.MaxEndAt))
	}
	if query.Priority != nil {
		andExprs = append(andExprs, sb.Equal("priority", *query.Priority))
	}
	if query.Status != "" {
		andExprs = append(andExprs, sb.Equal("status", query.Status))
	}
	if query.RecurrenceMode != "" {
		andExprs = append(andExprs, sb.Equal("recurrence_mode", query.RecurrenceMode))
	}
	if query.RecurrenceDay != nil {
		andExprs = append(andExprs, sb.Equal("recurrence_day", *query.RecurrenceDay))
	}
	if query.RecurrenceMonth != nil {
		andExprs = append(andExprs, sb.Equal("recurrence_month", *query.RecurrenceMonth))
	}
	if query.ResourceTag != nil {
		andExprs = append(andExprs, sb.Equal("resource_tag", *query.ResourceTag))
	}
	if query.Members != nil && len(query.Members) > 0 {
		subSB := sqlbuilder.PostgreSQL.NewSelectBuilder()
		andExprs = append(andExprs, sb.Exists(
			subSB.Select("1").
				From("reminder_members").
				Where(
					subSB.Equal("reminder_members.reminder_id", "reminders.id"),
					subSB.In("user_id", sqlbuilder.List(query.Members)),
				),
		))
	}
	if len(andExprs) > 0 {
		sb.Where(andExprs...)
	} else {
		sb.Where(sb.Equal("creator_id", userId))
	}

	sql, args := sb.Build()
	rows, err := r.dao.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i database.Reminder
		if err = rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.StartAt,
			&i.EndAt,
			&i.Note,
			&i.Location,
			&i.RecurrenceDay,
			&i.RecurrenceMonth,
			&i.RecurrenceMode,
			&i.Priority,
			&i.Status,
			&i.ResourceTag,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, model.ToReminderModel(&i))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
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

func (r *repo) CheckOverlappingReminder(ctx context.Context, userID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	return r.dao.CheckOverlappingReminder(ctx, database.CheckOverlappingReminderParams{
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime,
	})
}
