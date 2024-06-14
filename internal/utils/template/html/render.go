package html

import (
	"bytes"
	"html/template"
	"path"
	"path/filepath"
	"time"
)

type HTMLTime struct {
	Date   int
	Month  int
	Year   int
	Hour   int
	Minute int
}

func NewHTMLTime(t time.Time) HTMLTime {
	return HTMLTime{
		Date:   t.Day(),
		Month:  int(t.Month()),
		Year:   t.Year(),
		Hour:   t.Hour(),
		Minute: t.Minute(),
	}
}

// Render HTML byte slice from template file and data.
// data is a struct holding data for replacing placeholders
// in the target template.
// htmlTemplateFile is the relative path to the template file.
func RenderHtml(data any, htmlTemplateFile string) ([]byte, error) {
	htmlTemplateFile, err := filepath.Abs(htmlTemplateFile)
	if err != nil {
		return nil, err
	}

	fileName := path.Base(htmlTemplateFile)
	tmpl, err := template.New(fileName).ParseFiles(htmlTemplateFile)
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
