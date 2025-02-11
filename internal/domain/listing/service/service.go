package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"

	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
	"github.com/user2410/rrms-backend/pkg/ds/set"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/domain/listing/model"
	listing_utils "github.com/user2410/rrms-backend/internal/domain/listing/utils"
	payment_dto "github.com/user2410/rrms-backend/internal/domain/payment/dto"
	payment_model "github.com/user2410/rrms-backend/internal/domain/payment/model"
	payment_service "github.com/user2410/rrms-backend/internal/domain/payment/service"
	property_dto "github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/interfaces/rest/requests"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

type Service interface {
	CreateListing(data *dto.CreateListing) (*dto.CreateListingResponse, error)
	SearchListingCombination(data *dto.SearchListingCombinationQuery, userId uuid.UUID) (*dto.SearchListingCombinationResponse, error)
	GetListingByID(id uuid.UUID) (*model.ListingModel, error)
	GetListingsByIds(uid uuid.UUID, ids []uuid.UUID, fields []string) ([]model.ListingModel, error)
	GetListingsOfUser(userId uuid.UUID, query *dto.GetListingsQuery) (int, []model.ListingModel, error)
	GetListingPayments(id uuid.UUID) ([]payment_model.PaymentModel, error)

	UpdateListing(id uuid.UUID, data *dto.UpdateListing) error
	DeleteListing(id uuid.UUID) error
	CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckListingVisibility(lid uuid.UUID, uid uuid.UUID) (bool, error)
	CheckValidUnitForListing(lid uuid.UUID, uid uuid.UUID) (bool, error)
	CreateApplicationLink(data *dto.CreateApplicationLink) (string, error)
	VerifyApplicationLink(query *dto.VerifyApplicationLink) (bool, error)
	UpgradeListing(userId uuid.UUID, lid uuid.UUID, priority int) (*payment_model.PaymentModel, error)
	UpdateListingStatus(id uuid.UUID, active bool) error
	UpdateListingExpiration(id uuid.UUID, duration int64) error
	UpdateListingPriority(id uuid.UUID, priority int) error
	ExtendListing(userId uuid.UUID, lid uuid.UUID, duration int) (*payment_model.PaymentModel, error)
}

type service struct {
	hashSecret string
	domainRepo repos.DomainRepo

	esClient *es.ElasticSearchClient
}

func NewService(
	domainRepo repos.DomainRepo,
	hashSecret string,
	esClient *es.ElasticSearchClient,
) Service {
	return &service{
		hashSecret: hashSecret,
		domainRepo: domainRepo,
		esClient:   esClient,
	}
}

func (s *service) GetListingByID(id uuid.UUID) (*model.ListingModel, error) {
	return s.domainRepo.ListingRepo.GetListingByID(context.Background(), id)
}

func (s *service) GetListingsByIds(uid uuid.UUID, ids []uuid.UUID, fields []string) ([]model.ListingModel, error) {
	visibleIDS, err := s.FilterVisibleListings(ids, uid)
	if err != nil {
		return nil, err
	}
	return s.domainRepo.ListingRepo.GetListingsByIds(context.Background(), visibleIDS, fields)
}

func (s *service) DeleteListing(id uuid.UUID) error {
	return s.domainRepo.ListingRepo.DeleteListing(context.Background(), id)
}

func (s *service) CheckListingOwnership(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.domainRepo.ListingRepo.CheckListingOwnership(context.Background(), lid, uid)
}

func (s *service) CheckValidUnitForListing(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.domainRepo.ListingRepo.CheckValidUnitForListing(context.Background(), lid, uid)
}

func (s *service) GetListingsOfUser(userId uuid.UUID, query *dto.GetListingsQuery) (int, []model.ListingModel, error) {
	sspq := requests.SearchSortPaginationQuery{
		Limit:  types.Ptr[int32](math.MaxInt32),
		Offset: utils.Ternary(query.Offset == nil, types.Ptr[int32](0), query.Offset),
		SortBy: utils.Ternary(len(query.SortBy) == 0, []string{"listings.created_at"}, query.SortBy),
		Order:  utils.Ternary(len(query.Order) == 0, []string{"DESC"}, query.Order),
	}

	myListings, err := s.domainRepo.ListingRepo.SearchListingCombination(context.Background(), &dto.SearchListingCombinationQuery{
		SearchListingQuery: dto.SearchListingQuery{
			LCreatorID: types.Ptr(userId.String()),
		},
		SearchSortPaginationQuery: sspq,
	})
	if err != nil {
		return 0, nil, err
	}

	managedListings, err := s.domainRepo.ListingRepo.SearchListingCombination(context.Background(), &dto.SearchListingCombinationQuery{
		SearchPropertyQuery: property_dto.SearchPropertyQuery{
			PManagerIDS:  []string{userId.String()},
			PManagerRole: types.Ptr("OWNER"),
		},
		SearchSortPaginationQuery: sspq,
	})
	if err != nil {
		return 0, nil, err
	}

	lids := func() []uuid.UUID {
		ids := set.NewSet[uuid.UUID]()
		for _, listing := range myListings.Items {
			ids.Add(listing.LId)
		}
		for _, listing := range managedListings.Items {
			ids.Add(listing.LId)
		}
		return ids.ToSlice()
	}()

	total := len(lids)
	var actualLength int
	if query.Limit == nil {
		actualLength = total
	} else {
		actualLength = utils.Ternary(total > int(*query.Limit), int(*query.Limit), total)
	}
	items, err := s.domainRepo.ListingRepo.GetListingsByIds(context.Background(), lids[0:actualLength], query.Fields)

	return total, items, err
}

