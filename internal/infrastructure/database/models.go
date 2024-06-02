// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type APPLICATIONSTATUS string

const (
	APPLICATIONSTATUSPENDING               APPLICATIONSTATUS = "PENDING"
	APPLICATIONSTATUSAPPROVED              APPLICATIONSTATUS = "APPROVED"
	APPLICATIONSTATUSCONDITIONALLYAPPROVED APPLICATIONSTATUS = "CONDITIONALLY_APPROVED"
	APPLICATIONSTATUSREJECTED              APPLICATIONSTATUS = "REJECTED"
	APPLICATIONSTATUSWITHDRAWN             APPLICATIONSTATUS = "WITHDRAWN"
)

func (e *APPLICATIONSTATUS) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = APPLICATIONSTATUS(s)
	case string:
		*e = APPLICATIONSTATUS(s)
	default:
		return fmt.Errorf("unsupported scan type for APPLICATIONSTATUS: %T", src)
	}
	return nil
}

type NullAPPLICATIONSTATUS struct {
	APPLICATIONSTATUS APPLICATIONSTATUS `json:"APPLICATION_STATUS"`
	Valid             bool              `json:"valid"` // Valid is true if APPLICATIONSTATUS is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAPPLICATIONSTATUS) Scan(value interface{}) error {
	if value == nil {
		ns.APPLICATIONSTATUS, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.APPLICATIONSTATUS.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAPPLICATIONSTATUS) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.APPLICATIONSTATUS), nil
}

type CONTRACTSTATUS string

const (
	CONTRACTSTATUSPENDINGA CONTRACTSTATUS = "PENDING_A"
	CONTRACTSTATUSPENDINGB CONTRACTSTATUS = "PENDING_B"
	CONTRACTSTATUSPENDING  CONTRACTSTATUS = "PENDING"
	CONTRACTSTATUSSIGNED   CONTRACTSTATUS = "SIGNED"
	CONTRACTSTATUSREJECTED CONTRACTSTATUS = "REJECTED"
	CONTRACTSTATUSCANCELED CONTRACTSTATUS = "CANCELED"
)

func (e *CONTRACTSTATUS) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CONTRACTSTATUS(s)
	case string:
		*e = CONTRACTSTATUS(s)
	default:
		return fmt.Errorf("unsupported scan type for CONTRACTSTATUS: %T", src)
	}
	return nil
}

type NullCONTRACTSTATUS struct {
	CONTRACTSTATUS CONTRACTSTATUS `json:"CONTRACTSTATUS"`
	Valid          bool           `json:"valid"` // Valid is true if CONTRACTSTATUS is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCONTRACTSTATUS) Scan(value interface{}) error {
	if value == nil {
		ns.CONTRACTSTATUS, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CONTRACTSTATUS.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCONTRACTSTATUS) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CONTRACTSTATUS), nil
}

type MEDIATYPE string

const (
	MEDIATYPEIMAGE MEDIATYPE = "IMAGE"
	MEDIATYPEVIDEO MEDIATYPE = "VIDEO"
)

func (e *MEDIATYPE) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MEDIATYPE(s)
	case string:
		*e = MEDIATYPE(s)
	default:
		return fmt.Errorf("unsupported scan type for MEDIATYPE: %T", src)
	}
	return nil
}

type NullMEDIATYPE struct {
	MEDIATYPE MEDIATYPE `json:"MEDIATYPE"`
	Valid     bool      `json:"valid"` // Valid is true if MEDIATYPE is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMEDIATYPE) Scan(value interface{}) error {
	if value == nil {
		ns.MEDIATYPE, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MEDIATYPE.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMEDIATYPE) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MEDIATYPE), nil
}

type MESSAGESTATUS string

const (
	MESSAGESTATUSACTIVE  MESSAGESTATUS = "ACTIVE"
	MESSAGESTATUSDELETED MESSAGESTATUS = "DELETED"
)

func (e *MESSAGESTATUS) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MESSAGESTATUS(s)
	case string:
		*e = MESSAGESTATUS(s)
	default:
		return fmt.Errorf("unsupported scan type for MESSAGESTATUS: %T", src)
	}
	return nil
}

type NullMESSAGESTATUS struct {
	MESSAGESTATUS MESSAGESTATUS `json:"MESSAGESTATUS"`
	Valid         bool          `json:"valid"` // Valid is true if MESSAGESTATUS is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMESSAGESTATUS) Scan(value interface{}) error {
	if value == nil {
		ns.MESSAGESTATUS, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MESSAGESTATUS.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMESSAGESTATUS) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MESSAGESTATUS), nil
}

type MESSAGETYPE string

const (
	MESSAGETYPETEXT  MESSAGETYPE = "TEXT"
	MESSAGETYPEIMAGE MESSAGETYPE = "IMAGE"
	MESSAGETYPEFILE  MESSAGETYPE = "FILE"
)

