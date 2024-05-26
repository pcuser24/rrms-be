package service

import (
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
)

type SuggestedListing struct {
}

func (s *service) GetListingSuggestions(data *dto.SearchListingCombinationQuery) ([]SuggestedListing, error) {
	return nil, nil
}
