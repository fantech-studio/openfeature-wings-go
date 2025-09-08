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
	maxRetries    uint
	retryInterval time.Duration
}

type client interface {
	do(ctx context.Context, path, method string, req *internal.EvalRequest) (*internal.EvalResponse, error)
}

type httpClient struct {
	cli    *http.Client
	config *config
}

func newClient(config *config) client {
	return &httpClient{
		cli:    new(http.Client),
		config: config,
	}
}

func (c *httpClient) do(
	ctx context.Context, path, method string, req *internal.EvalRequest,
) (*internal.EvalResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(c.config.host)
	if err != nil {
		return nil, err
	}
	u.Scheme = "https"
	u = u.JoinPath(path)

	var buf *bytes.Buffer
	ope := func() (*internal.EvalResponse, error) {
		hReq, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(reqBytes))
		if err != nil {
			return nil, err
		}

		resp, err := c.cli.Do(hReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var res *internal.EvalResponse
			err = json.NewDecoder(resp.Body).Decode(res)
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

		if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
			return nil, backoff.Permanent(resolutionErr)
		}

		return nil, resolutionErr
	}

	res, err := backoff.Retry(ctx, ope,
		backoff.WithMaxTries(c.config.maxRetries),
		backoff.WithBackOff(backoff.NewConstantBackOff(c.config.retryInterval)))

	return res, err
}

func resolveStatusCode(statusCode int) func(string) of.ResolutionError {
	switch statusCode {
	case http.StatusBadRequest:
		return of.NewInvalidContextResolutionError
	case http.StatusNotFound:
		return of.NewFlagNotFoundResolutionError
	default:
		return of.NewGeneralResolutionError
	}
}
