package ms_auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	clientID    = "e4789891-9cae-424f-baba-4f8aa9cf90c1"
	redirectURI = "https://syaniclauncher.github.io/auth_callback"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	Err          error
}

func StartMicrosoftBrowserAuth() (string, <-chan AuthResult, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0") // bind on free port
	if err != nil {
		return "", nil, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	state := "port" + strconv.Itoa(port) //port<port>

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", redirectURI)
	params.Set("response_mode", "query")
	params.Set("scope", "XboxLive.signin offline_access")
	params.Set("state", state)

	authURL := "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?" + params.Encode()

	result := make(chan AuthResult, 1)
	mux := http.NewServeMux()
	server := &http.Server{Handler: mux}

	mux.HandleFunc("/auth_callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			return
		}

		token, err := exchangeCode(code)
		if err != nil {
			http.Error(w, "Token exchange failed", http.StatusBadGateway)
			result <- AuthResult{Err: err}
		} else {
			fmt.Fprintln(w, "Authentication successful! You can close this window.")
			result <- token
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = server.Shutdown(ctx)
		}()
	})

	go func() {
		if err := server.Serve(listener); err != nil &&
			err != http.ErrServerClosed {
			log.Println("Auth server:", err)
		}
	}()

	// close temporary server if user never completes authentication
	go func() {
		time.Sleep(5 * time.Minute)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	return authURL, result, nil
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Error        string `json:"error"`
	Description  string `json:"error_description"`
}

func exchangeCode(code string) (AuthResult, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://login.microsoftonline.com/consumers/oauth2/v2.0/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return AuthResult{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return AuthResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AuthResult{}, err
	}

	var token tokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return AuthResult{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return AuthResult{}, fmt.Errorf(
			"Microsoft OAuth: %s: %s",
			token.Error, token.Description,
		)
	}

	return AuthResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}
