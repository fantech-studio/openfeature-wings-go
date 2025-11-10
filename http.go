package openfeaturewings

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v5"
	of "github.com/open-feature/go-sdk/openfeature"

	"github.com/fantech-studio/openfeature-wings-go/internal"
)

type config struct {
	host          string
	cli           *http.Client
	maxRetries    uint
	retryInterval time.Duration
	apiKey        string
	apiKeyID      string
}

type client interface {
	do(ctx context.Context, path, method string, req *internal.EvalRequest) (*internal.EvalResponse, error)
}

type httpClient struct {
	cli           *http.Client
	host          string
	maxRetries    uint
	retryInterval time.Duration
	apiKey        string
	apiKeyID      string
}

func newClient(config *config) client {
	return &httpClient{
		host:          config.host,
		cli:           config.cli,
		maxRetries:    config.maxRetries,
		retryInterval: config.retryInterval,
		apiKey:        config.apiKey,
		apiKeyID:      config.apiKeyID,
	}
}

func (c *httpClient) do(
	ctx context.Context, path, method string, req *internal.EvalRequest,
) (*internal.EvalResponse, error) {
	// #nosec G101
	const (
		apiKeyHeader   = "x-api-key"
		apiKeyIDHeader = "x-api-key-id"
	)
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	u.Scheme = "https"
	u = u.JoinPath(path)

	buf := &bytes.Buffer{}
	ope := func() (*internal.EvalResponse, error) {
		hReq, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(reqBytes))
		if err != nil {
			return nil, err
		}

		hReq.Header.Add(apiKeyHeader, c.apiKey)
		hReq.Header.Add(apiKeyIDHeader, c.apiKeyID)

		resp, err := c.cli.Do(hReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var res *internal.EvalResponse
			err = json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				return nil, err
			}
			return res, nil
		}

		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return nil, err
		}

		resolutionErr := resolveStatusCode(resp.StatusCode)(buf.String())
		buf.Reset()

		if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
			return nil, backoff.Permanent(resolutionErr)
		}

		return nil, resolutionErr
	}

	res, err := backoff.Retry(ctx, ope,
		backoff.WithMaxTries(c.maxRetries),
		backoff.WithBackOff(backoff.NewConstantBackOff(c.retryInterval)))

	return res, err
}

func resolveStatusCode(statusCode int) func(string) of.ResolutionError {
	switch statusCode {
	case http.StatusBadRequest:
		return of.NewInvalidContextResolutionError
	case http.StatusNotFound:
		return of.NewFlagNotFoundResolutionError
	default:
		return func(msg string) of.ResolutionError {
			return of.NewGeneralResolutionError(msg)
		}
	}
}
