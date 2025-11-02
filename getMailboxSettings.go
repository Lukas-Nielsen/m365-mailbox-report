package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MailboxSettings struct {
	ForwardingAddress string `json:"forwardingSmtpAddress"`
	IsForwarding      bool   `json:"isForwardingEnabled"`
	UserPurpose       string `json:"userPurpose"`
}

var httpClient = &http.Client{}

func getMailboxSettings(token, userID string) (MailboxSettings, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/mailboxSettings", userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return MailboxSettings{}, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return MailboxSettings{}, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return MailboxSettings{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MailboxSettings{}, fmt.Errorf("reading response: %w", err)
	}

	var ms MailboxSettings
	if err := json.Unmarshal(body, &ms); err != nil {
		return MailboxSettings{}, fmt.Errorf("parsing JSON: %w", err)
	}
	return ms, nil
}
