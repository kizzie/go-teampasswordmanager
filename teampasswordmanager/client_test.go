package teampasswordmanager

import "testing"

func TestNewClient(t *testing.T) {
	clientConfig := ClientConfig{
		BaseURL:   "localhost",
		AuthToken: "1234",
	}

	expectedAPIURL := "localhost/index.php/api/v4/"

	client, _ := NewClient(&clientConfig)

	if client.apiURL != expectedAPIURL {
		t.Errorf("API Url is not formed correctly, expected (%s) got (%s)", expectedAPIURL, client.apiURL)
	}
}
