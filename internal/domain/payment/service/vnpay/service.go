package vnpay

import (
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	"github.com/user2410/rrms-backend/internal/domain/payment/repo"
	"github.com/user2410/rrms-backend/internal/domain/payment/service"
)

type VnPayService struct {
	service.PaymentService
	repo          repo.Repo
	lRepo         listing_repo.Repo
	vnpTmnCode    string
	vnpHashSecret string
	vnpUrl        string
	vnpApi        string
}

func NewVnpayService(
	repo repo.Repo, lRepo listing_repo.Repo,
	vnpTmnCode string, vnpHashSecret string, vnpUrl string, vnpApi string,
) service.Service {
	return &VnPayService{
		PaymentService: service.NewPaymentService(repo),
		repo:           repo,
		lRepo:          lRepo,
		vnpTmnCode:     vnpTmnCode,
		vnpHashSecret:  vnpHashSecret,
		vnpUrl:         vnpUrl,
		vnpApi:         vnpApi,
	}
}
