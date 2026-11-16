package ms_auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RefreshMicrosoftAuthResult(refreshToken string) (AuthResult, error) {
	if refreshToken == "" {
		return AuthResult{}, fmt.Errorf("Microsoft refresh token is empty")
	}

	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	form.Set("scope", microsoftAuthScope)

	req, err := http.NewRequest(http.MethodPost, microsoftTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return AuthResult{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return AuthResult{}, err
	}
	defer resp.Body.Close()

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return AuthResult{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return AuthResult{}, fmt.Errorf("Microsoft token refresh: %s: %s", token.Error, token.Description)
	}
	if token.RefreshToken == "" {
		token.RefreshToken = refreshToken
	}

	return AuthResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}
