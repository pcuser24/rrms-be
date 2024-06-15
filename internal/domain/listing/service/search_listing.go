package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func (s *service) SearchListingCombination(q *dto.SearchListingCombinationQuery, userId uuid.UUID) (*dto.SearchListingCombinationResponse, error) {
	if len(q.SortBy) == 0 {
		q.SortBy = append(q.SortBy, "listings.created_at", "listings.priority")
		q.Order = append(q.Order, "desc", "desc")
	}
	q.Limit = types.Ptr(utils.PtrDerefence(q.Limit, 1000))
	q.Offset = types.Ptr(utils.PtrDerefence(q.Offset, 0))
	q.LActive = types.Ptr(true)
	q.PIsPublic = types.Ptr(true)
	q.LMinExpiredAt = types.Ptr(time.Now())
	return s.domainRepo.ListingRepo.SearchListingCombination(context.Background(), q)
}
