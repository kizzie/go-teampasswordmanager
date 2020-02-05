package teampasswordmanager

import (
	"net/http"
)

// ClientConfig stores the config for the team password manager http client
type ClientConfig struct {
	BaseURL   string
	AuthToken string
}

// Client is the http client, api and auth token for team password manager
type Client struct {
	httpClient *http.Client
	apiURL     string
	authToken  string
}

// NewClient returns a configure team password manager client
func NewClient(config *ClientConfig) (Client, error) {
	httpClient := &http.Client{}
	apiURL := config.BaseURL + "/index.php/api/v4/"

	return Client{
		httpClient: httpClient,
		apiURL:     apiURL,
		authToken:  config.AuthToken,
	}, nil
}
