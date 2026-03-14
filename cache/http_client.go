package cache

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
)

const (
	httpTimeout    = 30 * time.Second
	maxAttempts    = 5
	initialBackoff = 5 * time.Second
)

var httpClient = &http.Client{Timeout: httpTimeout}
var sem = make(chan struct{}, 5)

func MakeRequest(logPrefix string, req *http.Request) ([]byte, error) {
	sem <- struct{}{}
	defer func() { <-sem }()

	attempts := 0
	body, err := retry.DoWithData(func() ([]byte, error) {
		attempts++
		log.Printf("[Info] (%s) Requesting (attempt %d/%d)\n", logPrefix, attempts, maxAttempts)
		response, rqErr := httpClient.Do(req)
		if rqErr != nil {
			log.Printf("[Error] (%s) Requesting: %v\n", logPrefix, rqErr)
			return nil, rqErr
		}
		defer func(Body io.ReadCloser) {
			if closeErr := Body.Close(); closeErr != nil {
				log.Printf("[Warn] (%s) Closing response body: %v\n", logPrefix, closeErr)
			}
		}(response.Body)

		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			rqBody, _ := io.ReadAll(response.Body)
			return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(rqBody))
		}

		rqBody, rqErr := io.ReadAll(response.Body)
		if rqErr != nil {
			log.Printf("[Error] (%s) Reading body: %v\n", logPrefix, rqErr)
			return nil, rqErr
		}

		return rqBody, nil
	}, retry.Attempts(maxAttempts), retry.Delay(initialBackoff), retry.DelayType(retry.BackOffDelay))

	if err != nil {
		log.Printf("[Error] (%s) Request failed: %v\n", logPrefix, err)
		return nil, err
	}

	log.Printf("[Info] (%s) Request successful\n", logPrefix)
	return body, nil
}
