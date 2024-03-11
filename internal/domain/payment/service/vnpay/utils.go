package vnpay

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
)

func sortObject(obj map[string]string) map[string]string {
	sorted := make(map[string]string)
	var keys []string

	for key := range obj {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		sorted[key] = obj[key]
	}

	return sorted
}

func stringify(vnpParams map[string]string) string {
	values := url.Values{}
	for key, value := range vnpParams {
		values.Add(key, value)
	}
	return values.Encode()
}

func sendHttpRequest(endpoint, method string, body any) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		method,
		endpoint,
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
