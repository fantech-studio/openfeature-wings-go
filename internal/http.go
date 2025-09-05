package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v5"
	of "github.com/open-feature/go-sdk/openfeature"
)

type EvalRequest struct {
	ID   string         `json:"id" validate:"required"`
	Meta map[string]any `json:"meta"`
}

type (
	EvalResponse struct {
		Variant string       `json:"variant"`
		Bool    *BoolValue   `json:"bool,omitempty"`
		Int     *IntValue    `json:"int,omitempty"`
		Float   *FloatValue  `json:"float,omitempty"`
		String  *StringValue `json:"string,omitempty"`
		Object  *ObjectValue `json:"object,omitempty"`
	}

	BoolValue struct {
		Value bool `json:"value"`
	}

	IntValue struct {
		Value int64 `json:"value"`
	}

	FloatValue struct {
		Value float64 `json:"value"`
	}

	StringValue struct {
		Value string `json:"value"`
	}

	ObjectValue struct {
		Value map[string]any `json:"value"`
	}
)

type Config struct {
	Host          string
	MaxRetries    uint
	RetryInterval time.Duration
}

type Client interface {
	Do(ctx context.Context, path, method string, req *EvalRequest) (*EvalResponse, error)
}

type client struct {
	cli    *http.Client
	config *Config
}

func NewClient(config *Config) Client {
	return &client{
		cli:    new(http.Client),
		config: config,
	}
}

func (c *client) Do(
	ctx context.Context, path, method string, req *EvalRequest,
) (*EvalResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(c.config.Host)
	if err != nil {
		return nil, err
	}
	u.Scheme = "https"
	u = u.JoinPath(path)

	var buf *bytes.Buffer
	ope := func() (*EvalResponse, error) {
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
			var res *EvalResponse
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
		backoff.WithMaxTries(c.config.MaxRetries),
		backoff.WithBackOff(backoff.NewConstantBackOff(c.config.RetryInterval)))

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
