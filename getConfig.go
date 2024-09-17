package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	url                 string
	username            string
	password            string
	account             string
	serviceAccount      []byte
	serviceAccountEmail string
}

func getConfig() (*Config, error) {
	url := os.Getenv("URL")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	account := os.Getenv("ACCOUNT")
	serviceAccountBase64 := os.Getenv("SERVICE_ACCOUNT")

	if url == "" {
		return nil, errors.New("Url is not set")
	}

	if username == "" {
		return nil, errors.New("Username is not set")
	}

	if password == "" {
		return nil, errors.New("Password is not set")
	}

	if account == "" {
		return nil, errors.New("Account is not set")
	}

	if serviceAccountBase64 == "" {
		return nil, errors.New("Service account is not set")
	}

	serviceAccount, err := base64.StdEncoding.DecodeString(serviceAccountBase64)

	if err != nil {
		return nil, err
	}

	type ServiceAccount struct {
		ClientEmail string `json:"client_email"`
	}

	parsedServiceAccount := ServiceAccount{}
	err = json.Unmarshal(serviceAccount, &parsedServiceAccount)

	if err != nil {
		return nil, err
	}

	serviceAccountEmail := parsedServiceAccount.ClientEmail

	return &Config{
		url,
		username,
		password,
		account,
		serviceAccount,
		serviceAccountEmail,
	}, nil
}
