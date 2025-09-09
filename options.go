package openfeaturewings

import (
	"net/http"
	"time"
)

type Option interface {
	apply(*config)
}

func WithMaxRetries(maxRetries uint) Option {
	return withMaxRetries(maxRetries)
}

type withMaxRetries uint

func (w withMaxRetries) apply(config *config) {
	config.maxRetries = uint(w)
}

func WithRetryInterval(retryInterval time.Duration) Option {
	return withRetryInterval(retryInterval)
}

type withRetryInterval time.Duration

func (w withRetryInterval) apply(config *config) {
	config.retryInterval = time.Duration(w)
}

func WithCredentianls(creds *Credentials) Option {
	return (*withCredentials)(creds)
}

type Credentials struct {
	APIKey   string
	APIKeyID string
}

type withCredentials Credentials

func (w *withCredentials) apply(config *config) {
	config.apiKey = w.APIKey
	config.apiKeyID = w.APIKeyID
}

func WithHTTPClient(cli *http.Client) Option {
	return &withHTTPClient{cli: cli}
}

type withHTTPClient struct{ cli *http.Client }

func (w *withHTTPClient) apply(config *config) {
	config.cli = w.cli
}