func (e *MESSAGETYPE) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MESSAGETYPE(s)
	case string:
		*e = MESSAGETYPE(s)
	default:
		return fmt.Errorf("unsupported scan type for MESSAGETYPE: %T", src)
	}
	return nil
}

type NullMESSAGETYPE struct {
	MESSAGETYPE MESSAGETYPE `json:"MESSAGETYPE"`
	Valid       bool        `json:"valid"` // Valid is true if MESSAGETYPE is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMESSAGETYPE) Scan(value interface{}) error {
	if value == nil {
		ns.MESSAGETYPE, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MESSAGETYPE.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMESSAGETYPE) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MESSAGETYPE), nil
}

type PAYMENTSTATUS string

const (
	PAYMENTSTATUSPENDING PAYMENTSTATUS = "PENDING"
	PAYMENTSTATUSSUCCESS PAYMENTSTATUS = "SUCCESS"
	PAYMENTSTATUSFAILED  PAYMENTSTATUS = "FAILED"
)

func (e *PAYMENTSTATUS) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PAYMENTSTATUS(s)
	case string:
		*e = PAYMENTSTATUS(s)
	default:
		return fmt.Errorf("unsupported scan type for PAYMENTSTATUS: %T", src)
	}
	return nil
}

type NullPAYMENTSTATUS struct {
	PAYMENTSTATUS PAYMENTSTATUS `json:"PAYMENTSTATUS"`
	Valid         bool          `json:"valid"` // Valid is true if PAYMENTSTATUS is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPAYMENTSTATUS) Scan(value interface{}) error {
	if value == nil {
		ns.PAYMENTSTATUS, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PAYMENTSTATUS.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPAYMENTSTATUS) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PAYMENTSTATUS), nil
}

type PROPERTYTYPE string

const (
	PROPERTYTYPEAPARTMENT     PROPERTYTYPE = "APARTMENT"
	PROPERTYTYPEPRIVATE       PROPERTYTYPE = "PRIVATE"
	PROPERTYTYPEROOM          PROPERTYTYPE = "ROOM"
	PROPERTYTYPESTORE         PROPERTYTYPE = "STORE"
	PROPERTYTYPEOFFICE        PROPERTYTYPE = "OFFICE"
	PROPERTYTYPEVILLA         PROPERTYTYPE = "VILLA"
	PROPERTYTYPEMINIAPARTMENT PROPERTYTYPE = "MINIAPARTMENT"
)

func (e *PROPERTYTYPE) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PROPERTYTYPE(s)
	case string:
		*e = PROPERTYTYPE(s)
	default:
		return fmt.Errorf("unsupported scan type for PROPERTYTYPE: %T", src)
	}
	return nil
}

type NullPROPERTYTYPE struct {
	PROPERTYTYPE PROPERTYTYPE `json:"PROPERTYTYPE"`
	Valid        bool         `json:"valid"` // Valid is true if PROPERTYTYPE is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPROPERTYTYPE) Scan(value interface{}) error {
	if value == nil {
		ns.PROPERTYTYPE, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PROPERTYTYPE.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPROPERTYTYPE) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PROPERTYTYPE), nil
}

type REMINDERRECURRENCEMODE string

const (
	REMINDERRECURRENCEMODENONE    REMINDERRECURRENCEMODE = "NONE"
	REMINDERRECURRENCEMODEWEEKLY  REMINDERRECURRENCEMODE = "WEEKLY"
	REMINDERRECURRENCEMODEMONTHLY REMINDERRECURRENCEMODE = "MONTHLY"
)

func (e *REMINDERRECURRENCEMODE) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = REMINDERRECURRENCEMODE(s)
	case string:
		*e = REMINDERRECURRENCEMODE(s)
	default:
		return fmt.Errorf("unsupported scan type for REMINDERRECURRENCEMODE: %T", src)
	}
	return nil
}

type NullREMINDERRECURRENCEMODE struct {
	REMINDERRECURRENCEMODE REMINDERRECURRENCEMODE `json:"REMINDERRECURRENCEMODE"`
	Valid                  bool                   `json:"valid"` // Valid is true if REMINDERRECURRENCEMODE is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullREMINDERRECURRENCEMODE) Scan(value interface{}) error {
	if value == nil {
		ns.REMINDERRECURRENCEMODE, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.REMINDERRECURRENCEMODE.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullREMINDERRECURRENCEMODE) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.REMINDERRECURRENCEMODE), nil
}

type RENTALCOMPLAINTSTATUS string

const (
	RENTALCOMPLAINTSTATUSPENDING  RENTALCOMPLAINTSTATUS = "PENDING"
	RENTALCOMPLAINTSTATUSRESOLVED RENTALCOMPLAINTSTATUS = "RESOLVED"
	RENTALCOMPLAINTSTATUSCLOSED   RENTALCOMPLAINTSTATUS = "CLOSED"
)

