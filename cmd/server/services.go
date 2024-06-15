package server

import (
	repos "github.com/user2410/rrms-backend/internal/domain/_repos"
	application_service "github.com/user2410/rrms-backend/internal/domain/application/service"
	auth_service "github.com/user2410/rrms-backend/internal/domain/auth/service"
	"github.com/user2410/rrms-backend/internal/domain/chat"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	vnp_service "github.com/user2410/rrms-backend/internal/domain/payment/service/vnpay"
	property_service "github.com/user2410/rrms-backend/internal/domain/property/service"
	"github.com/user2410/rrms-backend/internal/domain/reminder"
	rental_service "github.com/user2410/rrms-backend/internal/domain/rental/service"
	statistic_service "github.com/user2410/rrms-backend/internal/domain/statistic/service"
	unit_service "github.com/user2410/rrms-backend/internal/domain/unit/service"

	application_repo "github.com/user2410/rrms-backend/internal/domain/application/repo"
	auth_repo "github.com/user2410/rrms-backend/internal/domain/auth/repo"
	chat_repo "github.com/user2410/rrms-backend/internal/domain/chat/repo"
	listing_repo "github.com/user2410/rrms-backend/internal/domain/listing/repo"
	misc_repo "github.com/user2410/rrms-backend/internal/domain/misc/repo"
	payment_repo "github.com/user2410/rrms-backend/internal/domain/payment/repo"
	property_repo "github.com/user2410/rrms-backend/internal/domain/property/repo"
	reminder_repo "github.com/user2410/rrms-backend/internal/domain/reminder/repo"
	rental_repo "github.com/user2410/rrms-backend/internal/domain/rental/repo"
	statistic_repo "github.com/user2410/rrms-backend/internal/domain/statistic/repo"
	unit_repo "github.com/user2410/rrms-backend/internal/domain/unit/repo"
)

func (c *serverCommand) setupInternalServices() {
	var domainRepo repos.DomainRepo
	// Initialize repositories
	domainRepo.AuthRepo = auth_repo.NewRepo(c.dao)
	domainRepo.PropertyRepo = property_repo.NewRepo(c.dao)
	domainRepo.UnitRepo = unit_repo.NewRepo(c.dao)
	domainRepo.ListingRepo = listing_repo.NewRepo(c.dao)
	domainRepo.RentalRepo = rental_repo.NewRepo(c.dao)
	domainRepo.ApplicationRepo = application_repo.NewRepo(c.dao)
	domainRepo.PaymentRepo = payment_repo.NewRepo(c.dao)
	domainRepo.ChatRepo = chat_repo.NewRepo(c.dao)
	domainRepo.ReminderRepo = reminder_repo.NewRepo(c.dao)
	domainRepo.StatisticRepo = statistic_repo.NewRepo(c.dao)
	domainRepo.MiscRepo = misc_repo.NewRepo(c.dao)

	// Initialize internal services
	c.internalServices.MiscService = misc_service.NewService(domainRepo, c.notificationEndpoint, c.cronScheduler)
	c.internalServices.AuthService = auth_service.NewService(
		domainRepo,
		c.tokenMaker, c.config.AccessTokenTTL, c.config.RefreshTokenTTL,
	)
	c.internalServices.PropertyService = property_service.NewService(
		domainRepo,
		c.s3Client, c.config.AWSS3ImageBucket,
	)
	c.internalServices.UnitService = unit_service.NewService(domainRepo, c.s3Client, c.config.AWSS3ImageBucket)
	c.internalServices.ListingService = listing_service.NewService(
		domainRepo,
		c.config.TokenSecreteKey,
	)
	c.internalServices.RentalService = rental_service.NewService(
		domainRepo,
		c.internalServices.MiscService,
		c.cronScheduler,
		c.s3Client, c.config.AWSS3ImageBucket,
	)
	c.internalServices.ReminderService = reminder.NewService(
		domainRepo,
	)
	c.internalServices.ApplicationService = application_service.NewService(
		domainRepo,
		c.internalServices.ReminderService, c.internalServices.MiscService,
		c.s3Client, c.config.AWSS3ImageBucket,
		c.config.FESite,
	)
	c.internalServices.PaymentService = vnp_service.NewVnpayService(
		domainRepo,
		c.internalServices.ListingService,
		c.config.VnpTmnCode, c.config.VnpHashSecret, c.config.VnpUrl, c.config.VnpApi,
	)
	c.internalServices.ChatService = chat.NewService(domainRepo.ChatRepo)
	c.internalServices.StatisticService = statistic_service.NewService(
		domainRepo,
		c.elasticsearch,
	)
}
