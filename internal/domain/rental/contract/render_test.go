package contract

import (
	"testing"

	"github.com/stretchr/testify/require"
	application_model "github.com/user2410/rrms-backend/internal/domain/application/model"
	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	"github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
)

func TestRenderHtml(t *testing.T) {
	var (
		pr    model.PrerentalModel
		a     application_model.ApplicationModel
		p     property_model.PropertyModel
		unit  unit_model.UnitModel
		owner auth_model.UserModel
	)
	p.Type = "APARTMENT"
	res, err := RenderContractTemplate(&pr, &a, &p, &unit, &owner)
	require.NoError(t, err)
	t.Log(res)

	p.Type = "ROOM"
	res, err = RenderContractTemplate(&pr, &a, &p, &unit, &owner)
	require.NoError(t, err)
	t.Log(res)

	p.Type = "PRIVATE"
	res, err = RenderContractTemplate(&pr, &a, &p, &unit, &owner)
	require.NoError(t, err)
	t.Log(res)

	p.Type = "OFFICE"
	res, err = RenderContractTemplate(&pr, &a, &p, &unit, &owner)
	require.NoError(t, err)
	t.Log(res)

}
