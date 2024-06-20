// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	AddPropertyManager(ctx context.Context, arg AddPropertyManagerParams) error
	CheckApplicationUpdatabilty(ctx context.Context, arg CheckApplicationUpdatabiltyParams) (bool, error)
	CheckApplicationVisibility(ctx context.Context, arg CheckApplicationVisibilityParams) (bool, error)
	CheckListingExpired(ctx context.Context, id uuid.UUID) (pgtype.Bool, error)
	CheckListingOwnership(ctx context.Context, arg CheckListingOwnershipParams) (int64, error)
	CheckListingVisibility(ctx context.Context, arg CheckListingVisibilityParams) (bool, error)
	CheckMsgGroupMembership(ctx context.Context, arg CheckMsgGroupMembershipParams) (bool, error)
	CheckPaymentAccessible(ctx context.Context, arg CheckPaymentAccessibleParams) (bool, error)
	CheckReminderVisibility(ctx context.Context, arg CheckReminderVisibilityParams) (bool, error)
	CheckRentalVisibility(ctx context.Context, arg CheckRentalVisibilityParams) (bool, error)
	CheckUnitManageability(ctx context.Context, arg CheckUnitManageabilityParams) (int64, error)
	CheckUnitOfProperty(ctx context.Context, arg CheckUnitOfPropertyParams) (int64, error)
	CheckValidUnitForListing(ctx context.Context, arg CheckValidUnitForListingParams) (int64, error)
	CreateApplication(ctx context.Context, arg CreateApplicationParams) (Application, error)
	CreateApplicationCoap(ctx context.Context, arg CreateApplicationCoapParams) (ApplicationCoap, error)
	CreateApplicationMinor(ctx context.Context, arg CreateApplicationMinorParams) (ApplicationMinor, error)
	CreateApplicationPet(ctx context.Context, arg CreateApplicationPetParams) (ApplicationPet, error)
	CreateApplicationVehicle(ctx context.Context, arg CreateApplicationVehicleParams) (ApplicationVehicle, error)
	CreateContract(ctx context.Context, arg CreateContractParams) (Contract, error)
	CreateListing(ctx context.Context, arg CreateListingParams) (Listing, error)
	CreateListingPolicy(ctx context.Context, arg CreateListingPolicyParams) (ListingPolicy, error)
	CreateListingTag(ctx context.Context, arg CreateListingTagParams) (ListingTag, error)
	CreateListingUnit(ctx context.Context, arg CreateListingUnitParams) (ListingUnit, error)
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	CreateMsgGroup(ctx context.Context, arg CreateMsgGroupParams) (MsgGroup, error)
	CreateMsgGroupMember(ctx context.Context, arg CreateMsgGroupMemberParams) (MsgGroupMember, error)
	CreateNewPropertyManagerRequest(ctx context.Context, arg CreateNewPropertyManagerRequestParams) (NewPropertyManagerRequest, error)
	CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error)
	CreateNotificationDevice(ctx context.Context, arg CreateNotificationDeviceParams) (UserNotificationDevice, error)
	CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error)
	CreatePaymentItem(ctx context.Context, arg CreatePaymentItemParams) (PaymentItem, error)
	CreateProperty(ctx context.Context, arg CreatePropertyParams) (Property, error)
	CreatePropertyFeature(ctx context.Context, arg CreatePropertyFeatureParams) (PropertyFeature, error)
	CreatePropertyManager(ctx context.Context, arg CreatePropertyManagerParams) (PropertyManager, error)
	CreatePropertyMedia(ctx context.Context, arg CreatePropertyMediaParams) (PropertyMedium, error)
	CreatePropertyTag(ctx context.Context, arg CreatePropertyTagParams) (PropertyTag, error)
	CreatePropertyVerificationRequest(ctx context.Context, arg CreatePropertyVerificationRequestParams) (PropertyVerificationRequest, error)
	CreateReminder(ctx context.Context, arg CreateReminderParams) (Reminder, error)
	CreateRental(ctx context.Context, arg CreateRentalParams) (Rental, error)
	CreateRentalCoap(ctx context.Context, arg CreateRentalCoapParams) (RentalCoap, error)
	CreateRentalComplaint(ctx context.Context, arg CreateRentalComplaintParams) (RentalComplaint, error)
	CreateRentalComplaintReply(ctx context.Context, arg CreateRentalComplaintReplyParams) (RentalComplaintReply, error)
	CreateRentalMinor(ctx context.Context, arg CreateRentalMinorParams) (RentalMinor, error)
	CreateRentalPayment(ctx context.Context, arg CreateRentalPaymentParams) (RentalPayment, error)
	CreateRentalPet(ctx context.Context, arg CreateRentalPetParams) (RentalPet, error)
	CreateRentalPolicy(ctx context.Context, arg CreateRentalPolicyParams) (RentalPolicy, error)
	CreateRentalService(ctx context.Context, arg CreateRentalServiceParams) (RentalService, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUnit(ctx context.Context, arg CreateUnitParams) (Unit, error)
	CreateUnitAmenity(ctx context.Context, arg CreateUnitAmenityParams) (UnitAmenity, error)
	CreateUnitMedia(ctx context.Context, arg CreateUnitMediaParams) (UnitMedium, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteApplication(ctx context.Context, id int64) error
	DeleteExpiredTokens(ctx context.Context, interval int32) error
	DeleteListing(ctx context.Context, id uuid.UUID) error
	DeleteListingPolicies(ctx context.Context, listingID uuid.UUID) error
	DeleteListingTags(ctx context.Context, listingID uuid.UUID) error
	DeleteListingUnits(ctx context.Context, listingID uuid.UUID) error
	DeleteMsgGroup(ctx context.Context, groupID int64) error
	DeleteMsgGroupMember(ctx context.Context, arg DeleteMsgGroupMemberParams) error
	DeleteNotificationDeviceToken(ctx context.Context, arg DeleteNotificationDeviceTokenParams) error
	DeletePayment(ctx context.Context, id int64) error
	DeleteProperty(ctx context.Context, id uuid.UUID) error
	DeletePropertyFeature(ctx context.Context, arg DeletePropertyFeatureParams) error
	DeletePropertyManager(ctx context.Context, arg DeletePropertyManagerParams) error
	DeletePropertyMedia(ctx context.Context, arg DeletePropertyMediaParams) error
	DeletePropertyTag(ctx context.Context, arg DeletePropertyTagParams) error
	DeleteReminder(ctx context.Context, id int64) error
	DeleteRental(ctx context.Context, id int64) error
	DeleteUnit(ctx context.Context, id uuid.UUID) error
	DeleteUnitAmenity(ctx context.Context, arg DeleteUnitAmenityParams) error
	DeleteUnitMedia(ctx context.Context, arg DeleteUnitMediaParams) error
	GetAllPropertyFeatures(ctx context.Context) ([]PFeature, error)
	GetAllRentalPolicies(ctx context.Context) ([]LPolicy, error)
	GetAllUnitAmenities(ctx context.Context) ([]UAmenity, error)
	GetApplicationByID(ctx context.Context, id int64) (Application, error)
	GetApplicationCoaps(ctx context.Context, applicationID int64) ([]ApplicationCoap, error)
	GetApplicationMinors(ctx context.Context, applicationID int64) ([]ApplicationMinor, error)
	GetApplicationPets(ctx context.Context, applicationID int64) ([]ApplicationPet, error)
	GetApplicationVehicles(ctx context.Context, applicationID int64) ([]ApplicationVehicle, error)
	GetApplicationsByUserId(ctx context.Context, arg GetApplicationsByUserIdParams) ([]int64, error)
	GetApplicationsOfListing(ctx context.Context, listingID uuid.UUID) ([]int64, error)
	GetApplicationsToUser(ctx context.Context, arg GetApplicationsToUserParams) ([]int64, error)
	GetContractByID(ctx context.Context, id int64) (Contract, error)
	GetContractByRentalID(ctx context.Context, rentalID int64) (Contract, error)
	GetLeastRentedProperties(ctx context.Context, arg GetLeastRentedPropertiesParams) ([]GetLeastRentedPropertiesRow, error)
	GetLeastRentedUnits(ctx context.Context, arg GetLeastRentedUnitsParams) ([]GetLeastRentedUnitsRow, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (Listing, error)
	GetListingPolicies(ctx context.Context, listingID uuid.UUID) ([]ListingPolicy, error)
	GetListingTags(ctx context.Context, listingID uuid.UUID) ([]ListingTag, error)
	GetListingUnits(ctx context.Context, listingID uuid.UUID) ([]ListingUnit, error)
	GetListingsCountByCity(ctx context.Context, city string) (int64, error)
	// Get expired / active listings
	GetListingsOfProperty(ctx context.Context, arg GetListingsOfPropertyParams) ([]uuid.UUID, error)
	GetMaintenanceRequests(ctx context.Context, arg GetMaintenanceRequestsParams) ([]int64, error)
	GetManagedPropertiesByRole(ctx context.Context, arg GetManagedPropertiesByRoleParams) ([]uuid.UUID, error)
	GetManagedRentals(ctx context.Context, arg GetManagedRentalsParams) ([]int64, error)
	GetManagedUnits(ctx context.Context, managerID uuid.UUID) ([]uuid.UUID, error)
	GetMessagesOfGroup(ctx context.Context, arg GetMessagesOfGroupParams) ([]Message, error)
	GetMostRentedProperties(ctx context.Context, arg GetMostRentedPropertiesParams) ([]GetMostRentedPropertiesRow, error)
	GetMostRentedUnits(ctx context.Context, arg GetMostRentedUnitsParams) ([]GetMostRentedUnitsRow, error)
	GetMsgGroup(ctx context.Context, groupID int64) (MsgGroup, error)
	GetMsgGroupByName(ctx context.Context, arg GetMsgGroupByNameParams) (MsgGroup, error)
	GetMsgGroupMembers(ctx context.Context, groupID int64) ([]GetMsgGroupMembersRow, error)
	GetMyRentals(ctx context.Context, arg GetMyRentalsParams) ([]int64, error)
	GetNewApplications(ctx context.Context, arg GetNewApplicationsParams) ([]int64, error)
	GetNewPropertyManagerRequest(ctx context.Context, id int64) (NewPropertyManagerRequest, error)
	GetNewPropertyManagerRequestsToUser(ctx context.Context, arg GetNewPropertyManagerRequestsToUserParams) ([]NewPropertyManagerRequest, error)
	GetNotification(ctx context.Context, id int64) (Notification, error)
	GetNotificationDevice(ctx context.Context, arg GetNotificationDeviceParams) (UserNotificationDevice, error)
	GetNotificationsOfUser(ctx context.Context, arg GetNotificationsOfUserParams) ([]Notification, error)
	GetOccupiedProperties(ctx context.Context, managerID uuid.UUID) ([]uuid.UUID, error)
	GetOccupiedUnits(ctx context.Context, managerID uuid.UUID) ([]uuid.UUID, error)
	GetPaymentById(ctx context.Context, id int64) (Payment, error)
	GetPaymentItemsByPaymentId(ctx context.Context, paymentID int64) ([]PaymentItem, error)
	GetPaymentsOfRental(ctx context.Context, rentalID int64) ([]RentalPayment, error)
	GetPaymentsOfUser(ctx context.Context, arg GetPaymentsOfUserParams) ([]Payment, error)
	GetPaymentsStatistic(ctx context.Context, arg GetPaymentsStatisticParams) (float32, error)
	GetPropertiesWithActiveListing(ctx context.Context, managerID uuid.UUID) ([]uuid.UUID, error)
	GetPropertyById(ctx context.Context, id uuid.UUID) (Property, error)
	GetPropertyFeatures(ctx context.Context, propertyID uuid.UUID) ([]PropertyFeature, error)
	GetPropertyManagers(ctx context.Context, propertyID uuid.UUID) ([]PropertyManager, error)
	GetPropertyMedia(ctx context.Context, propertyID uuid.UUID) ([]PropertyMedium, error)
	GetPropertyTags(ctx context.Context, propertyID uuid.UUID) ([]PropertyTag, error)
	GetPropertyVerificationRequest(ctx context.Context, id int64) (PropertyVerificationRequest, error)
	GetPropertyVerificationRequestsOfProperty(ctx context.Context, arg GetPropertyVerificationRequestsOfPropertyParams) ([]PropertyVerificationRequest, error)
	GetRecentListings(ctx context.Context, limit int32) ([]uuid.UUID, error)
	GetReminderById(ctx context.Context, id int64) (Reminder, error)
	GetRemindersByCreator(ctx context.Context, creatorID uuid.UUID) ([]Reminder, error)
	GetRemindersInDate(ctx context.Context, dateTrunc pgtype.Interval) ([]Reminder, error)
	GetRemindersOfUserWithResourceTag(ctx context.Context, arg GetRemindersOfUserWithResourceTagParams) ([]Reminder, error)
	GetRental(ctx context.Context, id int64) (Rental, error)
	GetRentalByApplicationId(ctx context.Context, applicationID pgtype.Int8) (Rental, error)
	GetRentalCoapsByRentalID(ctx context.Context, rentalID int64) ([]RentalCoap, error)
	GetRentalComplaint(ctx context.Context, id int64) (RentalComplaint, error)
	GetRentalComplaintReplies(ctx context.Context, arg GetRentalComplaintRepliesParams) ([]RentalComplaintReply, error)
	GetRentalComplaintStatistics(ctx context.Context, arg GetRentalComplaintStatisticsParams) (int64, error)
	GetRentalComplaintsByRentalId(ctx context.Context, rentalID int64) ([]RentalComplaint, error)
	GetRentalComplaintsOfUser(ctx context.Context, arg GetRentalComplaintsOfUserParams) ([]RentalComplaint, error)
	GetRentalContractsOfUser(ctx context.Context, arg GetRentalContractsOfUserParams) ([]int64, error)
	GetRentalMinorsByRentalID(ctx context.Context, rentalID int64) ([]RentalMinor, error)
	GetRentalPayment(ctx context.Context, id int64) (RentalPayment, error)
	GetRentalPaymentArrears(ctx context.Context, arg GetRentalPaymentArrearsParams) ([]GetRentalPaymentArrearsRow, error)
	GetRentalPaymentIncomes(ctx context.Context, arg GetRentalPaymentIncomesParams) (float32, error)
	GetRentalPetsByRentalID(ctx context.Context, rentalID int64) ([]RentalPet, error)
	GetRentalPoliciesByRentalID(ctx context.Context, rentalID int64) ([]RentalPolicy, error)
	GetRentalServicesByRentalID(ctx context.Context, rentalID int64) ([]RentalService, error)
	// Get rental side: Side A (lanlord and managers) and Side B (tenant). Otherwise return C
	GetRentalSide(ctx context.Context, arg GetRentalSideParams) (string, error)
	GetRentalsOfProperty(ctx context.Context, arg GetRentalsOfPropertyParams) ([]int64, error)
	GetRentedProperties(ctx context.Context, tenantID pgtype.UUID) ([]uuid.UUID, error)
	GetSessionById(ctx context.Context, id uuid.UUID) (Session, error)
	GetSomeListings(ctx context.Context, arg GetSomeListingsParams) ([]Listing, error)
	GetTenantExpenditure(ctx context.Context, arg GetTenantExpenditureParams) (float32, error)
	GetTenantPendingPayments(ctx context.Context, arg GetTenantPendingPaymentsParams) ([]GetTenantPendingPaymentsRow, error)
	GetTotalTenantPendingPayments(ctx context.Context, userID pgtype.UUID) (float32, error)
	GetUnitAmenities(ctx context.Context, unitID uuid.UUID) ([]UnitAmenity, error)
	GetUnitById(ctx context.Context, id uuid.UUID) (Unit, error)
	GetUnitManagers(ctx context.Context, id uuid.UUID) ([]PropertyManager, error)
	GetUnitMedia(ctx context.Context, unitID uuid.UUID) ([]UnitMedium, error)
	GetUnitsOfProperty(ctx context.Context, propertyID uuid.UUID) ([]Unit, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (User, error)
	IsPropertyVisible(ctx context.Context, arg IsPropertyVisibleParams) (pgtype.Bool, error)
	IsUnitPublic(ctx context.Context, id uuid.UUID) (bool, error)
	PingContractByRentalID(ctx context.Context, rentalID int64) (PingContractByRentalIDRow, error)
	PlanRentalPayment(ctx context.Context, rentalID int64) ([]int64, error)
	PlanRentalPayments(ctx context.Context) ([]int64, error)
	UpdateApplicationStatus(ctx context.Context, arg UpdateApplicationStatusParams) ([]int64, error)
	UpdateContract(ctx context.Context, arg UpdateContractParams) error
	UpdateContractContent(ctx context.Context, arg UpdateContractContentParams) error
	UpdateFinePayments(ctx context.Context) error
	UpdateFinePaymentsOfRental(ctx context.Context, rentalID int64) error
	UpdateListing(ctx context.Context, arg UpdateListingParams) error
	UpdateListingPriority(ctx context.Context, arg UpdateListingPriorityParams) error
	UpdateListingStatus(ctx context.Context, arg UpdateListingStatusParams) error
	UpdateMessage(ctx context.Context, arg UpdateMessageParams) ([]int64, error)
	UpdateNewPropertyManagerRequest(ctx context.Context, arg UpdateNewPropertyManagerRequestParams) error
	UpdateNotificationDeviceTokenTimestamp(ctx context.Context, arg UpdateNotificationDeviceTokenTimestampParams) error
	UpdatePayment(ctx context.Context, arg UpdatePaymentParams) error
	UpdateProperty(ctx context.Context, arg UpdatePropertyParams) error
	UpdatePropertyVerificationRequest(ctx context.Context, arg UpdatePropertyVerificationRequestParams) error
	UpdateReminder(ctx context.Context, arg UpdateReminderParams) ([]Reminder, error)
	UpdateRental(ctx context.Context, arg UpdateRentalParams) error
	UpdateRentalComplaint(ctx context.Context, arg UpdateRentalComplaintParams) error
	UpdateRentalPayment(ctx context.Context, arg UpdateRentalPaymentParams) error
	UpdateSessionBlockingStatus(ctx context.Context, arg UpdateSessionBlockingStatusParams) error
	UpdateUnit(ctx context.Context, arg UpdateUnitParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdatedNotification(ctx context.Context, arg UpdatedNotificationParams) error
}

var _ Querier = (*Queries)(nil)
