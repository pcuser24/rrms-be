package repo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func (r *repo) CreatePropertyVerificationRequest(ctx context.Context, data *property_dto.CreatePropertyVerificationRequest) (property_model.PropertyVerificationRequest, error) {
	res, err := r.dao.CreatePropertyVerificationRequest(ctx, data.ToCreatePropertyVerificationRequestParams())
	if err != nil {
		return property_model.PropertyVerificationRequest{}, err
	}
	return property_model.ToPropertyVerificationRequest(&res), nil
}

func (r *repo) GetPropertyVerificationRequest(ctx context.Context, id int64) (property_model.PropertyVerificationRequest, error) {
	res, err := r.dao.GetPropertyVerificationRequest(ctx, id)
	if err != nil {
		return property_model.PropertyVerificationRequest{}, err
	}
	return property_model.ToPropertyVerificationRequest(&res), nil
}

func (r *repo) GetPropertyVerificationRequests(ctx context.Context, query *property_dto.GetPropertyVerificationRequestsQuery) (*property_dto.GetPropertyVerificationRequestsResponse, error) {
	makeWhereExprs := func(sb *sqlbuilder.SelectBuilder) []string {
		var exprs []string = make([]string, 0)
		if query.CreatorID != uuid.Nil {
			exprs = append(exprs, sb.Equal("creator_id", query.CreatorID))
		}
		if query.PropertyID != uuid.Nil {
			exprs = append(exprs, sb.Equal("property_id", query.PropertyID))
		}
		if len(query.Status) > 0 {
			exprs = append(exprs, sb.In("status", sqlbuilder.List(query.Status)))
		}
		return exprs
	}

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("count(*) OVER() AS full_count")
	sb.From("property_verification_requests")
	sb.Where(makeWhereExprs(sb)...)

	sql, args := sb.Build()
	row := r.dao.QueryRow(ctx, sql, args...)
	var fullCount int64
	if err := row.Scan(&fullCount); err != nil {
		return nil, err
	}

	sb = sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "creator_id", "property_id", "video_url", "house_ownership_certificate", "certificate_of_landuse_right", "front_idcard", "back_idcard", "note", "status", "created_at", "updated_at")
	sb.From("property_verification_requests")
	sb.Where(makeWhereExprs(sb)...)
	if query.SortBy != "" {
		if query.Order == "asc" {
			sb.OrderBy(query.SortBy).Asc()
		} else {
			sb.OrderBy(query.SortBy).Desc()
		}
	}
	sb.Offset(int(query.Offset))
	sb.Limit(int(query.Limit))

	sql, args = sb.Build()
	rows, err := r.dao.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]property_model.PropertyVerificationRequest, 0, query.Limit)
	for rows.Next() {
		var i database.PropertyVerificationRequest
		if err := rows.Scan(&i.ID, &i.CreatorID, &i.PropertyID, &i.VideoUrl, &i.HouseOwnershipCertificate, &i.CertificateOfLanduseRight, &i.FrontIdcard, &i.BackIdcard, &i.Note, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, property_model.ToPropertyVerificationRequest(&i))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &property_dto.GetPropertyVerificationRequestsResponse{
		FullCount: fullCount,
		Items:     items,
	}, nil
}

func (r *repo) GetPropertyVerificationRequestsOfProperty(ctx context.Context, pid uuid.UUID, limit, offset int32) ([]property_model.PropertyVerificationRequest, error) {
	res, err := r.dao.GetPropertyVerificationRequestsOfProperty(ctx, database.GetPropertyVerificationRequestsOfPropertyParams{
		PropertyID: pid,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, err
	}
	var items []property_model.PropertyVerificationRequest
	for _, i := range res {
		items = append(items, property_model.ToPropertyVerificationRequest(&i))
	}
	return items, nil
}

func (r *repo) GetPropertiesVerificationStatus(ctx context.Context, ids []uuid.UUID) ([]property_dto.GetPropertyVerificationStatus, error) {
	buildSQL := func(pid uuid.UUID) (string, []interface{}) {
		sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
		sb.Select("property_id", "status")
		sb.From("property_verification_requests")
		sb.Where(sb.Equal("property_id", pid))
		sb.OrderBy("updated_at").Desc()
		sb.Limit(1)

		return sb.Build()
	}

	res := make([]property_dto.GetPropertyVerificationStatus, 0, len(ids))

	queries := make([]database.BatchedQueryRow, 0, len(ids))
	for _, id := range ids {
		sql, args := buildSQL(id)
		queries = append(queries, database.BatchedQueryRow{
			SQL:    sql,
			Params: args,
			Fn: func(row pgx.Row) error {
				var vs property_dto.GetPropertyVerificationStatus
				if err := row.Scan(&vs.PropertyID, &vs.Status); err != nil {
					if errors.Is(err, database.ErrRecordNotFound) {
						return nil
					}
					return err
				}
				res = append(res, vs)
				return nil
			},
		})
	}

	err := r.dao.QueryRowBatch(ctx, queries)

	return res, err
}

func (r *repo) UpdatePropertyVerificationRequestStatus(ctx context.Context, id int64, data *property_dto.UpdatePropertyVerificationRequestStatus) error {
	return r.dao.UpdatePropertyVerificationRequest(ctx, database.UpdatePropertyVerificationRequestParams{
		ID: id,
		Status: database.NullPROPERTYVERIFICATIONSTATUS{
			PROPERTYVERIFICATIONSTATUS: data.Status,
			Valid:                      data.Status != "",
		},
		Feedback: types.StrN(data.Feedback),
	})
}
