package cache

import (
	"context"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const (
	httpTimeout = 30 * time.Second
)

type RLHTTPClient struct {
	client      *http.Client
	RateLimiter *rate.Limiter
}

func (c *RLHTTPClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.RateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewClient(rl *rate.Limiter) *RLHTTPClient {
	c := &RLHTTPClient{
		client:      &http.Client{Timeout: httpTimeout},
		RateLimiter: rl,
	}
	return c
}
