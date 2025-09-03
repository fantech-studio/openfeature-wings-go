package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

type Client struct {
	Cli  *http.Client
	Host string
}

func (c *Client) Do(
	ctx context.Context, path, method string, req *EvalRequest,
) (*EvalResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", c.Host, path)
	hReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}

	resp, err := c.Cli.Do(hReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		res := new(EvalResponse)
		err = json.NewDecoder(resp.Body).Decode(res)
		return res, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	resolutionErr := resolveStatusCode(resp.StatusCode)
	return nil, resolutionErr(buf.String())
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
