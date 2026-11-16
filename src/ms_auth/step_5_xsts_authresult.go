package ms_auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type xstsErrorResponse struct {
	XErr    int64  `json:"XErr"`
	Message string `json:"Message"`
}

func AuthorizeXSTS(xboxToken string) (XSTSAuthResult, error) {
	return requestXSTSAuthorization(xstsAuthorizeURL, xboxToken)
}

func requestXSTSAuthorization(endpoint, xboxToken string) (XSTSAuthResult, error) {
	if xboxToken == "" {
		return XSTSAuthResult{}, fmt.Errorf("Xbox Live token is empty")
	}

	payload := map[string]any{
		"Properties": map[string]any{
			"SandboxId": "RETAIL",
			"UserTokens": []string{
				xboxToken,
			},
		},
		"RelyingParty": minecraftRelyingParty,
		"TokenType":    "JWT",
	}
	body, status, err := requestAuthAPI(http.MethodPost, endpoint, payload, "")
	if err != nil {
		return XSTSAuthResult{}, err
	}
	if status != http.StatusOK {
		var errorResponse xstsErrorResponse
		if json.Unmarshal(body, &errorResponse) == nil && errorResponse.XErr != 0 {
			message := errorResponse.Message
			if message == "" {
				message = xstsErrorMessage(errorResponse.XErr)
			}
			return XSTSAuthResult{}, fmt.Errorf("XSTS authorization failed (HTTP %d): %s", status, message)
		}
		message := strings.TrimSpace(string(body))
		if message == "" {
			return XSTSAuthResult{}, fmt.Errorf("XSTS authorization failed with HTTP %d", status)
		}
		return XSTSAuthResult{}, fmt.Errorf("XSTS authorization failed with HTTP %d: %s", status, message)
	}
	var response xboxAuthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return XSTSAuthResult{}, err
	}
	if response.Token == "" || len(response.DisplayClaims.XUI) == 0 || response.DisplayClaims.XUI[0].UserHash == "" {
		return XSTSAuthResult{}, fmt.Errorf("XSTS response is missing token or user hash")
	}
	return XSTSAuthResult{
		Token:    response.Token,
		UserHash: response.DisplayClaims.XUI[0].UserHash,
	}, nil
}

func xstsErrorMessage(code int64) string {
	if message := map[int64]string{
		2148916233: "No Xbox profile found. Complete Xbox account setup at xbox.com.",
		2148916235: "Xbox Live is unavailable for this account or region.",
		2148916236: "Adult verification is required for this Microsoft account.",
		2148916237: "Age verification is required for this Microsoft account.",
		2148916238: "This child account must be added to a Microsoft Family by an adult.",
	}[code]; message != "" {
		return message
	}
	return fmt.Sprintf("Unknown XSTS error: %d", code)
}
