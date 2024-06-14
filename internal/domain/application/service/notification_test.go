package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	text_util "github.com/user2410/rrms-backend/internal/utils/template/text"
)

func TestXxx(t *testing.T) {
	pushContent, err := text_util.RenderText(
		struct {
			Status database.APPLICATIONSTATUS
		}{
			Status: database.APPLICATIONSTATUSAPPROVED,
		},
		"templates/title/update_application.txt",
	)
	require.NoError(t, err)
	t.Log(string(pushContent))
}
