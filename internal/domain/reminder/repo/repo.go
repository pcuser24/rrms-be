package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/redis/go-redis/v9"
	"github.com/user2410/rrms-backend/internal/domain/reminder/dto"
	"github.com/user2410/rrms-backend/internal/domain/reminder/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/redisd"
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
	dao         database.DAO
	redisClient redisd.RedisClient
}

func NewRepo(d database.DAO, redisClient redisd.RedisClient) Repo {
	return &repo{
		dao:         d,
		redisClient: redisClient,
	}
}

func (r *repo) CreateReminder(ctx context.Context, data *dto.CreateReminder) (model.ReminderModel, error) {
	rmdb, err := r.dao.CreateReminder(ctx, data.ToCreateReminderDB())
	if err != nil {
		return model.ReminderModel{}, err
	}

	rm := model.ToReminderModel(&rmdb)

	// save reminder to cache
	r.saveReminderToCache(ctx, &rm)

	return rm, nil
}

func (r *repo) GetRemindersOfUser(ctx context.Context, userId uuid.UUID, query *dto.GetRemindersQuery) ([]model.ReminderModel, error) {
	// read from cache
	reminders, err := r.getRemindersByStartAtRange(ctx, query.MinStartAt, query.MaxStartAt)
	if err == nil && len(reminders) > 0 {
		return nil, err
	}

	var (
		andExprs []string              = make([]string, 0)
		res      []model.ReminderModel = make([]model.ReminderModel, 0)
	)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "creator_id", "title", "start_at", "end_at", "note", "location", "created_at", "updated_at")
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

	// save to cache
	for _, reminder := range res {
		r.saveReminderToCache(ctx, &reminder)
	}

	return res, nil
}

func (r *repo) GetReminder(ctx context.Context, id int64) (model.ReminderModel, error) {
	// read from cache
	reminder, err := r.readReminderFromCache(ctx, id)
	if err == nil && reminder != nil {
		return *reminder, nil
	}

	res, err := r.dao.GetReminderById(ctx, id)
	if err != nil {
		return model.ReminderModel{}, err
	}

	rm := model.ToReminderModel(&res)

	// save to cache
	r.saveReminderToCache(ctx, &rm)

	return rm, nil
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
		ID:        id,
		CreatorID: userId,
	})
}

func (r *repo) CheckOverlappingReminder(ctx context.Context, userID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	subSB := sqlbuilder.PostgreSQL.NewSelectBuilder()
	subSB.Select("1").From("reminders").Where(
		subSB.Equal("creator_id", userID),
		subSB.Or(
			subSB.And(
				subSB.GTE("start_at", startTime),
				subSB.LTE("start_at", endTime),
			),
			subSB.And(
				subSB.GTE("end_at", startTime),
				subSB.LTE("end_at", endTime),
			),
		),
	)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select(sb.Exists(subSB))

	sql, args := sb.Build()
	row := r.dao.QueryRow(ctx, sql, args...)
	var res bool
	err := row.Scan(&res)
	return res, err
}

func (r *repo) saveReminderToCache(ctx context.Context, reminder *model.ReminderModel) error {
	dataMap, err := json.Marshal(reminder)
	if err != nil {
		return err
	}
	expiration := reminder.EndAt.Sub(time.Now().AddDate(0, 0, 1))
	if expiration < 0 {
		return nil
	}
	setRes := r.redisClient.Set(ctx, fmt.Sprintf("reminder:%d", reminder.ID), dataMap, expiration)
	if err := setRes.Err(); err != nil {
		return err
	}

	if err := r.redisClient.ZAdd(ctx, "reminders:startAt", redis.Z{
		Score:  float64(reminder.StartAt.Unix()),
		Member: reminder.ID,
	}).Err(); err != nil {
		return err
	}

	if err := r.redisClient.ZAdd(ctx, "reminders:endAt", redis.Z{
		Score:  float64(reminder.EndAt.Unix()),
		Member: reminder.ID,
	}).Err(); err != nil {
		return err
	}

	return nil
}

func (r *repo) readReminderFromCache(ctx context.Context, id int64) (*model.ReminderModel, error) {
	cacheRes := r.redisClient.Get(ctx, fmt.Sprintf("reminder:%d", id))
	if err := cacheRes.Err(); err != nil {
		return nil, err
	}
	dataStr := cacheRes.Val()
	if dataStr == "" {
		return nil, nil
	}
	var p model.ReminderModel
	if err := json.Unmarshal([]byte(dataStr), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repo) getRemindersByStartAtRange(ctx context.Context, minStartAt, maxStartAt time.Time) ([]model.ReminderModel, error) {
	// Convert the time range to Unix timestamps
	minScore := minStartAt.Unix()
	maxScore := maxStartAt.Unix()

	// Query the sorted set for reminder IDs within the range
	ids, err := r.redisClient.ZRangeByScore(ctx, "reminders:startAt", &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", minScore),
		Max: fmt.Sprintf("%d", maxScore),
	}).Result()
	if err != nil {
		return nil, err
	}

	// Fetch reminder data from hashes
	var reminders []model.ReminderModel
	res := r.redisClient.MGet(ctx, ids...)
	if err := res.Err(); err != nil {
		return nil, err
	}
	for _, datum := range res.Val() {
		dataStr, ok := datum.(string)
		if !ok {
			continue
		}

		var reminder model.ReminderModel
		if err := json.Unmarshal([]byte(dataStr), &reminder); err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}
