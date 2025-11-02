package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Permission struct {
	Roles   []string `json:"roles"`
	Grantee struct {
		Email string `json:"emailAddress"`
	} `json:"grantee"`
}

type PermissionsResponse struct {
	Value []Permission `json:"value"`
}

func trimTrailingComma(s string) string {
	return strings.TrimSuffix(s, ", ")
}

func contains(slice []string, v string) bool {
	for _, s := range slice {
		if s == v {
			return true
		}
	}
	return false
}

func getDelegation(token, userID string) (fullAccess, sendAs, sendOnBehalf string, err error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/mailFolders/Inbox/permissions", userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var permResp PermissionsResponse
	if err = json.Unmarshal(body, &permResp); err != nil {
		return
	}

	var fb strings.Builder
	var sb strings.Builder
	var hb strings.Builder

	for _, p := range permResp.Value {
		email := p.Grantee.Email
		if email == "" {
			continue
		}
		if contains(p.Roles, "Owner") {
			fb.WriteString(email)
			fb.WriteString(", ")
		}
		if contains(p.Roles, "SendAs") {
			sb.WriteString(email)
			sb.WriteString(", ")
		}
		if contains(p.Roles, "SendOnBehalf") {
			hb.WriteString(email)
			hb.WriteString(", ")
		}
	}

	fullAccess = trimTrailingComma(fb.String())
	sendAs = trimTrailingComma(sb.String())
	sendOnBehalf = trimTrailingComma(hb.String())
	return
}
