package main

import (
	"fmt"
	"os"
)

type Config struct {
	AzureTenantID        string
	AzureClientID        string
	AzureClientSecret    string
	BookstackPageID      string
	BookstackTokenID     string
	BookstackTokenSecret string
	BookstackBaseURL     string
}

var cfg Config

func init() {
	cfg = Config{
		AzureTenantID:        getEnv("AZURE_TENANT_ID"),
		AzureClientID:        getEnv("AZURE_CLIENT_ID"),
		AzureClientSecret:    getEnv("AZURE_CLIENT_SECRET"),
		BookstackPageID:      getEnv("BOOKSTACK_PAGE_ID"),
		BookstackTokenID:     getEnv("BOOKSTACK_TOKEN_ID"),
		BookstackTokenSecret: getEnv("BOOKSTACK_TOKEN_SECRET"),
		BookstackBaseURL:     getEnv("BOOKSTACK_BASEURL"),
	}

	// Validate required fields
	if cfg.AzureTenantID == "" ||
		cfg.AzureClientID == "" ||
		cfg.AzureClientSecret == "" ||
		cfg.BookstackPageID == "" ||
		cfg.BookstackTokenID == "" ||
		cfg.BookstackTokenSecret == "" ||
		cfg.BookstackBaseURL == "" {
		panic("missing required environment variables")
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func main() {
	token, err := getAccessToken()
	if err != nil {
		panic(err)
	}

	users, err := getUsers(token)
	if err != nil {
		panic(err)
	}

	md := buildMarkdown(users, token)

	if err := UpdateBookstack(md); err != nil {
		panic(err)
	}

	fmt.Println("BookStack Seite erfolgreich aktualisiert")
}
