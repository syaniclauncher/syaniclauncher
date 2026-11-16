package ms_auth

// to get AuthResult via Oauth way

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

func StartMicrosoftBrowserAuth() (string, <-chan AuthResult, error) {
	return StartMicrosoftBrowserAuthContext(context.Background())
}

func StartMicrosoftBrowserAuthContext(ctx context.Context) (string, <-chan AuthResult, error) {
	if err := ctx.Err(); err != nil {
		return "", nil, err
	}
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
	params.Set("scope", microsoftAuthScope)
	params.Set("state", state)

	authURL := microsoftAuthorizeURL + "?" + params.Encode()

	result := make(chan AuthResult, 1)
	mux := http.NewServeMux()
	server := &http.Server{Handler: mux}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `<body style="margin:0;height:100vh;display:grid;place-content:center;text-align:center;background:#000;color:#fff;font:14px/1.5 system-ui,sans-serif"><img src="https://syaniclauncher.github.io/misc/ourport.png" style="max-width:90vw;max-height:80vh"><small style="color:#aaa;letter-spacing:.3px">This port is temporarily used for ms oauth</small>`)
	})
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

		token, err := exchangeOauthCode(ctx, code)
		if err != nil {
			http.Error(w, "Token exchange failed", http.StatusBadGateway)
			select {
			case result <- AuthResult{Err: err}:
			default:
			}
		} else {
			http.Redirect(w, r, successURI, http.StatusFound) //redirect to hide code in browser
			select {
			case result <- token:
			default:
			}
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

	go func() {
		timer := time.NewTimer(5 * time.Minute)
		defer timer.Stop()
		select {
		case <-ctx.Done():
		case <-timer.C:
			select {
			case result <- AuthResult{Err: fmt.Errorf("Microsoft browser sign-in timed out")}:
			default:
			}
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
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

func exchangeOauthCode(ctx context.Context, code string) (AuthResult, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		microsoftTokenURL,
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
