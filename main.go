package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Overload()

	config, err := getConfig()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loginUrl := config.url + "/bakaweb/Login"

	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.PostForm(loginUrl, url.Values{
		"username":   {config.username},
		"password":   {config.password},
		"persistent": {"true"},
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	content, _ := io.ReadAll(resp.Body)

	fmt.Println(string(content))
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
