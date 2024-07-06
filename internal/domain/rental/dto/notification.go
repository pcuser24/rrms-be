package dto

import (
	"github.com/google/uuid"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type NotifyCreatePreRental struct {
	Rental *rental_model.RentalModel `json:"rental"`
	Secret string                    `json:"secret"`
}

type NotifyUpdatePreRental struct {
	PreRental  *rental_model.PreRental   `json:"preRental"`
	Rental     *rental_model.RentalModel `json:"rental"`
	UpdateData *UpdatePreRental          `json:"updateData"`
}

type NotifyUpdatePayments struct {
	Rental *rental_model.RentalModel `json:"rental"`
	// old rental payment data before update
	RentalPayment *rental_model.RentalPayment `json:"rentalPayment"`
	UpdateData    *UpdateRentalPayment        `json:"updateData"`
}

type NotifyCreateRentalPayment struct {
	Rental        *rental_model.RentalModel   `json:"rental"`
	RentalPayment *rental_model.RentalPayment `json:"rentalPayment"`
}

type NotifyCreateContract struct {
	Contract *rental_model.ContractModel `json:"contract"`
	Rental   *rental_model.RentalModel   `json:"rental"`
}

type NotifyUpdateContract struct {
	Contract *rental_model.ContractModel `json:"contract"`
	Rental   *rental_model.RentalModel   `json:"rental"`
	Side     string                      `json:"side"`
}

type NotifyCreateRentalComplaint struct {
	Complaint *rental_model.RentalComplaint `json:"complaint"`
	Rental    *rental_model.RentalModel     `json:"rental"`
}

type NotifyCreateComplaintReply struct {
	Complaint      *rental_model.RentalComplaint      `json:"complaint"`
	ComplaintReply *rental_model.RentalComplaintReply `json:"complaintReply"`
	Rental         *rental_model.RentalModel          `json:"rental"`
}

type NotifyUpdateComplaintStatus struct {
	Complaint *rental_model.RentalComplaint  `json:"complaint"`
	Rental    *rental_model.RentalModel      `json:"rental"`
	Status    database.RENTALCOMPLAINTSTATUS `json:"status"`
	UpdatedBy uuid.UUID                      `json:"updatedBy"`
}
