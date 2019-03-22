package client

import (
	"net/http"
	"time"
)

// NewClient returns a volusnap client
func NewClient(serverURL string) *Client {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	return &Client{
		serverURL:  serverURL,
		httpClient: httpClient,
	}
}

// Client struct
type Client struct {
	serverURL  string
	httpClient *http.Client
}
