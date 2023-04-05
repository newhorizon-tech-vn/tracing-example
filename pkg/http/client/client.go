package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	client *http.Client
}

type HTTPClientConfig struct {
	Timeout             time.Duration
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
}

func NewClient(config *HTTPClientConfig) *Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = config.MaxIdleConns
	t.MaxConnsPerHost = config.MaxConnsPerHost
	t.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost

	return &Client{
		client: &http.Client{
			Timeout:   config.Timeout,
			Transport: t,
		},
	}
}

func DefaultClient() *Client {
	return NewClient(&HTTPClientConfig{
		Timeout:             10 * time.Second,
		MaxIdleConns:        100,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 100,
	})
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, int, error) {
	resp, err := otelhttp.Get(ctx, url)
	if err != nil {
		return nil, 0, err
	}

	bytes, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return bytes, resp.StatusCode, nil
}

func (c *Client) GetV2(ctx context.Context, url string) ([]byte, int, error) {
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	bytes, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return bytes, resp.StatusCode, nil
}
