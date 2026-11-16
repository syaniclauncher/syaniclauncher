package ms_auth

import (
	"fmt"
	"net/http"
)

type minecraftAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func AuthenticateMinecraft(userHash, xstsToken string) (MinecraftAuthResult, error) {
	return requestMinecraftAuthentication(minecraftLoginURL, userHash, xstsToken)
}

func requestMinecraftAuthentication(endpoint, userHash, xstsToken string) (MinecraftAuthResult, error) {
	if userHash == "" || xstsToken == "" {
		return MinecraftAuthResult{}, fmt.Errorf("Xbox user hash and XSTS token are required")
	}

	payload := map[string]string{
		"identityToken": fmt.Sprintf("XBL3.0 x=%s;%s", userHash, xstsToken),
	}
	body, status, err := requestAuthAPI(http.MethodPost, endpoint, payload, "")
	if err != nil {
		return MinecraftAuthResult{}, err
	}
	var response minecraftAuthResponse
	if err := decodeAuthResponse("Minecraft authentication", status, body, &response); err != nil {
		return MinecraftAuthResult{}, err
	}
	if response.AccessToken == "" {
		return MinecraftAuthResult{}, fmt.Errorf("Minecraft response is missing access token")
	}
	return MinecraftAuthResult{
		AccessToken: response.AccessToken,
		ExpiresIn:   response.ExpiresIn,
	}, nil
}
