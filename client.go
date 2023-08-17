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
	"time"

	"github.com/carlmjohnson/requests"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

const (
	baseURL = "https://api.lexoffice.io"
)

// Client is to define the request data
type Client struct {
	baseUrl string

	// httpClient is the client used to make HTTP requests.
	httpClient *http.Client

	rate *rate.Limiter
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

func WithRate(opPerSecond int) func(*Client) {
	return func(c *Client) {
		c.rate = rate.NewLimiter(rate.Every(time.Second/time.Duration(opPerSecond)), 1)
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

	// if rate is set, wrap transport with a rate limiter
	if client.rate != nil {
		client.httpClient.Transport = rateTransport{
			limiter: client.rate,
			base:    client.httpClient.Transport,
		}
	}

	return client
}

func (c *Client) Request(path string) *requests.Builder {
	return requests.
		URL(c.baseUrl).
		Path(path).
		Accept("application/json").
		Client(c.httpClient)
}

func (c *Client) Requestf(path string, args ...any) *requests.Builder {
	return requests.
		URL(c.baseUrl).
		Pathf(path, args...).
		Accept("application/json").
		Client(c.httpClient)
}

type rateTransport struct {
	limiter *rate.Limiter
	base    http.RoundTripper
}

func (t rateTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}
