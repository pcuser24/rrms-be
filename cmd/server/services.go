package server

import (
	"github.com/user2410/rrms-backend/internal/domain/application"
	application_asynctask "github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	auth_asynctask "github.com/user2410/rrms-backend/internal/domain/auth/asynctask"
	"github.com/user2410/rrms-backend/internal/domain/chat"
	"github.com/user2410/rrms-backend/internal/domain/listing"
	payment_service "github.com/user2410/rrms-backend/internal/domain/payment/service"
	vnp_service "github.com/user2410/rrms-backend/internal/domain/payment/service/vnpay"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/rental"
	"github.com/user2410/rrms-backend/internal/domain/storage"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"

	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	payment_repo "github.com/user2410/rrms-backend/internal/domain/payment/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"

	"github.com/hibiken/asynq"
)

func (c *serverCommand) setupInternalServices(
	dao database.DAO,
	s3Client *s3.S3Client,
) {
	c.asyncTaskDistributor = asynctask.NewRedisTaskDistributor(asynq.RedisClientOpt{
		Addr: c.config.AsynqRedisAddress,
	})

	authRepo := auth_repo.NewRepo(dao)
	authTaskDistributor := auth_asynctask.NewTaskDistributor(c.asyncTaskDistributor)
	c.internalServices.AuthService = auth.NewService(
		authRepo,
		c.tokenMaker, c.config.AccessTokenTTL, c.config.RefreshTokenTTL,
		authTaskDistributor,
	)
	propertyRepo := property_repo.NewRepo(dao)
	unitRepo := unit_repo.NewRepo(dao)
	listingRepo := listing_repo.NewRepo(dao)
	rentalRepo := rental.NewRepo(dao)
	applicationRepo := application_repo.NewRepo(dao)
	paymentRepo := payment_repo.NewRepo(dao)
	chatRepo := chat_repo.NewRepo(dao)

	s := storage.NewStorage(s3Client, c.config.AWSS3ImageBucket)

	c.internalServices.PropertyService = property.NewService(propertyRepo, unitRepo, listingRepo, applicationRepo)
	c.internalServices.UnitService = unit.NewService(unitRepo)
	c.internalServices.ListingService = listing.NewService(listingRepo, propertyRepo, paymentRepo, c.config.TokenSecreteKey)
	c.internalServices.RentalService = rental.NewService(rentalRepo)
	applicationTaskDistributor := application_asynctask.NewTaskDistributor(c.asyncTaskDistributor)
	c.internalServices.ApplicationService = application.NewService(
		applicationRepo,
		chatRepo,
		listingRepo,
		propertyRepo,
		applicationTaskDistributor,
	)
	c.internalServices.PaymentService = payment_service.NewService(paymentRepo)
	c.internalServices.StorageService = storage.NewService(s)
	c.internalServices.VnpService = vnp_service.NewVnpayService(
		paymentRepo,
		c.config.VnpTmnCode,
		c.config.VnpHashSecret,
		c.config.VnpUrl,
		c.config.VnpApi,
	)
	c.internalServices.ChatService = chat.NewService(chatRepo)
}
