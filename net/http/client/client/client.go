package client

import (
	"net/http"
	"time"
)

var DefaultClient *http.Client

func init() {
	DefaultClient = http.DefaultClient
}

type ClientOption func(c *http.Client)

func NewClient(options ...ClientOption) *http.Client {
	c := &http.Client{}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithTimeout(t time.Duration) ClientOption {
	return func(c *http.Client) {
		c.Timeout = t
	}
}
