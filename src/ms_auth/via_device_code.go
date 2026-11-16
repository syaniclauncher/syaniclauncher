package ms_auth

// to get AuthResult via Device Code way

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type deviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int64  `json:"interval"`
	ExpiresIn       int64  `json:"expires_in"`
	Error           string `json:"error"`
	Description     string `json:"error_description"`
}

func StartMicrosoftDeviceCodeAuth() (result <-chan AuthResult, err error, verificationURI, userCode string) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("scope", microsoftAuthScope)

	req, err := http.NewRequest(http.MethodPost, microsoftDeviceCodeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err, "", ""
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err, "", ""
	}
	defer resp.Body.Close()

	var deviceCode deviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceCode); err != nil {
		return nil, err, "", ""
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Microsoft device code: %s: %s",
			deviceCode.Error, deviceCode.Description,
		), "", ""
	}
	if deviceCode.DeviceCode == "" || deviceCode.UserCode == "" || deviceCode.VerificationURI == "" {
		return nil, fmt.Errorf("Microsoft device code response is missing required fields"), "", ""
	}
	if deviceCode.Interval == 0 {
		deviceCode.Interval = 5
	}
	if deviceCode.Interval < 0 || deviceCode.ExpiresIn <= 0 {
		return nil, fmt.Errorf("Microsoft device code response has invalid timing values"), "", ""
	}
	maxSeconds := int64((1<<63 - 1) / int64(time.Second))
	if deviceCode.Interval > maxSeconds || deviceCode.ExpiresIn > maxSeconds {
		return nil, fmt.Errorf("Microsoft device code response has invalid timing values"), "", ""
	}

	resultChan := make(chan AuthResult, 1)
	expiresAt := time.Now().Add(time.Duration(deviceCode.ExpiresIn) * time.Second)
	interval := time.Duration(deviceCode.Interval) * time.Second
	go func() {
		authResult, err := pollMicrosoftDeviceCode(client, microsoftTokenURL, deviceCode.DeviceCode, interval, expiresAt)
		if err != nil {
			resultChan <- AuthResult{Err: err}
			return
		}
		resultChan <- authResult
	}()

	return resultChan, nil, deviceCode.VerificationURI, deviceCode.UserCode
}

func pollMicrosoftDeviceCode(client *http.Client, endpoint, deviceCode string, interval time.Duration, expiresAt time.Time) (AuthResult, error) {
	for {
		remaining := time.Until(expiresAt)
		if remaining <= 0 {
			return AuthResult{}, fmt.Errorf("Microsoft device code expired")
		}
		wait := interval
		if wait > remaining {
			wait = remaining
		}
		time.Sleep(wait)
		if time.Until(expiresAt) <= 0 {
			return AuthResult{}, fmt.Errorf("Microsoft device code expired")
		}

		form := url.Values{}
		form.Set("client_id", clientID)
		form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		form.Set("device_code", deviceCode)
		req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
		if err != nil {
			return AuthResult{}, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return AuthResult{}, err
		}
		var token tokenResponse
		decodeErr := json.NewDecoder(resp.Body).Decode(&token)
		resp.Body.Close()
		if decodeErr != nil {
			return AuthResult{}, decodeErr
		}
		if resp.StatusCode == http.StatusOK {
			return AuthResult{
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				ExpiresIn:    token.ExpiresIn,
			}, nil
		}

		switch token.Error {
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5 * time.Second
		default:
			return AuthResult{}, fmt.Errorf(
				"Microsoft device code: %s: %s",
				token.Error, token.Description,
			)
		}
	}
}
