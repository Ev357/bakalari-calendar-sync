package main

import (
	"errors"
	"os"
)

type Config struct {
	url          string
	username     string
	password     string
	clientId     string
	clientSecret string
	refreshToken string
}

func getConfig() (*Config, error) {
	url := os.Getenv("URL")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	refreshToken := os.Getenv("REFRESH_TOKEN")

	if url == "" {
		return nil, errors.New("Url is not set")
	}

	if username == "" {
		return nil, errors.New("Username is not set")
	}

	if password == "" {
		return nil, errors.New("Password is not set")
	}

	if clientId == "" {
		return nil, errors.New("Client ID is not set")
	}

	if clientSecret == "" {
		return nil, errors.New("Client secret is not set")
	}

	if refreshToken == "" {
		return nil, errors.New("Refresh token is not set")
	}

	return &Config{
		url,
		username,
		password,
		clientId,
		clientSecret,
		refreshToken,
	}, nil
}
