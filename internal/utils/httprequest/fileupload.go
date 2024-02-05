package httprequest

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// FileUploadMultipart uploads a file to a server using multipart form
// filePath: the absolute path to the file to be uploaded
// url: the server url to upload the file to
func FileUploadMultipart(filePath string, url string, method string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	return client.Do(r)
}

// FileUpload uploads a file to a server
// filePath: the absolute path to the file to be uploaded
// url: the server url to upload the file to
func FileUpload(filePath string, url string, method string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r, err := http.NewRequest(method, url, file)
	if err != nil {
		return nil, err
	}
	fstat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", http.DetectContentType(make([]byte, fstat.Size())))

	client := &http.Client{}
	return client.Do(r)
}