func (e *RENTALCOMPLAINTSTATUS) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RENTALCOMPLAINTSTATUS(s)
	case string:
		*e = RENTALCOMPLAINTSTATUS(s)
	default:
		return fmt.Errorf("unsupported scan type for RENTALCOMPLAINTSTATUS: %T", src)
	}
	return nil
}

type NullRENTALCOMPLAINTSTATUS struct {
	RENTALCOMPLAINTSTATUS RENTALCOMPLAINTSTATUS `json:"RENTALCOMPLAINTSTATUS"`
	Valid                 bool                  `json:"valid"` // Valid is true if RENTALCOMPLAINTSTATUS is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRENTALCOMPLAINTSTATUS) Scan(value interface{}) error {
	if value == nil {
		ns.RENTALCOMPLAINTSTATUS, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RENTALCOMPLAINTSTATUS.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRENTALCOMPLAINTSTATUS) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RENTALCOMPLAINTSTATUS), nil
}

type RENTALCOMPLAINTTYPE string

const (
	RENTALCOMPLAINTTYPEREPORT     RENTALCOMPLAINTTYPE = "REPORT"
	RENTALCOMPLAINTTYPESUGGESTION RENTALCOMPLAINTTYPE = "SUGGESTION"
)

func (e *RENTALCOMPLAINTTYPE) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RENTALCOMPLAINTTYPE(s)
	case string:
		*e = RENTALCOMPLAINTTYPE(s)
	default:
		return fmt.Errorf("unsupported scan type for RENTALCOMPLAINTTYPE: %T", src)
	}
	return nil
}

type NullRENTALCOMPLAINTTYPE struct {
	RENTALCOMPLAINTTYPE RENTALCOMPLAINTTYPE `json:"RENTALCOMPLAINTTYPE"`
	Valid               bool                `json:"valid"` // Valid is true if RENTALCOMPLAINTTYPE is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRENTALCOMPLAINTTYPE) Scan(value interface{}) error {
	if value == nil {
		ns.RENTALCOMPLAINTTYPE, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RENTALCOMPLAINTTYPE.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRENTALCOMPLAINTTYPE) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RENTALCOMPLAINTTYPE), nil
}

type RENTALPAYMENTSTATUS string

const (
	RENTALPAYMENTSTATUSPLAN        RENTALPAYMENTSTATUS = "PLAN"
	RENTALPAYMENTSTATUSISSUED      RENTALPAYMENTSTATUS = "ISSUED"
	RENTALPAYMENTSTATUSPENDING     RENTALPAYMENTSTATUS = "PENDING"
	RENTALPAYMENTSTATUSREQUEST2PAY RENTALPAYMENTSTATUS = "REQUEST2PAY"
	RENTALPAYMENTSTATUSPAID        RENTALPAYMENTSTATUS = "PAID"
	RENTALPAYMENTSTATUSCANCELLED   RENTALPAYMENTSTATUS = "CANCELLED"
)

func (e *RENTALPAYMENTSTATUS) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RENTALPAYMENTSTATUS(s)
	case string:
		*e = RENTALPAYMENTSTATUS(s)
	default:
		return fmt.Errorf("unsupported scan type for RENTALPAYMENTSTATUS: %T", src)
	}
	return nil
}

type NullRENTALPAYMENTSTATUS struct {
	RENTALPAYMENTSTATUS RENTALPAYMENTSTATUS `json:"RENTALPAYMENTSTATUS"`
	Valid               bool                `json:"valid"` // Valid is true if RENTALPAYMENTSTATUS is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRENTALPAYMENTSTATUS) Scan(value interface{}) error {
	if value == nil {
		ns.RENTALPAYMENTSTATUS, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RENTALPAYMENTSTATUS.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRENTALPAYMENTSTATUS) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RENTALPAYMENTSTATUS), nil
}

type RENTALPAYMENTTYPE string

const (
	RENTALPAYMENTTYPEPREPAID  RENTALPAYMENTTYPE = "PREPAID"
	RENTALPAYMENTTYPEPOSTPAID RENTALPAYMENTTYPE = "POSTPAID"
)

func (e *RENTALPAYMENTTYPE) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RENTALPAYMENTTYPE(s)
	case string:
		*e = RENTALPAYMENTTYPE(s)
	default:
		return fmt.Errorf("unsupported scan type for RENTALPAYMENTTYPE: %T", src)
	}
	return nil
}

type NullRENTALPAYMENTTYPE struct {
	RENTALPAYMENTTYPE RENTALPAYMENTTYPE `json:"RENTALPAYMENTTYPE"`
	Valid             bool              `json:"valid"` // Valid is true if RENTALPAYMENTTYPE is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRENTALPAYMENTTYPE) Scan(value interface{}) error {
	if value == nil {
		ns.RENTALPAYMENTTYPE, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RENTALPAYMENTTYPE.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRENTALPAYMENTTYPE) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RENTALPAYMENTTYPE), nil
}

type RENTALSTATUS string

const (
	RENTALSTATUSINPROGRESS RENTALSTATUS = "INPROGRESS"
	RENTALSTATUSEND        RENTALSTATUS = "END"
)

