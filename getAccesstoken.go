package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

// getAccessToken retrieves an Azure AD token using client credentials flow.
func getAccessToken() (string, error) {

	form := url.Values{}
	form.Set("client_id", cfg.AzureClientID)
	form.Set("scope", "https://graph.microsoft.com/.default")
	form.Set("client_secret", cfg.AzureClientSecret)
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", cfg.AzureTenantID),
		strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}
	return tr.AccessToken, nil
}
