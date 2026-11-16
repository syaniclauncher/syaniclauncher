package ms_auth

import (
	"context"
	"fmt"
	"net/http"
)

type minecraftProfileResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetMinecraftProfile(accessToken string) (MinecraftProfile, error) {
	return GetMinecraftProfileContext(context.Background(), accessToken)
}

func GetMinecraftProfileContext(ctx context.Context, accessToken string) (MinecraftProfile, error) {
	return requestMinecraftProfile(ctx, minecraftProfileURL, accessToken)
}

func requestMinecraftProfile(ctx context.Context, endpoint, accessToken string) (MinecraftProfile, error) {
	if accessToken == "" {
		return MinecraftProfile{}, fmt.Errorf("Minecraft access token is empty")
	}

	body, status, err := requestAuthAPIContext(ctx, http.MethodGet, endpoint, nil, accessToken)
	if err != nil {
		return MinecraftProfile{}, err
	}
	if status == http.StatusNotFound {
		return MinecraftProfile{}, fmt.Errorf("Minecraft Java Edition profile was not found")
	}
	var response minecraftProfileResponse
	if err := decodeAuthResponse("Minecraft profile lookup", status, body, &response); err != nil {
		return MinecraftProfile{}, err
	}
	if response.ID == "" || response.Name == "" {
		return MinecraftProfile{}, fmt.Errorf("Minecraft profile response is missing ID or name")
	}
	return MinecraftProfile{ID: response.ID, Name: response.Name}, nil
}
