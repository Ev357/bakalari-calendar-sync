package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	// "net/http"
)

func main() {
	godotenv.Overload()

	config, err := getConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(config)
}

type Config struct {
	url      string
	username string
	password string
}

func getConfig() (*Config, error) {
	url := os.Getenv("URL")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	if url == "" {
		return nil, errors.New("Url is not set")
	}

	if username == "" {
		return nil, errors.New("Username is not set")
	}

	if password == "" {
		return nil, errors.New("Password is not set")
	}

	return &Config{
		url,
		username,
		password,
	}, nil
}
