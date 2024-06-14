package html

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	// iterate over template files in template subdirectory
	// files, err := os.ReadDir("templates")
	// require.NoError(t, err)

	data := struct {
		Date          HTMLTime
		Name          string
		ApplicationId string
		ListingTitle  string
	}{
		Date:          NewHTMLTime(time.Now()),
		Name:          "Tehc's School",
		ApplicationId: "123456",
		ListingTitle:  "A test listing",
	}
	buf, err := RenderHtml(data, "templates/test1.html")
	require.NoError(t, err)
	t.Log("test1.html:", string(buf))
}