func (e *RENTALSTATUS) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RENTALSTATUS(s)
	case string:
		*e = RENTALSTATUS(s)
	default:
		return fmt.Errorf("unsupported scan type for RENTALSTATUS: %T", src)
	}
	return nil
}

type NullRENTALSTATUS struct {
	RENTALSTATUS RENTALSTATUS `json:"RENTALSTATUS"`
	Valid        bool         `json:"valid"` // Valid is true if RENTALSTATUS is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRENTALSTATUS) Scan(value interface{}) error {
	if value == nil {
		ns.RENTALSTATUS, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RENTALSTATUS.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRENTALSTATUS) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RENTALSTATUS), nil
}

type TENANTTYPE string

const (
	TENANTTYPEINDIVIDUAL   TENANTTYPE = "INDIVIDUAL"
	TENANTTYPEFAMILY       TENANTTYPE = "FAMILY"
	TENANTTYPEORGANIZATION TENANTTYPE = "ORGANIZATION"
)

func (e *TENANTTYPE) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TENANTTYPE(s)
	case string:
		*e = TENANTTYPE(s)
	default:
		return fmt.Errorf("unsupported scan type for TENANTTYPE: %T", src)
	}
	return nil
}

type NullTENANTTYPE struct {
	TENANTTYPE TENANTTYPE `json:"TENANTTYPE"`
	Valid      bool       `json:"valid"` // Valid is true if TENANTTYPE is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTENANTTYPE) Scan(value interface{}) error {
	if value == nil {
		ns.TENANTTYPE, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TENANTTYPE.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTENANTTYPE) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TENANTTYPE), nil
}

type UNITTYPE string

const (
	UNITTYPEROOM      UNITTYPE = "ROOM"
	UNITTYPEAPARTMENT UNITTYPE = "APARTMENT"
	UNITTYPESTUDIO    UNITTYPE = "STUDIO"
)

func (e *UNITTYPE) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UNITTYPE(s)
	case string:
		*e = UNITTYPE(s)
	default:
		return fmt.Errorf("unsupported scan type for UNITTYPE: %T", src)
	}
	return nil
}

type NullUNITTYPE struct {
	UNITTYPE UNITTYPE `json:"UNITTYPE"`
	Valid    bool     `json:"valid"` // Valid is true if UNITTYPE is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUNITTYPE) Scan(value interface{}) error {
	if value == nil {
		ns.UNITTYPE, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UNITTYPE.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUNITTYPE) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UNITTYPE), nil
}

type USERROLE string

const (
	USERROLEADMIN    USERROLE = "ADMIN"
	USERROLELANDLORD USERROLE = "LANDLORD"
	USERROLETENANT   USERROLE = "TENANT"
)

func (e *USERROLE) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = USERROLE(s)
	case string:
		*e = USERROLE(s)
	default:
		return fmt.Errorf("unsupported scan type for USERROLE: %T", src)
	}
	return nil
}

type NullUSERROLE struct {
	USERROLE USERROLE `json:"USERROLE"`
	Valid    bool     `json:"valid"` // Valid is true if USERROLE is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUSERROLE) Scan(value interface{}) error {
	if value == nil {
		ns.USERROLE, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.USERROLE.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUSERROLE) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.USERROLE), nil
}

type Account struct {
	ID                uuid.UUID   `json:"id"`
	UserId            uuid.UUID   `json:"userId"`
	Type              string      `json:"type"`
	Provider          string      `json:"provider"`
	ProviderAccountId string      `json:"providerAccountId"`
	RefreshToken      pgtype.Text `json:"refresh_token"`
	AccessToken       pgtype.Text `json:"access_token"`
	ExpiresAt         pgtype.Int4 `json:"expires_at"`
	TokenType         pgtype.Text `json:"token_type"`
	Scope             pgtype.Text `json:"scope"`
	IDToken           pgtype.Text `json:"id_token"`
	SessionState      pgtype.Text `json:"session_state"`
}

