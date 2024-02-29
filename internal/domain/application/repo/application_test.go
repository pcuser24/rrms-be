package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/application/dto"
	"github.com/user2410/rrms-backend/internal/domain/application/model"
	"github.com/user2410/rrms-backend/internal/utils/random"
)

func TestCreateApplication(t *testing.T) {
	NewRandomApplicationDB(t, testAuthRepo, testListingRepo, testPropertyRepo, testUnitRepo, testApplicationRepo)
}

func TestGetApplicationById(t *testing.T) {
	a := NewRandomApplicationDB(t, testAuthRepo, testListingRepo, testPropertyRepo, testUnitRepo, testApplicationRepo)

	a_1, err := testApplicationRepo.GetApplicationById(context.Background(), a.ID)
	require.NoError(t, err)
	sameApplications(t, a, a_1)
}

func TestGetApplicationByIds(t *testing.T) {
	selectedFields := random.RandomlyPickNFromSlice[string](dto.GetRetrievableFields(), 5)
	a1 := NewRandomApplicationDB(t, testAuthRepo, testListingRepo, testPropertyRepo, testUnitRepo, testApplicationRepo)
	a2 := NewRandomApplicationDB(t, testAuthRepo, testListingRepo, testPropertyRepo, testUnitRepo, testApplicationRepo)

	as, err := testApplicationRepo.GetApplicationsByIds(
		context.Background(),
		[]int64{a1.ID, a2.ID},
		selectedFields,
	)
	require.NoError(t, err)
	require.Equal(t, len(as), 2)

	compareFn := func(p_1, p_2 *model.ApplicationModel) {
		for _, f := range selectedFields {
			switch f {
			case "creator_id":
				require.Equal(t, p_1.CreatorID, p_2.CreatorID)
			case "listing_id":
				require.Equal(t, p_1.ListingID, p_2.ListingID)
			case "property_id":
				require.Equal(t, p_1.PropertyID, p_2.PropertyID)
			case "status":
				require.Equal(t, p_1.Status, p_2.Status)
			case "created_at":
				require.Equal(t, p_1.CreatedAt, p_2.CreatedAt)
			case "updated_at":
				require.Equal(t, p_1.UpdatedAt, p_2.UpdatedAt)
			case "full_name":
				require.Equal(t, p_1.FullName, p_2.FullName)
			case "email":
				require.Equal(t, p_1.Email, p_2.Email)
			case "phone":
				require.Equal(t, p_1.Phone, p_2.Phone)
			case "dob":
				require.Equal(t, p_1.Dob, p_2.Dob)
			case "profile_image":
				require.Equal(t, p_1.ProfileImage, p_2.ProfileImage)
			case "movein_date":
				require.Equal(t, p_1.MoveinDate, p_2.MoveinDate)
			case "preferred_term":
				require.Equal(t, p_1.PreferredTerm, p_2.PreferredTerm)
			case "rental_intention":
				require.Equal(t, p_1.RentalIntention, p_2.RentalIntention)
			case "rh_address":
				require.Equal(t, p_1.RhAddress, p_2.RhAddress)
			case "rh_city":
				require.Equal(t, p_1.RhCity, p_2.RhCity)
			case "rh_district":
				require.Equal(t, p_1.RhDistrict, p_2.RhDistrict)
			case "rh_ward":
				require.Equal(t, p_1.RhWard, p_2.RhWard)
			case "rh_rental_duration":
				require.Equal(t, p_1.RhRentalDuration, p_2.RhRentalDuration)
			case "rh_monthly_payment":
				require.Equal(t, p_1.RhMonthlyPayment, p_2.RhMonthlyPayment)
			case "rh_reason_for_leaving":
				require.Equal(t, p_1.RhReasonForLeaving, p_2.RhReasonForLeaving)
			case "employment_status":
				require.Equal(t, p_1.EmploymentStatus, p_2.EmploymentStatus)
			case "employment_company_name":
				require.Equal(t, p_1.EmploymentCompanyName, p_2.EmploymentCompanyName)
			case "employment_position":
				require.Equal(t, p_1.EmploymentPosition, p_2.EmploymentPosition)
			case "employment_monthly_income":
				require.Equal(t, p_1.EmploymentMonthlyIncome, p_2.EmploymentMonthlyIncome)
			case "employment_comment":
				require.Equal(t, p_1.EmploymentComment, p_2.EmploymentComment)
			case "identity_type":
				require.Equal(t, p_1.IdentityType, p_2.IdentityType)
			case "identity_number":
				require.Equal(t, p_1.IdentityNumber, p_2.IdentityNumber)
			case "units":
				require.Equal(t, (p_1.Units), (p_2.Units))
			case "minors":
				require.Equal(t, (p_1.Minors), (p_2.Minors))
			case "coaps":
				require.Equal(t, (p_1.Coaps), (p_2.Coaps))
			case "pets":
				require.Equal(t, (p_1.Pets), (p_2.Pets))
			case "vehicles":
				require.Equal(t, (p_1.Vehicles), (p_2.Vehicles))
			}
		}
	}
	for _, a := range as {
		if a.ID == a1.ID {
			compareFn(&a, a1)
		} else if a.ID == a2.ID {
			compareFn(&a, a2)
		} else {
			t.Error("unexpected property", a)
		}
	}
}
