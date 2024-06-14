package text

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	// iterate over template files in template subdirectory
	files, err := os.ReadDir("templates")
	require.NoError(t, err)

	data := struct {
		Name string
	}{
		Name: "Tehc's School",
	}
	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			buf, err := RenderText(data, "templates/"+file.Name())
			require.NoError(t, err)
			t.Log(string(buf))
		})
	}
}
