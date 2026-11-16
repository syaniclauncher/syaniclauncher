package ms_auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func requestAuthAPI(method, endpoint string, payload any, bearerToken string) ([]byte, int, error) {
	var requestBody io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, err
		}
		requestBody = bytes.NewReader(encoded)
	}
	req, err := http.NewRequest(method, endpoint, requestBody)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

func decodeAuthResponse(service string, status int, body []byte, result any) error {
	if status != http.StatusOK {
		message := strings.TrimSpace(string(body))
		if message == "" {
			return fmt.Errorf("%s failed with HTTP %d", service, status)
		}
		return fmt.Errorf("%s failed with HTTP %d: %s", service, status, message)
	}
	return json.Unmarshal(body, result)
}
