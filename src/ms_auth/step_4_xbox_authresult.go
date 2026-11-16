package ms_auth

import (
	"fmt"
	"net/http"
)

type xboxAuthResponse struct {
	Token         string `json:"Token"`
	DisplayClaims struct {
		XUI []struct {
			UserHash string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

func AuthenticateXboxLive(microsoftAccessToken string) (XboxAuthResult, error) {
	return requestXboxLiveAuthentication(xboxAuthURL, microsoftAccessToken)
}

func requestXboxLiveAuthentication(endpoint, microsoftAccessToken string) (XboxAuthResult, error) {
	if microsoftAccessToken == "" {
		return XboxAuthResult{}, fmt.Errorf("Microsoft access token is empty")
	}

	payload := map[string]any{
		"Properties": map[string]string{
			"AuthMethod": "RPS",
			"SiteName":   xboxSiteName,
			"RpsTicket":  "d=" + microsoftAccessToken,
		},
		"RelyingParty": xboxRelyingParty,
		"TokenType":    "JWT",
	}
	body, status, err := requestAuthAPI(http.MethodPost, endpoint, payload, "")
	if err != nil {
		return XboxAuthResult{}, err
	}
	var response xboxAuthResponse
	if err := decodeAuthResponse("Xbox Live authentication", status, body, &response); err != nil {
		return XboxAuthResult{}, err
	}
	if response.Token == "" || len(response.DisplayClaims.XUI) == 0 || response.DisplayClaims.XUI[0].UserHash == "" {
		return XboxAuthResult{}, fmt.Errorf("Xbox Live response is missing token or user hash")
	}
	return XboxAuthResult{
		Token:    response.Token,
		UserHash: response.DisplayClaims.XUI[0].UserHash,
	}, nil
}
