package text

import (
	"bytes"
	"path"
	"path/filepath"
	"text/template"
)

func RenderText(data any, textTemplateFile string) ([]byte, error) {
	textTemplateFile, err := filepath.Abs(textTemplateFile)
	if err != nil {
		return nil, err
	}

	fileName := path.Base(textTemplateFile)
	tmpl, err := template.New(fileName).ParseFiles(textTemplateFile)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)

	err = tmpl.Execute(buffer, data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
