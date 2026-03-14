package ade

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
)

var (
	InitialBackoff      = 2 * time.Second
	HttpTimeout         = 30 * time.Second
	MaxAttempts    uint = 3
	HTTPClient          = &http.Client{Timeout: HttpTimeout}
)

var sem = make(chan struct{}, 5)

func makeRequest(req *http.Request) ([]byte, error) {
	attempts := 0
	body, err := retry.DoWithData(func() ([]byte, error) {
		attempts++
		sem <- struct{}{}
		defer func() { <-sem }()

		slog.Info("requesting", "url", req.URL, "attempt", attempts, "maxAttempts", MaxAttempts)
		response, rqErr := HTTPClient.Do(req)
		if rqErr != nil {
			slog.Error("request failed", "url", req.URL, "err", rqErr)
			return nil, rqErr
		}
		defer func(Body io.ReadCloser) {
			if closeErr := Body.Close(); closeErr != nil {
				slog.Warn("closing response body failed", "url", req.URL, "err", closeErr)
			}
		}(response.Body)

		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			rqBody, readErr := io.ReadAll(response.Body)
			if readErr != nil {
				slog.Warn("reading error response body failed", "url", req.URL, "err", readErr)
				return nil, fmt.Errorf("unexpected status code %d: failed to read error response body: %w", response.StatusCode, readErr)
			}
			return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(rqBody))
		}

		rqBody, rqErr := io.ReadAll(response.Body)
		if rqErr != nil {
			slog.Error("reading response body failed", "url", req.URL, "err", rqErr)
			return nil, rqErr
		}

		return rqBody, nil
	}, retry.Attempts(MaxAttempts), retry.Delay(InitialBackoff), retry.DelayType(retry.BackOffDelay))

	if err != nil {
		slog.Error("all retry attempts exhausted", "url", req.URL, "err", err)
		return nil, err
	}

	slog.Info("request successful", "url", req.URL)
	return body, nil
}