type Application struct {
	ID                      int64             `json:"id"`
	CreatorID               pgtype.UUID       `json:"creator_id"`
	ListingID               uuid.UUID         `json:"listing_id"`
	PropertyID              uuid.UUID         `json:"property_id"`
	UnitID                  uuid.UUID         `json:"unit_id"`
	ListingPrice            float32           `json:"listing_price"`
	OfferedPrice            float32           `json:"offered_price"`
	Status                  APPLICATIONSTATUS `json:"status"`
	CreatedAt               time.Time         `json:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at"`
	TenantType              TENANTTYPE        `json:"tenant_type"`
	FullName                string            `json:"full_name"`
	Email                   string            `json:"email"`
	Phone                   string            `json:"phone"`
	Dob                     pgtype.Date       `json:"dob"`
	ProfileImage            string            `json:"profile_image"`
	MoveinDate              pgtype.Date       `json:"movein_date"`
	PreferredTerm           int32             `json:"preferred_term"`
	RentalIntention         string            `json:"rental_intention"`
	OrganizationName        pgtype.Text       `json:"organization_name"`
	OrganizationHqAddress   pgtype.Text       `json:"organization_hq_address"`
	OrganizationScale       pgtype.Text       `json:"organization_scale"`
	RhAddress               pgtype.Text       `json:"rh_address"`
	RhCity                  pgtype.Text       `json:"rh_city"`
	RhDistrict              pgtype.Text       `json:"rh_district"`
	RhWard                  pgtype.Text       `json:"rh_ward"`
	RhRentalDuration        pgtype.Int4       `json:"rh_rental_duration"`
	RhMonthlyPayment        pgtype.Int8       `json:"rh_monthly_payment"`
	RhReasonForLeaving      pgtype.Text       `json:"rh_reason_for_leaving"`
	EmploymentStatus        string            `json:"employment_status"`
	EmploymentCompanyName   pgtype.Text       `json:"employment_company_name"`
	EmploymentPosition      pgtype.Text       `json:"employment_position"`
	EmploymentMonthlyIncome pgtype.Int8       `json:"employment_monthly_income"`
	EmploymentComment       pgtype.Text       `json:"employment_comment"`
}

type ApplicationCoap struct {
	ApplicationID int64       `json:"application_id"`
	FullName      string      `json:"full_name"`
	Dob           time.Time   `json:"dob"`
	Job           string      `json:"job"`
	Income        int32       `json:"income"`
	Email         pgtype.Text `json:"email"`
	Phone         pgtype.Text `json:"phone"`
	Description   pgtype.Text `json:"description"`
}

type ApplicationMinor struct {
	ApplicationID int64       `json:"application_id"`
	FullName      string      `json:"full_name"`
	Dob           time.Time   `json:"dob"`
	Email         pgtype.Text `json:"email"`
	Phone         pgtype.Text `json:"phone"`
	Description   pgtype.Text `json:"description"`
}

type ApplicationPet struct {
	ApplicationID int64         `json:"application_id"`
	Type          string        `json:"type"`
	Weight        pgtype.Float4 `json:"weight"`
	Description   pgtype.Text   `json:"description"`
}

type ApplicationVehicle struct {
	ApplicationID int64       `json:"application_id"`
	Type          string      `json:"type"`
	Model         pgtype.Text `json:"model"`
	Code          string      `json:"code"`
	Description   pgtype.Text `json:"description"`
}

type Contract struct {
	ID                        int64          `json:"id"`
	RentalID                  int64          `json:"rental_id"`
	AFullname                 string         `json:"a_fullname"`
	ADob                      pgtype.Date    `json:"a_dob"`
	APhone                    string         `json:"a_phone"`
	AAddress                  string         `json:"a_address"`
	AHouseholdRegistration    string         `json:"a_household_registration"`
	AIdentity                 string         `json:"a_identity"`
	AIdentityIssuedBy         string         `json:"a_identity_issued_by"`
	AIdentityIssuedAt         pgtype.Date    `json:"a_identity_issued_at"`
	ADocuments                []string       `json:"a_documents"`
	ABankAccount              pgtype.Text    `json:"a_bank_account"`
	ABank                     pgtype.Text    `json:"a_bank"`
	ARegistrationNumber       string         `json:"a_registration_number"`
	BFullname                 string         `json:"b_fullname"`
	BOrganizationName         pgtype.Text    `json:"b_organization_name"`
	BOrganizationHqAddress    pgtype.Text    `json:"b_organization_hq_address"`
	BOrganizationCode         pgtype.Text    `json:"b_organization_code"`
	BOrganizationCodeIssuedAt pgtype.Date    `json:"b_organization_code_issued_at"`
	BOrganizationCodeIssuedBy pgtype.Text    `json:"b_organization_code_issued_by"`
	BDob                      pgtype.Text    `json:"b_dob"`
	BPhone                    string         `json:"b_phone"`
	BAddress                  pgtype.Text    `json:"b_address"`
	BHouseholdRegistration    pgtype.Text    `json:"b_household_registration"`
	BIdentity                 pgtype.Text    `json:"b_identity"`
	BIdentityIssuedBy         pgtype.Text    `json:"b_identity_issued_by"`
	BIdentityIssuedAt         pgtype.Date    `json:"b_identity_issued_at"`
	BBankAccount              pgtype.Text    `json:"b_bank_account"`
	BBank                     pgtype.Text    `json:"b_bank"`
	BTaxCode                  pgtype.Text    `json:"b_tax_code"`
	PaymentMethod             string         `json:"payment_method"`
	PaymentDay                int32          `json:"payment_day"`
	NCopies                   int32          `json:"n_copies"`
	CreatedAtPlace            string         `json:"created_at_place"`
	Content                   string         `json:"content"`
	Status                    CONTRACTSTATUS `json:"status"`
	CreatedAt                 time.Time      `json:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at"`
	CreatedBy                 uuid.UUID      `json:"created_by"`
	UpdatedBy                 uuid.UUID      `json:"updated_by"`
}

