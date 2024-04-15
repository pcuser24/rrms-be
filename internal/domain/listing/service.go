package listing

import (
	"context"
	"fmt"
	"net/url"

	"github.com/user2410/rrms-backend/internal/domain/listing/repo"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	listing_utils "github.com/user2410/rrms-backend/internal/domain/listing/utils"
	payment_dto "github.com/user2410/rrms-backend/internal/domain/payment/dto"
	payment_model "github.com/user2410/rrms-backend/internal/domain/payment/model"
	payment_repo "github.com/user2410/rrms-backend/internal/domain/payment/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/requests"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Service interface {
	CreateListing(data *dto.CreateListing) (*model.ListingModel, error)
	SearchListingCombination(data *dto.SearchListingCombinationQuery, userId uuid.UUID) (*dto.SearchListingCombinationResponse, error)
	GetListingByID(id uuid.UUID) (*model.ListingModel, error)
	GetListingsByIds(ids []uuid.UUID, fields []string) ([]model.ListingModel, error)
	GetListingsOfUser(userId uuid.UUID, fields []string) ([]model.ListingModel, error)
	UpdateListing(data *dto.UpdateListing) error
	DeleteListing(id uuid.UUID) error
	CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckListingVisibility(lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckValidUnitForListing(lid uuid.UUID, uid uuid.UUID) (bool, error)
	CreateListingPayment(data *dto.CreateListingPayment) (*payment_model.PaymentModel, error)
	CreateApplicationLink(data *dto.CreateApplicationLink) (string, error)
	VerifyApplicationLink(query *dto.VerifyApplicationLink) (bool, error)
}

type service struct {
	hashSecret  string
	lRepo       repo.Repo
	pRepo       property_repo.Repo
	paymentRepo payment_repo.Repo
}

func NewService(
	lRepo repo.Repo, pRepo property_repo.Repo, paymentRepo payment_repo.Repo,
	hashSecret string,
) Service {
	return &service{
		hashSecret:  hashSecret,
		lRepo:       lRepo,
		pRepo:       pRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *service) CreateListing(data *dto.CreateListing) (*model.ListingModel, error) {
	return s.lRepo.CreateListing(context.Background(), data)
}

func (s *service) SearchListingCombination(q *dto.SearchListingCombinationQuery, userId uuid.UUID) (*dto.SearchListingCombinationResponse, error) {
	if len(q.SortBy) == 0 {
		q.SortBy = append(q.SortBy, "listings.created_at", "listings.priority")
		q.Order = append(q.Order, "desc", "desc")
	}
	q.Limit = types.Ptr(utils.PtrDerefence(q.Limit, 1000))
	q.Offset = types.Ptr(utils.PtrDerefence(q.Offset, 0))
	res, err := s.lRepo.SearchListingCombination(context.Background(), q)
	if err != nil {
		return nil, err
	}
	var items []dto.SearchListingCombinationItem
	for _, item := range res.Items {
		isVisible, err := s.CheckListingVisibility(item.LId, userId)
		if err != nil {
			return nil, err
		}
		if isVisible {
			items = append(items, item)
		}
	}
	res.Items = items
	return res, nil
}

func (s *service) GetListingByID(id uuid.UUID) (*model.ListingModel, error) {
	return s.lRepo.GetListingByID(context.Background(), id)
}

func (s *service) GetListingsByIds(ids []uuid.UUID, fields []string) ([]model.ListingModel, error) {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = id.String()
	}
	return s.lRepo.GetListingsByIds(context.Background(), idsStr, fields)
}

func (s *service) UpdateListing(data *dto.UpdateListing) error {
	return s.lRepo.UpdateListing(context.Background(), data)
}

func (s *service) DeleteListing(id uuid.UUID) error {
	return s.lRepo.DeleteListing(context.Background(), id)
}

func (s *service) CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.lRepo.CheckListingOwnership(context.Background(), lid, uid)
}

func (s *service) CheckValidUnitForListing(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.lRepo.CheckValidUnitForListing(context.Background(), lid, uid)
}

func (s *service) GetListingsOfUser(userId uuid.UUID, fields []string) ([]model.ListingModel, error) {
	myListings, err := s.lRepo.SearchListingCombination(context.Background(), &dto.SearchListingCombinationQuery{
		SearchListingQuery: dto.SearchListingQuery{
			LCreatorID: types.Ptr(userId.String()),
		},
		SearchSortPaginationQuery: requests.SearchSortPaginationQuery{
			Limit:  types.Ptr[int32](1000),
			Offset: types.Ptr[int32](0),
			// SortBy: types.Ptr[string]("created_at"),
			// Order:  types.Ptr[string]("desc"),
		},
	})
	if err != nil {
		return nil, err
	}

	var lids []string
	for _, listing := range myListings.Items {
		lids = append(lids, listing.LId.String())
	}

	return s.lRepo.GetListingsByIds(context.Background(), lids, fields)

}

func (s *service) CreateListingPayment(data *dto.CreateListingPayment) (*payment_model.PaymentModel, error) {
	// create a payment entry
	params := payment_dto.CreatePayment{UserId: data.UserId}
	amount, err := listing_utils.CalculateListingPrice(data.Priority, data.PostDuration)
	if err != nil {
		return nil, err
	}
	params.Amount = amount

	params.OrderInfo = fmt.Sprintf("[CREATELISTING-%s] Phi dang tin nha cho thue", data.ListingId.String())
	params.Items = []payment_dto.CreatePaymentItem{
		{
			Name:     "Phi dang tin",
			Price:    amount,
			Quantity: 1,
			Discount: 0,
		},
	}
	return s.paymentRepo.CreatePayment(context.Background(), &params)
}

func (s *service) CreateApplicationLink(data *dto.CreateApplicationLink) (string, error) {
	key, err := listing_utils.EncryptApplicationLink(s.hashSecret, data)
	if err != nil {
		return "", err
	}

	urlValues := url.Values{}
	urlValues.Add("listingId", data.ListingId.String())
	urlValues.Add("fullName", data.FullName)
	urlValues.Add("email", data.Email)
	urlValues.Add("phone", data.Phone)
	urlValues.Add("k", key)

	return urlValues.Encode(), nil
}

func (s *service) VerifyApplicationLink(query *dto.VerifyApplicationLink) (bool, error) {
	return listing_utils.VerifyApplicationLink(query, s.hashSecret)
}

func (s *service) CheckListingVisibility(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.lRepo.CheckListingVisibility(context.Background(), lid, uid)
}
