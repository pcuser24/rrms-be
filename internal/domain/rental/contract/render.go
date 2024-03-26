package contract

import (
	"fmt"
	"time"

	application_model "github.com/user2410/rrms-backend/internal/domain/application/model"
	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	html_util "github.com/user2410/rrms-backend/internal/utils/html"
	"github.com/user2410/rrms-backend/internal/utils/number"
)

const (
	shortPlaceHolder  = "................."
	mediumPlaceHolder = "..................................."
	longPlaceHolder   = "...................................................................."
)

var (
	pType2Template = map[database.PROPERTYTYPE]string{
		"APARTMENT":     "apartment_template.html",
		"PRIVATE":       "private_template.html",
		"ROOM":          "room_template.html",
		"STORE":         "private_template.html",
		"OFFICE":        "office_template.html",
		"MINIAPARTMENT": "apartment_template.html",
	}
)

var (
	basePath = utils.GetBasePath() + "/internal/domain/rental/contract/"
)

func RenderContractTemplate(
	prerental *model.PrerentalModel,
	application *application_model.ApplicationModel,
	property *property_model.PropertyModel,
	unit *unit_model.UnitModel,
	owner *auth_model.UserModel,
) (string, error) {
	templateFile := basePath + pType2Template[property.Type]
	data := struct {
		Date              html_util.HTMLTime
		OwnerName         string
		OwnerAddress      string
		OwnerPhone        string
		OwnerEmail        string
		TenantName        string
		TenantDOB         html_util.HTMLTime
		TenantIdentity    string
		TenantAddress     string
		TenantPhone       string
		TenantEmail       string
		MoveinDate        html_util.HTMLTime
		StartDate         html_util.HTMLTime
		EndDate           html_util.HTMLTime
		RentalPrice       string
		RentalPriceStr    string
		Deposit           string
		DepositStr        string
		RentalDuration    string
		RentalDurationStr string
		FullAddress       string
		NumberOfFloors    string
		PArea             string
		PBuilding         string
		PProject          string
		UName             string
		UArea             string
		UAreaStr          string
		UFloor            string
	}{
		Date:          html_util.NewHTMLTime(time.Now()),
		OwnerName:     mediumPlaceHolder,
		OwnerAddress:  mediumPlaceHolder,
		OwnerPhone:    mediumPlaceHolder,
		OwnerEmail:    mediumPlaceHolder,
		TenantName:    mediumPlaceHolder,
		TenantDOB:     html_util.NewHTMLTime(prerental.TenantDob),
		TenantAddress: mediumPlaceHolder,
		TenantPhone:   mediumPlaceHolder,
		TenantEmail:   mediumPlaceHolder,
		// MoveinDate:        mediumPlaceHolder,
		// StartDate:         mediumPlaceHolder,
		// EndDate:           mediumPlaceHolder,
		RentalPrice:       shortPlaceHolder,
		RentalPriceStr:    shortPlaceHolder,
		Deposit:           mediumPlaceHolder,
		DepositStr:        mediumPlaceHolder,
		RentalDuration:    mediumPlaceHolder,
		RentalDurationStr: mediumPlaceHolder,
		FullAddress:       property.FullAddress,
		NumberOfFloors:    mediumPlaceHolder,
		PArea:             fmt.Sprintf("%f", property.Area),
		PBuilding:         mediumPlaceHolder,
		PProject:          mediumPlaceHolder,
		UName:             unit.Name,
		UArea:             fmt.Sprintf("%f", unit.Area),
		UAreaStr:          shortPlaceHolder,
		UFloor:            shortPlaceHolder,
	}
	if application != nil {
		data.TenantName = application.FullName
		data.TenantEmail = application.Email
		data.TenantPhone = application.Phone
		data.MoveinDate = utils.Ternary(application.MoveinDate.IsZero(), html_util.NewHTMLTime(time.Now()), html_util.NewHTMLTime(application.MoveinDate))
		data.StartDate = data.Date
		data.EndDate = utils.Ternary(
			application.MoveinDate.IsZero(),
			html_util.NewHTMLTime(time.Now().AddDate(0, int(application.PreferredTerm), 0)),
			html_util.NewHTMLTime(application.MoveinDate.AddDate(0, int(application.PreferredTerm), 0)),
		)
		data.RentalDuration = fmt.Sprintf("%d", application.PreferredTerm)
		data.RentalDurationStr, _ = number.ToStr(int64(application.PreferredTerm))
	}
	if property.NumberOfFloors != nil {
		data.NumberOfFloors = fmt.Sprintf("%d", *property.NumberOfFloors)
	}
	if property.Building != nil {
		data.PBuilding = *property.Building
	}
	if property.Project != nil {
		data.PProject = *property.Project
	}

	buf, err := html_util.RenderHtml(data, templateFile)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
