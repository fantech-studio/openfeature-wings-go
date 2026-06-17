package openfeaturewings

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fantech-studio/openfeature-wings-go/internal"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// TestHTTPClientDoURL verifies that the request scheme is forced to https by
// default and to http when WithInsecure is supplied, regardless of the scheme
// contained in the configured host.
func TestHTTPClientDoURL(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		opts    []Option
		wantURL string
	}{
		{
			name:    "default forces https",
			host:    "http://wings.example.com",
			opts:    nil,
			wantURL: "https://wings.example.com/bool:evaluate",
		},
		{
			name:    "WithInsecure forces http",
			host:    "https://wings.internal:8080",
			opts:    []Option{WithInsecure()},
			wantURL: "http://wings.internal:8080/bool:evaluate",
		},
		{
			name:    "WithInsecure keeps an http host on http",
			host:    "http://wings.internal:8080",
			opts:    []Option{WithInsecure()},
			wantURL: "http://wings.internal:8080/bool:evaluate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotURL string
			rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
				gotURL = r.URL.String()
				body, err := json.Marshal(&internal.EvalResponse{
					Variant: "on",
					Bool:    &internal.BoolValue{Value: true},
				})
				if err != nil {
					return nil, err
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			})

			opts := append([]Option{WithHTTPClient(&http.Client{Transport: rt})}, tt.opts...)
			config := resolveOptions(opts...)
			config.host = tt.host
			c := newClient(config)

			resp, err := c.do(
				context.Background(), "/bool:evaluate", http.MethodPost,
				&internal.EvalRequest{ID: "flag"},
			)
			if err != nil {
				t.Fatalf("do returned error: %v", err)
			}
			if gotURL != tt.wantURL {
				t.Errorf("request URL = %q, want %q", gotURL, tt.wantURL)
			}
			if resp == nil || resp.Bool == nil || !resp.Bool.Value {
				t.Errorf("unexpected response: %+v", resp)
			}
		})
	}
}

// TestHTTPClientDoWithInsecure verifies that with WithInsecure the request
// actually reaches a plain HTTP server and the response is decoded correctly.
func TestHTTPClientDoWithInsecure(t *testing.T) {
	const path = "/bool:evaluate"

	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if err := json.NewEncoder(w).Encode(&internal.EvalResponse{
			Variant: "on",
			Bool:    &internal.BoolValue{Value: true},
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer srv.Close()

	config := resolveOptions(WithInsecure())
	config.host = srv.URL // plain HTTP (http://127.0.0.1:port)
	c := newClient(config)

	resp, err := c.do(
		context.Background(), path, http.MethodPost,
		&internal.EvalRequest{ID: "flag"},
	)
	if err != nil {
		t.Fatalf("do returned error: %v", err)
	}
	if gotPath != path {
		t.Errorf("server received path = %q, want %q", gotPath, path)
	}
	if resp == nil || resp.Bool == nil || !resp.Bool.Value {
		t.Errorf("unexpected response: %+v", resp)
	}
}
