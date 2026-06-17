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

func WithCredentials(creds *Credentials) Option {
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

// WithInsecure disables transport security so that evaluation requests are sent
// over plain HTTP instead of the default HTTPS. Any scheme present in the host
// passed to [NewProvider] is ignored; this option alone determines whether
// requests use https (the default) or http.
//
// Use it only for trusted environments, such as a wings endpoint exposed over
// plain HTTP inside a cluster (for example, a development cluster). Without this
// option the provider always uses HTTPS.
//
// WARNING: with transport security disabled, the credentials configured via
// [WithCredentials] are transmitted over the network in plaintext. Never enable
// this for endpoints reachable over an untrusted network.
func WithInsecure() Option {
	return withInsecure{}
}

type withInsecure struct{}

func (withInsecure) apply(config *config) {
	config.insecure = true
}
