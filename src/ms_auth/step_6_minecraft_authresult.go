package ms_auth

import (
	"context"
	"fmt"
	"net/http"
)

type minecraftAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func AuthenticateMinecraft(userHash, xstsToken string) (MinecraftAuthResult, error) {
	return AuthenticateMinecraftContext(context.Background(), userHash, xstsToken)
}

func AuthenticateMinecraftContext(ctx context.Context, userHash, xstsToken string) (MinecraftAuthResult, error) {
	return requestMinecraftAuthentication(ctx, minecraftLoginURL, userHash, xstsToken)
}

func requestMinecraftAuthentication(ctx context.Context, endpoint, userHash, xstsToken string) (MinecraftAuthResult, error) {
	if userHash == "" || xstsToken == "" {
		return MinecraftAuthResult{}, fmt.Errorf("Xbox user hash and XSTS token are required")
	}

	payload := map[string]string{
		"identityToken": fmt.Sprintf("XBL3.0 x=%s;%s", userHash, xstsToken),
	}
	body, status, err := requestAuthAPIContext(ctx, http.MethodPost, endpoint, payload, "")
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
