package services

import (
	application_service "github.com/user2410/rrms-backend/internal/domain/application/service"
	auth_service "github.com/user2410/rrms-backend/internal/domain/auth/service"
	"github.com/user2410/rrms-backend/internal/domain/chat"
	listing_service "github.com/user2410/rrms-backend/internal/domain/listing/service"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	payment_service "github.com/user2410/rrms-backend/internal/domain/payment/service"
	property_service "github.com/user2410/rrms-backend/internal/domain/property/service"
	"github.com/user2410/rrms-backend/internal/domain/reminder"
	rental_service "github.com/user2410/rrms-backend/internal/domain/rental/service"
	statistic_service "github.com/user2410/rrms-backend/internal/domain/statistic/service"
	unit_service "github.com/user2410/rrms-backend/internal/domain/unit/service"
)

type DomainServices struct {
	AuthService        auth_service.Service
	PropertyService    property_service.Service
	UnitService        unit_service.Service
	ListingService     listing_service.Service
	RentalService      rental_service.Service
	ApplicationService application_service.Service
	ReminderService    reminder.Service
	PaymentService     payment_service.Service
	ChatService        chat.Service
	StatisticService   statistic_service.Service
	MiscService        misc_service.Service
}