type LPolicy struct {
	ID     int64  `json:"id"`
	Policy string `json:"policy"`
}

type Listing struct {
	ID          uuid.UUID `json:"id"`
	CreatorID   uuid.UUID `json:"creator_id"`
	PropertyID  uuid.UUID `json:"property_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	ContactType string    `json:"contact_type"`
	// Rental price per month in vietnamese dong
	Price           float32       `json:"price"`
	PriceNegotiable bool          `json:"price_negotiable"`
	SecurityDeposit pgtype.Float4 `json:"security_deposit"`
	// Lease term in months
	LeaseTerm         pgtype.Int4 `json:"lease_term"`
	PetsAllowed       pgtype.Bool `json:"pets_allowed"`
	NumberOfResidents pgtype.Int4 `json:"number_of_residents"`
	// Priority of the listing, range from 1 to 5, 1 is the lowest
	Priority  int32     `json:"priority"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// The time when the listing is expired. The listing is expired if the current time is greater than this time.
	ExpiredAt time.Time `json:"expired_at"`
}

type ListingPolicy struct {
	ListingID uuid.UUID   `json:"listing_id"`
	PolicyID  int64       `json:"policy_id"`
	Note      pgtype.Text `json:"note"`
}

type ListingTag struct {
	ID        int64     `json:"id"`
	ListingID uuid.UUID `json:"listing_id"`
	Tag       string    `json:"tag"`
}

type ListingUnit struct {
	ListingID uuid.UUID `json:"listing_id"`
	UnitID    uuid.UUID `json:"unit_id"`
	Price     int64     `json:"price"`
}

type Message struct {
	ID        int64         `json:"id"`
	GroupID   int64         `json:"group_id"`
	FromUser  uuid.UUID     `json:"from_user"`
	Content   string        `json:"content"`
	Status    MESSAGESTATUS `json:"status"`
	Type      MESSAGETYPE   `json:"type"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type MsgGroup struct {
	GroupID   int64     `json:"group_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
}

type MsgGroupMember struct {
	GroupID int64     `json:"group_id"`
	UserID  uuid.UUID `json:"user_id"`
}

