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
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/requests"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Service interface {
	CreateListing(data *dto.CreateListing) (*model.ListingModel, error)
	SearchListingCombination(data *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error)
	GetListingByID(id uuid.UUID) (*model.ListingModel, error)
	GetListingsByIds(ids []uuid.UUID, fields []string) ([]model.ListingModel, error)
	GetListingsOfUser(userId uuid.UUID, fields []string) ([]model.ListingModel, error)
	UpdateListing(data *dto.UpdateListing) error
	DeleteListing(id uuid.UUID) error
	CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error)
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
	listing, err := s.lRepo.CreateListing(context.Background(), data)
	if err != nil {
		return nil, err
	}
	err = s.pRepo.UpdateProperty(context.Background(), &property_dto.UpdateProperty{
		ID:       listing.PropertyID,
		IsPublic: types.Ptr(true),
	})
	if err != nil {
		return listing, err
	}

	return listing, nil
}

func (s *service) SearchListingCombination(data *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error) {
	data.SortBy = types.Ptr(utils.PtrDerefence[string](data.SortBy, "created_at"))
	data.Order = types.Ptr(utils.PtrDerefence[string](data.Order, "desc"))
	data.Limit = types.Ptr(utils.PtrDerefence[int32](data.Limit, 1000))
	data.Offset = types.Ptr(utils.PtrDerefence[int32](data.Offset, 0))
	return s.lRepo.SearchListingCombination(context.Background(), data)
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
			LCreatorID: types.Ptr[string](userId.String()),
		},
		SearchSortPaginationQuery: requests.SearchSortPaginationQuery{
			Limit:  types.Ptr[int32](1000),
			Offset: types.Ptr[int32](0),
			SortBy: types.Ptr[string]("created_at"),
			Order:  types.Ptr[string]("desc"),
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

	params.OrderInfo = fmt.Sprintf("[CREATE_LISTING-%s] Phi dang tin nha cho thue", data.ListingId.String())
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
