package vnpay

import (
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	"github.com/user2410/rrms-backend/internal/domain/payment/repo"
	"github.com/user2410/rrms-backend/internal/domain/payment/service"
)

type VnPayService struct {
	service.PaymentService
	repo          repo.Repo
	lService      listing_service.Service
	vnpTmnCode    string
	vnpHashSecret string
	vnpUrl        string
	vnpApi        string
}

func NewVnpayService(
	repo repo.Repo, lService listing_service.Service,
	vnpTmnCode string, vnpHashSecret string, vnpUrl string, vnpApi string,
) service.Service {
	return &VnPayService{
		PaymentService: service.NewPaymentService(repo),
		repo:           repo,
		lService:       lService,
		vnpTmnCode:     vnpTmnCode,
		vnpHashSecret:  vnpHashSecret,
		vnpUrl:         vnpUrl,
		vnpApi:         vnpApi,
	}
}
