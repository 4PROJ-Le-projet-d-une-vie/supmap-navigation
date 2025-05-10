package gis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type ClientOptions struct {
	Timeout time.Duration
}

func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout: 7 * time.Second,
	}
}

func NewClient(baseURL string, options ...ClientOptions) *Client {
	opts := DefaultClientOptions()
	if len(options) > 0 {
		opts = options[0]
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: opts.Timeout},
	}
}

func (c *Client) CalculateRoute(ctx context.Context, routeRequest RouteRequest) (*Route, error) {
	reqURL, err := url.Parse(c.baseURL + "/route")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	body, err := json.Marshal(routeRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var routeResponse RouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&routeResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(routeResponse.Data) == 0 {
		return nil, fmt.Errorf("no route found")
	}

	return &routeResponse.Data[0], nil
}
