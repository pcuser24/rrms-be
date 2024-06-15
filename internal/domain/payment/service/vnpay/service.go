package vnpay

import (
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	"github.com/user2410/rrms-backend/internal/domain/payment/service"
)

type VnPayService struct {
	service.PaymentService
	domainRepo    repos.DomainRepo
	lService      listing_service.Service
	vnpTmnCode    string
	vnpHashSecret string
	vnpUrl        string
	vnpApi        string
}

func NewVnpayService(
	domainRepo repos.DomainRepo, lService listing_service.Service,
	vnpTmnCode string, vnpHashSecret string, vnpUrl string, vnpApi string,
) service.Service {
	return &VnPayService{
		PaymentService: service.NewPaymentService(domainRepo.PaymentRepo),
		domainRepo:     domainRepo,
		lService:       lService,
		vnpTmnCode:     vnpTmnCode,
		vnpHashSecret:  vnpHashSecret,
		vnpUrl:         vnpUrl,
		vnpApi:         vnpApi,
	}
}
