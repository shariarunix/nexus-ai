package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthService struct {
	SupabaseURL string
	AnonKey     string
}

func NewAuthService(supabaseURL, anonKey string) *AuthService {
	return &AuthService{
		SupabaseURL: supabaseURL,
		AnonKey:     anonKey,
	}
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
	Error string `json:"error_description,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

func (s *AuthService) SignUp(email, password string) (*AuthResponse, error) {
	url := fmt.Sprintf("%s/auth/v1/signup", s.SupabaseURL)
	return s.makeRequest(url, email, password)
}

func (s *AuthService) SignIn(email, password string) (*AuthResponse, error) {
	url := fmt.Sprintf("%s/auth/v1/token?grant_type=password", s.SupabaseURL)
	return s.makeRequest(url, email, password)
}

func (s *AuthService) makeRequest(url, email, password string) (*AuthResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.AnonKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		// Try to parse error message from body if possible, or return generic
		var errResp struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
			Msg         string `json:"msg"`
		}
		_ = json.Unmarshal(respBody, &errResp)
		msg := errResp.Description
		if msg == "" {
			msg = errResp.Msg // specific to some endpoints
		}
		if msg == "" {
			msg = errResp.Error
		}
		if msg == "" {
			msg = string(respBody)
		}
		return nil, fmt.Errorf("auth error: %s", msg)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}