func (s *service) GetListingPayments(id uuid.UUID) ([]payment_model.PaymentModel, error) {
	return s.domainRepo.ListingRepo.GetListingPayments(context.Background(), id)
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

func (s *service) FilterVisibleListings(lids []uuid.UUID, uid uuid.UUID) ([]uuid.UUID, error) {
	lidSet := set.NewSet[uuid.UUID]()
	lidSet.AddAll(lids...)
	return s.domainRepo.ListingRepo.FilterVisibleListings(context.Background(), lidSet.ToSlice(), uid)
}

func (s *service) CheckListingVisibility(lid uuid.UUID, uid uuid.UUID) (bool, error) {
	return s.domainRepo.ListingRepo.CheckListingVisibility(context.Background(), lid, uid)
}

func (s *service) UpgradeListing(userId uuid.UUID, lid uuid.UUID, priority int) (*payment_model.PaymentModel, error) {
	listing, err := s.domainRepo.ListingRepo.GetListingByID(context.Background(), lid)
	if err != nil {
		return nil, err
	}

	payments, err := s.domainRepo.ListingRepo.GetListingPaymentsByType(context.Background(), lid, payment_service.PAYMENTTYPE_UPGRADELISTING)
	if err != nil && errors.Is(err, database.ErrRecordNotFound) {
		return nil, err
	}
	if func() bool {
		found := false
		for _, payment := range payments {
			if payment.Status != database.PAYMENTSTATUSSUCCESS {
				found = true
				break
			}
		}
		return found
	}() {
		return nil, ErrUnpaidPayment
	}

	params := payment_dto.CreatePayment{UserId: userId}
	amount, discount, err := listing_utils.CalculateUpgradeListingPrice(listing, priority)
	if err != nil {
		return nil, err
	}
	params.Amount = amount
	params.OrderInfo = fmt.Sprintf(
		"[%s%s%s%s%d] Phi nang cap tin dang nha cho thue",
		payment_service.PAYMENTTYPE_UPGRADELISTING, payment_service.PAYMENTTYPE_DELIMITER, listing.ID.String(), payment_service.PAYMENTTYPE_DELIMITER, priority,
	)
	params.Items = []payment_dto.CreatePaymentItem{
		{
			Name:     "Phi nang cap",
			Price:    amount,
			Quantity: 1,
			Discount: int32(discount),
		},
	}

	payment, err := s.domainRepo.PaymentRepo.CreatePayment(context.Background(), &params)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *service) ExtendListing(userId uuid.UUID, lid uuid.UUID, duration int) (*payment_model.PaymentModel, error) {
	listing, err := s.domainRepo.ListingRepo.GetListingByID(context.Background(), lid)
	if err != nil {
		return nil, err
	}

	payments, err := s.domainRepo.ListingRepo.GetListingPaymentsByType(context.Background(), lid, payment_service.PAYMENTTYPE_EXTENDLISTING)
	if err != nil && errors.Is(err, database.ErrRecordNotFound) {
		return nil, err
	}
	if func() bool {
		found := false
		for _, payment := range payments {
			if payment.Status != database.PAYMENTSTATUSSUCCESS {
				found = true
				break
			}
		}
		return found
	}() {
		return nil, ErrUnpaidPayment
	}

	params := payment_dto.CreatePayment{UserId: userId}
	amount, discount, err := listing_utils.CalculateExtendListingPrice(listing, duration)
	if err != nil {
		return nil, err
	}
	params.Amount = amount
	params.OrderInfo = fmt.Sprintf(
		"[%s%s%s%s%d] Phi gia han tin dang nha cho thue",
		payment_service.PAYMENTTYPE_EXTENDLISTING, payment_service.PAYMENTTYPE_DELIMITER, listing.ID.String(), payment_service.PAYMENTTYPE_DELIMITER, duration,
	)
	params.Items = []payment_dto.CreatePaymentItem{
		{
			Name:     "Phi gia han",
			Price:    amount,
			Quantity: 1,
			Discount: int32(discount),
		},
	}

	payment, err := s.domainRepo.PaymentRepo.CreatePayment(context.Background(), &params)
	if err != nil {
		return nil, err
	}
	return payment, nil
}