type NewPropertyManagerRequest struct {
	ID         int64       `json:"id"`
	CreatorID  uuid.UUID   `json:"creator_id"`
	PropertyID uuid.UUID   `json:"property_id"`
	UserID     pgtype.UUID `json:"user_id"`
	Email      string      `json:"email"`
	Approved   bool        `json:"approved"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// Security guard, Parking, Gym, ...
type PFeature struct {
	ID      int64  `json:"id"`
	Feature string `json:"feature"`
}

type Payment struct {
	ID        int64         `json:"id"`
	UserID    uuid.UUID     `json:"user_id"`
	OrderID   string        `json:"order_id"`
	OrderInfo string        `json:"order_info"`
	Amount    float32       `json:"amount"`
	Status    PAYMENTSTATUS `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type PaymentItem struct {
	PaymentID int64   `json:"payment_id"`
	Name      string  `json:"name"`
	Price     float32 `json:"price"`
	Quantity  int32   `json:"quantity"`
	Discount  int32   `json:"discount"`
}

type Property struct {
	ID             uuid.UUID   `json:"id"`
	CreatorID      uuid.UUID   `json:"creator_id"`
	Name           string      `json:"name"`
	Building       pgtype.Text `json:"building"`
	Project        pgtype.Text `json:"project"`
	Area           float32     `json:"area"`
	NumberOfFloors pgtype.Int4 `json:"number_of_floors"`
	YearBuilt      pgtype.Int4 `json:"year_built"`
	// n,s,w,e,nw,ne,sw,se
	Orientation   pgtype.Text   `json:"orientation"`
	EntranceWidth pgtype.Float4 `json:"entrance_width"`
	Facade        pgtype.Float4 `json:"facade"`
	FullAddress   string        `json:"full_address"`
	City          string        `json:"city"`
	District      string        `json:"district"`
	Ward          pgtype.Text   `json:"ward"`
	Lat           pgtype.Float8 `json:"lat"`
	Lng           pgtype.Float8 `json:"lng"`
	PrimaryImage  pgtype.Int8   `json:"primary_image"`
	Description   pgtype.Text   `json:"description"`
	Type          PROPERTYTYPE  `json:"type"`
	IsPublic      bool          `json:"is_public"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type PropertyFeature struct {
	PropertyID  uuid.UUID   `json:"property_id"`
	FeatureID   int64       `json:"feature_id"`
	Description pgtype.Text `json:"description"`
}

type PropertyManager struct {
	PropertyID uuid.UUID `json:"property_id"`
	ManagerID  uuid.UUID `json:"manager_id"`
	Role       string    `json:"role"`
}

type PropertyMedium struct {
	ID          int64       `json:"id"`
	PropertyID  uuid.UUID   `json:"property_id"`
	Url         string      `json:"url"`
	Type        MEDIATYPE   `json:"type"`
	Description pgtype.Text `json:"description"`
}

type PropertyTag struct {
	ID         int64     `json:"id"`
	PropertyID uuid.UUID `json:"property_id"`
	Tag        string    `json:"tag"`
}

type Reminder struct {
	ID        int64       `json:"id"`
	CreatorID uuid.UUID   `json:"creator_id"`
	Title     string      `json:"title"`
	StartAt   time.Time   `json:"start_at"`
	EndAt     time.Time   `json:"end_at"`
	Note      pgtype.Text `json:"note"`
	Location  string      `json:"location"`
	// 7-bit integer representing days in a week (0-6) when the reminder should be triggered. 0 is Sunday, 1 is Monday, and so on.
	RecurrenceDay pgtype.Int4 `json:"recurrence_day"`
	// 32-bit integer representing days in a month (0-30) when the reminder should be triggered. 0 is the last day of the month, 1 is the first day of the month, and so on.
	RecurrenceMonth pgtype.Int4            `json:"recurrence_month"`
	RecurrenceMode  REMINDERRECURRENCEMODE `json:"recurrence_mode"`
	Priority        int32                  `json:"priority"`
	ResourceTag     string                 `json:"resource_tag"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type Rental struct {
	ID                      int64             `json:"id"`
	CreatorID               uuid.UUID         `json:"creator_id"`
	PropertyID              uuid.UUID         `json:"property_id"`
	UnitID                  uuid.UUID         `json:"unit_id"`
	ApplicationID           pgtype.Int8       `json:"application_id"`
	TenantID                pgtype.UUID       `json:"tenant_id"`
	ProfileImage            string            `json:"profile_image"`
	TenantType              TENANTTYPE        `json:"tenant_type"`
	TenantName              string            `json:"tenant_name"`
	TenantPhone             string            `json:"tenant_phone"`
	TenantEmail             string            `json:"tenant_email"`
	OrganizationName        pgtype.Text       `json:"organization_name"`
	OrganizationHqAddress   pgtype.Text       `json:"organization_hq_address"`
	StartDate               pgtype.Date       `json:"start_date"`
	MoveinDate              pgtype.Date       `json:"movein_date"`
	RentalPeriod            int32             `json:"rental_period"`
	PaymentType             RENTALPAYMENTTYPE `json:"payment_type"`
	RentalPrice             float32           `json:"rental_price"`
	RentalPaymentBasis      int32             `json:"rental_payment_basis"`
	RentalIntention         string            `json:"rental_intention"`
	Deposit                 float32           `json:"deposit"`
	DepositPaid             bool              `json:"deposit_paid"`
	NoticePeriod            pgtype.Int4       `json:"notice_period"`
	ElectricitySetupBy      string            `json:"electricity_setup_by"`
	ElectricityPaymentType  pgtype.Text       `json:"electricity_payment_type"`
	ElectricityCustomerCode pgtype.Text       `json:"electricity_customer_code"`
	ElectricityProvider     pgtype.Text       `json:"electricity_provider"`
	ElectricityPrice        pgtype.Float4     `json:"electricity_price"`
	WaterSetupBy            string            `json:"water_setup_by"`
	WaterPaymentType        pgtype.Text       `json:"water_payment_type"`
	WaterCustomerCode       pgtype.Text       `json:"water_customer_code"`
	WaterProvider           pgtype.Text       `json:"water_provider"`
	WaterPrice              pgtype.Float4     `json:"water_price"`
	Note                    pgtype.Text       `json:"note"`
	Status                  RENTALSTATUS      `json:"status"`
	CreatedAt               time.Time         `json:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at"`
}

type RentalCoap struct {
	RentalID    int64       `json:"rental_id"`
	FullName    pgtype.Text `json:"full_name"`
	Dob         pgtype.Date `json:"dob"`
	Job         pgtype.Text `json:"job"`
	Income      pgtype.Int4 `json:"income"`
	Email       pgtype.Text `json:"email"`
	Phone       pgtype.Text `json:"phone"`
	Description pgtype.Text `json:"description"`
}

type RentalComplaint struct {
	ID         int64                 `json:"id"`
	RentalID   int64                 `json:"rental_id"`
	CreatorID  uuid.UUID             `json:"creator_id"`
	Title      string                `json:"title"`
	Content    string                `json:"content"`
	Suggestion pgtype.Text           `json:"suggestion"`
	Media      []string              `json:"media"`
	OccurredAt time.Time             `json:"occurred_at"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
	UpdatedBy  uuid.UUID             `json:"updated_by"`
	Type       RENTALCOMPLAINTTYPE   `json:"type"`
	Status     RENTALCOMPLAINTSTATUS `json:"status"`
}

type RentalComplaintReply struct {
	ComplaintID int64     `json:"complaint_id"`
	ReplierID   uuid.UUID `json:"replier_id"`
	Reply       string    `json:"reply"`
	Media       []string  `json:"media"`
	CreatedAt   time.Time `json:"created_at"`
}

type RentalMinor struct {
	RentalID    int64       `json:"rental_id"`
	FullName    string      `json:"full_name"`
	Dob         pgtype.Date `json:"dob"`
	Email       pgtype.Text `json:"email"`
	Phone       pgtype.Text `json:"phone"`
	Description pgtype.Text `json:"description"`
}

type RentalPayment struct {
	ID int64 `json:"id"`
	// {payment.id}_{ELECTRICITY | WATER | RENTAL | DEPOSIT | SERVICES{id}}_{payment.created_at}
	Code       string      `json:"code"`
	RentalID   int64       `json:"rental_id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	StartDate  pgtype.Date `json:"start_date"`
	EndDate    pgtype.Date `json:"end_date"`
	ExpiryDate pgtype.Date `json:"expiry_date"`
	// the date the payment gets paid
	PaymentDate pgtype.Date         `json:"payment_date"`
	UpdatedBy   pgtype.UUID         `json:"updated_by"`
	Status      RENTALPAYMENTSTATUS `json:"status"`
	Amount      float32             `json:"amount"`
	Discount    pgtype.Float4       `json:"discount"`
	Penalty     pgtype.Float4       `json:"penalty"`
	Note        pgtype.Text         `json:"note"`
}

type RentalPet struct {
	RentalID    int64         `json:"rental_id"`
	Type        string        `json:"type"`
	Weight      pgtype.Float4 `json:"weight"`
	Description pgtype.Text   `json:"description"`
}

type RentalPolicy struct {
	RentalID int64  `json:"rental_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type RentalService struct {
	ID       int64  `json:"id"`
	RentalID int64  `json:"rental_id"`
	Name     string `json:"name"`
	// The party who set up the service, either "LANDLORD" or "TENANT"
	SetupBy  string        `json:"setup_by"`
	Provider pgtype.Text   `json:"provider"`
	Price    pgtype.Float4 `json:"price"`
}

type Session struct {
	ID           uuid.UUID   `json:"id"`
	SessionToken string      `json:"sessionToken"`
	UserId       uuid.UUID   `json:"userId"`
	Expires      time.Time   `json:"expires"`
	UserAgent    pgtype.Text `json:"user_agent"`
	ClientIp     pgtype.Text `json:"client_ip"`
	IsBlocked    bool        `json:"is_blocked"`
	CreatedAt    time.Time   `json:"created_at"`
}

// Air conditioner, Fridge, Washing machine, ...
type UAmenity struct {
	ID      int64  `json:"id"`
	Amenity string `json:"amenity"`
}

type Unit struct {
	ID                  uuid.UUID   `json:"id"`
	PropertyID          uuid.UUID   `json:"property_id"`
	Name                string      `json:"name"`
	Area                float32     `json:"area"`
	Floor               pgtype.Int4 `json:"floor"`
	NumberOfLivingRooms pgtype.Int4 `json:"number_of_living_rooms"`
	NumberOfBedrooms    pgtype.Int4 `json:"number_of_bedrooms"`
	NumberOfBathrooms   pgtype.Int4 `json:"number_of_bathrooms"`
	NumberOfToilets     pgtype.Int4 `json:"number_of_toilets"`
	NumberOfBalconies   pgtype.Int4 `json:"number_of_balconies"`
	NumberOfKitchens    pgtype.Int4 `json:"number_of_kitchens"`
	Type                UNITTYPE    `json:"type"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

type UnitAmenity struct {
	UnitID      uuid.UUID   `json:"unit_id"`
	AmenityID   int64       `json:"amenity_id"`
	Description pgtype.Text `json:"description"`
}

type UnitMedium struct {
	ID          int64       `json:"id"`
	UnitID      uuid.UUID   `json:"unit_id"`
	Url         string      `json:"url"`
	Type        MEDIATYPE   `json:"type"`
	Description pgtype.Text `json:"description"`
}

// User info table
type User struct {
	ID        uuid.UUID   `json:"id"`
	Email     string      `json:"email"`
	Password  pgtype.Text `json:"password"`
	GroupID   pgtype.UUID `json:"group_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	CreatedBy pgtype.UUID `json:"created_by"`
	UpdatedBy pgtype.UUID `json:"updated_by"`
	// 1: deleted, 0: not deleted
	DeletedF  bool        `json:"deleted_f"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Phone     pgtype.Text `json:"phone"`
	Avatar    pgtype.Text `json:"avatar"`
	Address   pgtype.Text `json:"address"`
	City      pgtype.Text `json:"city"`
	District  pgtype.Text `json:"district"`
	Ward      pgtype.Text `json:"ward"`
	Role      USERROLE    `json:"role"`
}

type VerificationToken struct {
	Identifier string    `json:"identifier"`
	Token      string    `json:"token"`
	Expires    time.Time `json:"expires"`
}
