//**********************************************************
//
// This file is part of lexoffice.
// All code may be used. Feel free and maybe code something better.
//
// Author: Jonas Kwiedor
//
//**********************************************************

package golexoffice

import (
	"net/http"

	"github.com/carlmjohnson/requests"
	"golang.org/x/oauth2"
)

const (
	baseURL = "https://api.lexoffice.io"
)

// Client is to define the request data
type Client struct {
	baseUrl string

	// httpClient is the client used to make HTTP requests.
	httpClient *http.Client
}

func WithClient(client *http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = client
	}
}

func WithBaseUrl(baseUrl string) func(*Client) {
	return func(c *Client) {
		c.baseUrl = baseUrl
	}
}

func NewClient(token string, o ...func(*Client)) *Client {
	client := &Client{
		httpClient: http.DefaultClient,
		baseUrl:    baseURL,
	}

	for _, option := range o {
		option(client)
	}

	// set client for auth
	client.httpClient.Transport = &oauth2.Transport{
		Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		Base:   client.httpClient.Transport,
	}

	return client
}

func (c *Client) Request(path string) *requests.Builder {
	return requests.
		URL(c.baseUrl).
		Path(path).
		Client(c.httpClient)
}
