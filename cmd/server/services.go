package server

import (
	application_service "github.com/user2410/rrms-backend/internal/domain/application/service"
	"github.com/user2410/rrms-backend/internal/domain/auth"
	"github.com/user2410/rrms-backend/internal/domain/chat"
	"github.com/user2410/rrms-backend/internal/domain/listing"
	payment_service "github.com/user2410/rrms-backend/internal/domain/payment/service"
	vnp_service "github.com/user2410/rrms-backend/internal/domain/payment/service/vnpay"
	"github.com/user2410/rrms-backend/internal/domain/property"
	"github.com/user2410/rrms-backend/internal/domain/reminder"
	rental_service "github.com/user2410/rrms-backend/internal/domain/rental/service"
	"github.com/user2410/rrms-backend/internal/domain/storage"
	"github.com/user2410/rrms-backend/internal/domain/unit"
	"github.com/user2410/rrms-backend/internal/infrastructure/asynctask"
	"github.com/user2410/rrms-backend/internal/infrastructure/aws/s3"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"

	application_asynctask "github.com/user2410/rrms-backend/internal/domain/application/asynctask"
	auth_asynctask "github.com/user2410/rrms-backend/internal/domain/auth/asynctask"

	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	payment_repo "github.com/user2410/rrms-backend/internal/domain/payment/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	reminder_repo "github.com/user2410/rrms-backend/internal/domain/reminder/repo"
	rental_repo "github.com/user2410/rrms-backend/internal/domain/rental/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
)

func (c *serverCommand) setupInternalServices(
	dao database.DAO,
	s3Client *s3.S3Client,
	taskDistributor asynctask.Distributor,
) {
	// Initialize repositories
	authRepo := auth_repo.NewRepo(dao)
	propertyRepo := property_repo.NewRepo(dao)
	unitRepo := unit_repo.NewRepo(dao)
	listingRepo := listing_repo.NewRepo(dao)
	rentalRepo := rental_repo.NewRepo(dao)
	applicationRepo := application_repo.NewRepo(dao)
	paymentRepo := payment_repo.NewRepo(dao)
	chatRepo := chat_repo.NewRepo(dao)
	reminderRepo := reminder_repo.NewRepo(dao)

	// Initialize storage services
	s := storage.NewStorage(s3Client, c.config.AWSS3ImageBucket)

	// Initialize async task distributor
	authTaskDistributor := auth_asynctask.NewTaskDistributor(taskDistributor)
	applicationTaskDistributor := application_asynctask.NewTaskDistributor(taskDistributor)

	// Initialize internal services
	c.internalServices.AuthService = auth.NewService(
		authRepo,
		c.tokenMaker, c.config.AccessTokenTTL, c.config.RefreshTokenTTL,
		authTaskDistributor,
	)
	c.internalServices.PropertyService = property.NewService(
		propertyRepo,
		unitRepo,
		listingRepo,
		applicationRepo,
	)
	c.internalServices.UnitService = unit.NewService(unitRepo)
	c.internalServices.ListingService = listing.NewService(
		listingRepo,
		propertyRepo,
		paymentRepo,
		c.config.TokenSecreteKey,
	)
	c.internalServices.RentalService = rental_service.NewService(
		rentalRepo,
		authRepo,
		applicationRepo,
		listingRepo,
		propertyRepo,
		unitRepo,
		c.cronScheduler,
	)
	c.internalServices.ReminderService = reminder.NewService(
		reminderRepo,
		c.wsNotificationAdapter,
	)
	c.internalServices.ApplicationService = application_service.NewService(
		applicationRepo,
		chatRepo,
		listingRepo,
		propertyRepo,
		c.internalServices.ReminderService,
		applicationTaskDistributor,
	)
	c.internalServices.PaymentService = payment_service.NewService(paymentRepo)
	c.internalServices.StorageService = storage.NewService(s)
	c.internalServices.VnpService = vnp_service.NewVnpayService(
		paymentRepo,
		listingRepo,
		c.config.VnpTmnCode, c.config.VnpHashSecret, c.config.VnpUrl, c.config.VnpApi,
	)
	c.internalServices.ChatService = chat.NewService(chatRepo)
}
