package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultPageID      = "7"
	defaultTokenID     = "Zw1FcdN64pbGL3QJgNPXIlFhnP1Dwzrc"
	defaultTokenSecret = "LaEBPcxZJ7obZN1yUliHHPSqhcZ9PATW"
	defaultBaseURL     = "https://bs.h.lukasnielsen.de"
)

type bookstackClient struct {
	baseURL     string
	pageID      string
	tokenID     string
	tokenSecret string
	httpClient  *http.Client
}

func newBookstackClient() *bookstackClient {
	return &bookstackClient{
		baseURL:     cfg.BookstackBaseURL,
		pageID:      cfg.BookstackPageID,
		tokenID:     cfg.BookstackTokenID,
		tokenSecret: cfg.BookstackTokenSecret,
		httpClient:  &http.Client{},
	}
}

func (c *bookstackClient) updateMarkdown(md string) error {
	url := fmt.Sprintf("%s/api/pages/%s", c.baseURL, c.pageID)
	payload := map[string]string{
		"markdown": md,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s:%s", c.tokenID, c.tokenSecret))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("BookStack update failed: %s", string(respBody))
	}
	return nil
}

// UpdateBookstack is the exported helper used by other packages.
func UpdateBookstack(md string) error {
	client := newBookstackClient()
	return client.updateMarkdown(md)
}
